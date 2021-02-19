package main

import (
	"fmt"
	"github.com/suuhui/v0.6/ziface"
	"github.com/suuhui/v0.6/znet"
)

type PingRouter struct {
	znet.BaseRouter
}

func (pr *PingRouter) Handle(request ziface.IRequest) {
	fmt.Println("call PingRouter")
	fmt.Printf("> handle: recv from client: msgId=%d, data=%s\n", request.GetMsgID(), request.GetData())
	err := request.GetConnection().SendMsg(0, []byte("ping..."))
	if err != nil {
		fmt.Println("call back ping error.")
	}
}

type HelloZinxRouter struct {
	znet.BaseRouter
}

func (hzr *HelloZinxRouter) Handle(request ziface.IRequest) {
	fmt.Println("Call HelloZinxHandler")
	fmt.Printf("> handle: recv from client: msgId=%d, data=%s\n", request.GetMsgID(), request.GetData())
	err := request.GetConnection().SendMsg(1, []byte("Hello zinx router V0.5"))
	if err != nil {
		fmt.Println(err)
	}
}

func main() {
	server := znet.NewServer()
	server.AddRouter(0, &PingRouter{})
	server.AddRouter(1, &HelloZinxRouter{})
	server.Serve()
}
