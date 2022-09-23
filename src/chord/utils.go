package chord

import (
	"math/rand"
	"net/rpc"
)

type Test struct{
	NodesID []int
	Nodes 	[]*Chord
	Network []*rpc.Client
}

var portStart = 9900

//
// generate nodes and network for tester to use. (Attention: NodeID == NodePort)
//
func (test *Test) makeNodes(num int) (bool) {
	test.Nodes = make([]*Chord, num)
	test.Network = make([]*rpc.Client, num)
	test.NodesID = make([]int, num)
	test.makeNodeID(num)										// genarating NodeID
	for i:=0;i<num;i++ {
		test.Nodes[i] = Make(test.NodesID[i], test.Network, test.NodesID)		// make n nodes
		registRPCserver(test.Nodes[i], test.NodesID[i])			// register rpc
		test.Network[i] = connectRPCserver(test.NodesID[i])		// connect rpc
	}
	return true
}

// 
// generate NodeID. (Attention: NodeID == NodePort)
// 
func (test *Test) makeNodeID(num int){		// NodeID与RPC端口绑定（同构）
	for i:=0;i<num;i++ {
	    test.NodesID[i] = rand.Intn(1000) + 9000
		for j:=0;j<i;j++ {					// 如果已经被用了，也换一个
			if test.NodesID[j] == test.NodesID[i]{
				i--
				continue
			}
		}
	}
}
