package chord

import (
	"math/rand"
	"net/rpc"
	"encoding/base64"
	crand "crypto/rand"
	"unsafe"

)

type Test struct{
	Nodes	[]*Chord
	Network	[]*rpc.Client
	NodesID	[]int
}

//
// generate nodes and network for tester to use. (Attention: NodeID == NodePort)
//
func (test *Test) makeNodes(num int) (bool) {
	test.Nodes = make([]*Chord, num)
	test.Network = make([]*rpc.Client, num)
	test.NodesID = make([]int, num)
	makeNodeIDs(num)										// genarating NodeID
	for i:=0;i<num;i++ {
		// test.Nodes[i] = Make(test.NodesID[i], test.Network, test.NodesID)		// make n nodes
		// registRPCserver(test.Nodes[i], test.NodesID[i])			// register rpc
		// test.Network[i] = connectRPCserver(test.NodesID[i])		// connect rpc
	}
	return true
}

// 
// generate NodeID. (Attention: NodeID == NodePort)
// 
func makeNodeIDs(num int) ([]string) {		
	NodeID := make([]string, num)
	for i:=0;i<num;i++ {
	    NodeID[i] = randstring(rand.Intn(65536))
		for j:=0;j<i;j++ {					// 如果已经被用了，换一个
			if NodeID[j] == NodeID[i]{
				i--
				continue
			}
		}
	}
	return NodeID
}

func randstring(n int) string {
	b := make([]byte, 2*n)
	crand.Read(b)
	s := base64.URLEncoding.EncodeToString(b)
	return s[0:n]
}

func BKDRHash(str string) (int) {
	seed := 131
	hash := 0
	x := (*[2]uintptr)(unsafe.Pointer(&str))
	h := [3]uintptr{x[0], x[1], x[1]}
	b := *(*[]byte)(unsafe.Pointer(&h))
	for i:=0;i<len(b);i++{
		hash = hash * seed + int(b[i])
	}
	if hash < 0{
		hash = 0 - hash
	}
	return hash
}
