package chord

import(
	// "net"
	"net/rpc"
	// "log"
)

type Chord struct{
	Peers	[]*rpc.Client 	// RPC end points of all peers
	NodeID	int				// 该节点的唯一ID
	Location int			// 该节点在环上的位置
	Prev	int				// 前置节点的位置
	Succ	int 			// 后继节点的位置
	
	Key		[]int			// 用于保存key
	Value	[]string		// 用于保存value
	Size	int 			// 表示该节点已经存储的大小

	FingerTable [8]int	// finger table
}

//
// RPC Args for chord
//
type GetLocationArgs struct{
	NodeID int
}

type GetLocationReply struct{
	Location int
}

func (chd *Chord) SendGetLocation(node int, args *GetLocationArgs , reply *GetLocationReply) error {
    ok := chd.Peers[node].Call("Chord.GetLocation", args, reply)
	return ok
}

func (chd *Chord) GetLocation(args *GetLocationArgs , reply *GetLocationReply){
	if (args.NodeID == chd.NodeID){
		return
	} else {
		reply.Location = chd.Location
	}
}
//
// The put function will find out which node to use for
// saving the given value 'value'
//
func put(key int, value string){

}

func get(key int){

}



//
// bit is a number indicating how many bits the hash function will use
//
func Make(nodeID int) *Chord {
	chd := &Chord{}
	chd.NodeID = nodeID
	// chd.Peers = peers
	chd.Location = naiveHash(nodeID)
	return chd
}

// 
// for a given node index(int), find his prev and succ node.
// first return value is the prev, second return value is the succ
//
func findLocation(index int)(int, int){
	prev := 0
	succ := 0

	return prev, succ
}

func fillFingerTable(){

}

func naiveHash(val int) int {
	return ((val * 1919810) / 114514) % 65536
}
