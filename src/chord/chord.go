package chord

import (
	"dht/labrpc"
	"sync"
	"time"
)

/*********************************
|    Hint: Feel free to add      |
|  any functions if you need :)  |
*********************************/

type Chord struct{
	mu 				sync.Mutex
	killed 			bool
	Peers			[]*labrpc.ClientEnd // RPC end points of all peers
	Others			[]string
	NodeID			string				
	Location 		int					
	PrevLocation 	int					
	SuccLocation 	int 				
	Map				map[string]string	
	// Your code here.

}

// RPC functions for CHORD
type GetLocationArgs struct {
	NodeID string
	Location int
}
type GetLocationReply struct {
	Location 	 int
	EndIndex	 int
}

func (chd *Chord) SendGetLocation(index int, args *GetLocationArgs, 
reply *GetLocationReply) bool {
    ok := chd.Peers[index].Call("Chord.GetLocation", args, reply)
	return ok
}
func (chd *Chord) GetLocation(args *GetLocationArgs, reply *GetLocationReply) {
	// Your code here.
}

// Please implement your hash code here.
func (chd *Chord) hash(v string) int {

}

// The put function will find out which node to use for
// saving the given value
func (chd *Chord) put(key string, value string) {
	// Your code here.
}

// The get function returns value with given key
func (chd *Chord) get(key string) string, error {
	// Your code here.
}

// Hint: You should put essential functions such as `stabilization()` 
// in this loop, chord will excute them regularly.
func (chd *Chord) routine() {
	Time := 100
	time.Sleep(time.Duration(Time)*time.Millisecond)
	for chd.killed == false {
		// Your code here.
	}
}

// Each node runs “stabilization” regularly to update successor & predecessor
func (chd *Chord) stabilization() {
	// Your code here.
} 

func Make(nodeID string, peers []*labrpc.ClientEnd, otherID []string, self int) *Chord {
	chd := &Chord{}
	chd.NodeID = nodeID
	chd.Peers = peers
	chd.Others = otherID
	chd.Location = BKDRHash(nodeID) % 65536
	chd.Map = make(map[string]string)
	chd.killed = false
	
	go chd.routine()

	return chd
}

 
/***********************************************
|       BELOW ARE CODES NEEDED FOR TEST        |
| 		!!! DO NOT MODIFY THIS PART !!!        |
***********************************************/
func (chd *Chord) Kill() {
	chd.killed = true
}

func (chd *Chord) Killed() bool {
	return chd.killed
}

func (chd *Chord) Tellself (args int, reply int) (string) {
	return chd.NodeID
}

func (chd *Chord) Tellothers (args int, reply int) ([]string) {
	return chd.Others
}