package main

import (
	"fmt"
	"github.com/suuhui/v0.3/ziface"
	"github.com/suuhui/v0.3/znet"
)

type PingRouter struct {
	znet.BaseRouter
}

func (pr *PingRouter) PreHandle(request ziface.IRequest)  {
	fmt.Println("Call PingRouter PreHandle")
	_, err := request.GetConnection().GetTCPConnection().Write([]byte("before ping...\n"))
	if err != nil {
		fmt.Println("call back ping error.")
	}
}

func (pr *PingRouter) Handle(request ziface.IRequest) {
	fmt.Println("Call PingRouter Handle")
	_, err := request.GetConnection().GetTCPConnection().Write([]byte("ping...\n"))
	if err != nil {
		fmt.Println("call back ping error.")
	}
}

func (pr *PingRouter) PostHandle(request ziface.IRequest)  {
	fmt.Println("Call PingRouter PostHandle")
	_, err := request.GetConnection().GetTCPConnection().Write([]byte("after ping...\n"))
	if err != nil {
		fmt.Println("call back ping error.")
	}
}

func main() {
	server := znet.NewServer("[zinx v0.3]")
	server.AddRouter(&PingRouter{})
	server.Serve()
}
