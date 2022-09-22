package chord

import (
	"log"
	"net/http"
	"net/rpc"
	"strconv"
	"net/rpc/jsonrpc"

)

func registRPCserver(chd *Chord, port int){
	// 1.注册服务
    // chd := new(Chord)
    // 注册一个rect的服务
    rpc.Register(chd)
    // 2.服务处理绑定到http协议上
    rpc.HandleHTTP()
    // 3.监听服务
    err := http.ListenAndServe(":"+strconv.Itoa(port), nil)
    if err != nil {
        log.Panicln(err)
    }
}

func connectRPCserver(port int) (*rpc.Client) {
	conn, err := jsonrpc.Dial("tcp", ":"+strconv.Itoa(port))
    if err != nil {
        log.Panicln(err)
    }
	return conn
}