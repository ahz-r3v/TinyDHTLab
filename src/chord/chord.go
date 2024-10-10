package chord

import (
	// "net"
	"dht/labrpc"
	// "fmt"
	"sync"
	"time"
	"errors"
	// "log"
)

type Chord struct{
	mu 				sync.Mutex
	killed 			bool
	Peers			[]*labrpc.ClientEnd // RPC end points of all peers
	selfIndex		int					// self index in the clients ends
	Others			[]string

	NodeID			string				// 该节点的唯一ID
	Location 		int					// 该节点在环上的位置
	PrevLocation 	int					// 前置节点的位置
	PrevEndIndex	int					// 前置节点的RPC索引
	SuccLocation 	int 				// 后继节点的位置
	SuccEndIndex	int					// 后继节点的RPC索引

	Map				map[string]string	// 节点内map
	Size			int 				// 表示该节点已经存储的大小

	FingerTable 	[8]int				// finger table
	Located			bool				// 该节点是否已经确定下位置
}

func (chd *Chord) Kill() {
	chd.killed = true
}
func (chd *Chord) Killed() bool {
	return chd.killed
}

func (chd *Chord) hash(k string) int {
	return BKDRHash(k)
}

//
// RPC methods for chord
//
type GetLocationArgs struct {
	NodeID string
	Location int
}
type GetLocationReply struct {
	Location 	 int
	EndIndex	 int
}

func (chd *Chord) SendGetLocation(index int, args *GetLocationArgs , reply *GetLocationReply) bool {
    ok := chd.Peers[index].Call("Chord.GetLocation", args, reply)
	return ok
}
func (chd *Chord) GetLocation(args *GetLocationArgs , reply *GetLocationReply) {
	if args.Location == chd.Location {
		reply.Location = -1
		reply.EndIndex = -1
		return
	}
	reply.Location = chd.Location
	reply.EndIndex = chd.selfIndex
	tl := args.Location
	location := chd.Location
	succ := chd.SuccLocation
	prev := chd.PrevLocation
	if tl > location {
		if succ > location && tl < succ {
			succ = tl
		} else if succ < location{
			succ = tl
		}
		if prev < location {

		} else if prev > location && tl > prev{
			prev = tl
		}
	} else if tl < location {
		if prev < location && tl > prev{
			prev = tl
		} else if prev > location{
			prev = tl
		}
		if succ > location{

		} else if succ < location && tl < succ {
			succ = tl
		}
	}
}

type HandleEntryArgs struct {
	From int
	Function string
	Key string
	Value string
}
type HandleEntryReply struct {
	Confirm bool
	Value string
	Err error
	Location int
	EndIndex int
}

func (chd *Chord) SendHandleEntry(index int, args *HandleEntryArgs , reply *HandleEntryReply) bool {
	ok := chd.Peers[index].Call("Chord.HandleEntry",args,reply)
	return ok
}

func (chd *Chord) HandleEntry(args *HandleEntryArgs, reply *HandleEntryReply){
	if args.Function == "put" {
		// fmt.Println("<RPC> Put request from:",args.From,"to:",chd.selfIndex)
		chd.put(args.Key, args.Value)
	} else if args.Function == "get" {
		reply.Value, reply.Err = chd.get(args.Key)
	}
	reply.Confirm = true
} 


//
// The put function will find out which node to use for
// saving the given value 'value'
//
func (chd *Chord) put(key string, value string) {
	// fmt.Println("Put request:",chd.selfIndex)
	hash := BKDRHash(key) % 65536
	switch {
		case chd.PrevLocation < chd.Location:
			if hash > chd.PrevLocation && hash < chd.Location {
				// fmt.Println("key:",key,"value:",value,"saved here:",chd.selfIndex)
				chd.Map[key] = value
			} else {
				args  := &HandleEntryArgs{
					From: chd.selfIndex,
					Function: "put",
					Key: key,
					Value: value,
				}
				reply := &HandleEntryReply{}
				chd.mu.Lock()
				// fmt.Println("key:",key,"value:",value,"throw to node:",chd.SuccEndIndex)
				chd.SendHandleEntry(chd.SuccEndIndex, args, reply)
				chd.mu.Unlock()
			}
			break
		case chd.PrevLocation > chd.Location:
			if hash > chd.PrevLocation || hash < chd.Location {
				// fmt.Println("key:",key,"value:",value,"saved here:",chd.selfIndex)
				chd.Map[key] = value
			} else {
				args  := &HandleEntryArgs{
					From: chd.selfIndex,
					Function: "put",
					Key: key,
					Value: value,
				}
				reply := &HandleEntryReply{}
				chd.mu.Lock()
				// fmt.Println("key:",key,"value:",value,"throw to node:",chd.SuccEndIndex)
				chd.SendHandleEntry(chd.SuccEndIndex, args, reply)
				chd.mu.Unlock()
			}
			break
		default:
			// fmt.Println("This should not be printed.")
	}
}

func (chd *Chord) get(key string) (string, error) {
	hash := BKDRHash(key) % 65536
	switch {
	case chd.PrevLocation < chd.Location:
		if hash > chd.PrevLocation && hash < chd.Location {
			v, ok := chd.Map[key]
			if ok {
				return v, nil
			}
		} else {
			args  := &HandleEntryArgs{
				From: chd.selfIndex,
				Function: "get",
				Key: key,
			}
			reply := &HandleEntryReply{}
			chd.mu.Lock()
			ok := chd.SendHandleEntry(chd.SuccEndIndex, args, reply)
			chd.mu.Unlock()
			if ok {
				return reply.Value, reply.Err
			}
		}
		break
	case chd.PrevLocation > chd.Location:
		if hash > chd.PrevLocation || hash < chd.Location {
			v, ok := chd.Map[key]
			if ok {
				return v, nil
			}
		} else {
			args  := &HandleEntryArgs{
				From: chd.selfIndex,
				Function: "get",
				Key: key,
			}
			reply := &HandleEntryReply{}
			chd.mu.Lock()
			ok := chd.SendHandleEntry(chd.SuccEndIndex, args, reply)
			chd.mu.Unlock()
			if ok {
				return reply.Value, reply.Err
			}
		}
		break
	default:
		// fmt.Println("This should not be printed.")
	}
	return "", errors.New("No record!")
}

// func (chd *Chord) find(key string){
// 	hash := BKDRHash(key) % 65536
	
// }



//
// You should put essential functions such as `stabilization()` in
// this loop, chord should excute them regularly
//
func (chd *Chord) routine() {
	Time := 100
	time.Sleep(time.Duration(Time)*time.Millisecond)
	for chd.killed == false {
		// Your code here.
		chd.stabilization()
		Time := 300
		time.Sleep(time.Duration(Time)*time.Millisecond)
		
	}
}

//
// 每个节点会定期执行“stabilization”以更新 successor/predecessor。
//
func (chd *Chord) stabilization() {
	length := len(chd.Peers)
	location := chd.Location
	// fmt.Println("me:",chd.selfIndex,",len:",len(chd.Peers))
	locs := make([]int, length-1)							// other nodes' locations\
	ends := make([]int, length-1)	
	j := 0
	for i:=0;i<length;i++{
		args := &GetLocationArgs{
			// NodeID: chd.NodeID,
			Location: chd.Location,
		}
		reply := &GetLocationReply{}
		chd.SendGetLocation(i, args, reply)
		if (reply.Location == -1 || reply.EndIndex == -1){
			continue
		}
		locs[j] = reply.Location
		ends[j] = reply.EndIndex
		j++
	}

	prev := locs[0]
	succ := locs[0]
	prevEndIndex := ends[0]
	succEndIndex := ends[0]

	for i:=0;i<len(locs);i++{

		tl := locs[i]
		ei := ends[i]
		if tl > location {
			if succ > location && tl < succ {
				succ = tl
				succEndIndex = ei
			} else if succ < location{
				succ = tl
				succEndIndex = ei
			}

			if prev < location {

			} else if prev > location && tl > prev{
				prev = tl
				prevEndIndex = ei
			}
		} else if tl < location {
			if prev < location && tl > prev{
				prev = tl
				prevEndIndex = ei
			} else if prev > location{
				prev = tl
				prevEndIndex = ei
			}

			if succ > location{

			} else if succ < location && tl < succ {
				succ = tl
				succEndIndex = ei
			}
		}
		
	}
	// fmt.Println("me:", chd.Location,", prev:", prev,", succ:", succ)
	// fmt.Println("me:", chd.selfIndex,", prev:", prevEndIndex,", succ:", succEndIndex)
	chd.mu.Lock()
	chd.PrevLocation = prev
	chd.SuccLocation = succ
	chd.PrevEndIndex = prevEndIndex
	chd.SuccEndIndex = succEndIndex
	chd.Located = true
	chd.mu.Unlock()
}

//
// 每个节点会定期执行“stabilization”以更新 finger table。
//
func (chd *Chord) fixFinger () {

}


//
// bit is a number indicating how many bits the hash function will use
//
func Make(nodeID string, peers []*labrpc.ClientEnd, otherID []string, self int) *Chord {
	chd := &Chord{}
	chd.Located = true
	chd.NodeID = nodeID
	chd.Peers = peers
	chd.Others = otherID
	chd.Location = BKDRHash(nodeID) % 65536
	chd.Map = make(map[string]string)
	chd.selfIndex = self
	chd.killed = false
	
	
	go chd.routine()

	return chd
}

 
//
// ******* BELOW ARE CODES ONLY FOR TEST ********
// for a given node index(int), find his prev and succ node.
// first return value is the prev, second return value is the succ
//
func getLocation(index int)(int, int){
	prev := 0
	succ := 0

	return prev, succ
}

func fillFingerTable(){

}

func (chd *Chord) Tellself (args int, reply int) (string) {
	return chd.NodeID
}
func (chd *Chord) Tellothers (args int, reply int) ([]string) {
	return chd.Others
}