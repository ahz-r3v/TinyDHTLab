package chord

import (
	"testing"
)




func TestMakeNode(t *testing.T){
	test := &Test{}
	test.makeNodes(3)
	test.Nodes[1].hello()
}
