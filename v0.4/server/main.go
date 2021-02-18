package main

import (
	"fmt"
	"github.com/suuhui/v0.4/ziface"
	"github.com/suuhui/v0.4/znet"
	"io"
	"net"
)

type PingRouter struct {
	znet.BaseRouter
}

func (pr *PingRouter) PreHandle(request ziface.IRequest)  {
	fmt.Printf("> prehandle\n")
}

func (pr *PingRouter) Handle(request ziface.IRequest) {
	fmt.Printf("> handle: recv from client: msgId=%d, data=%s\n", request.GetMsgID(), request.GetData())
	err := request.GetConnection().SendMsg(1, []byte("ping...\n"))
	if err != nil {
		fmt.Println("call back ping error.")
	}
}

func (pr *PingRouter) PostHandle(request ziface.IRequest)  {
	fmt.Println("> postHandle")
}

func main() {
	server := znet.NewServer()
	server.AddRouter(&PingRouter{})
	server.Serve()
}

func serverTest() {
	listener, err := net.Listen("tcp", "127.0.0.1:8888")
	if err != nil {
		fmt.Println("server listen err: ", err)
		return
	}
	for {
		conn, err := listener.Accept()
		if err != nil {
			fmt.Println("server accept error: ", err)
		}

		go func(conn net.Conn) {
			dp := znet.NewDataPack()
			for {
				headData := make([]byte, dp.GetHeadLen())
				_, err := io.ReadFull(conn, headData)
				if err != nil {
					fmt.Println("read head error")
					return
				}
				msgHead, err := dp.Unpack(headData)
				if err != nil {
					fmt.Println("server unpack error: ", err)
					return
				}
				if msgHead.GetDataLen() > 0 {
					msg := msgHead.(*znet.Message)
					msg.Data = make([]byte, msg.GetDataLen())
					_, err := io.ReadFull(conn, msg.Data)
					if err != nil {
						fmt.Println("server unpack data err: ", err)
						return
					}
					fmt.Printf("==> Recv Msg: ID=%d, len=%d, data=%s\n", msg.Id, msg.DataLen, msg.Data)
				}
			}
		}(conn)
	}
}
