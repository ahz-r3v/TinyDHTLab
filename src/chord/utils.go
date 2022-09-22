package chord

import (
	"bytes"
	"fmt"
	"math/rand"
	"os/exec"
	"regexp"
	"strconv"
	"strings"
)

type Test struct{
	NodesID []int
}

func (test *Test) makeNodeID(num int){		//NodeID与RPC端口绑定（同构）
	for i:=0;i<num;i++ {
	    test.NodesID[i] = rand.Intn(65535)
	    if PortAvailable(test.NodesID[i]) {	//如果端口被占用着，换一个
			i--
			continue
	    }
	}
}

//
// 检查端口是否被占用
//
func PortAvailable(port int) bool {
    checkStatement := fmt.Sprintf(`netstat -apt | grep -q %d ; echo $?`, port)
    output, err := exec.Command("sh", "-c", checkStatement).CombinedOutput()
    if err != nil {
        return false
    }
    result, err := strconv.Atoi(strings.TrimSuffix(string(output), "n"))
    if err != nil {
        return false
    }
    if result == 0 {	// 0代表不可用，1代表可用
        return false
    }

    return true
}