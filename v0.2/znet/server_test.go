package znet

import (
	"fmt"
	"net"
	"testing"
	"time"
)

func ClientTest() {
	fmt.Println("Client Test start...")

	//3秒之后发起测试请求，给服务端开启服务的机会
	time.Sleep(3 * time.Second)

	conn, err := net.Dial("tcp", "127.0.0.1:8888")
	if err != nil {
		fmt.Println("client start fail, exit!")
		return
	}
	for {
		msgByte := []byte("hello, zinx")
		_, err := conn.Write(msgByte)
		if err != nil {
			fmt.Println("write error: ", err)
			return
		}
		buf := make([]byte, len(msgByte))
		cnt, err := conn.Read(buf)
		if err != nil {
			fmt.Println("read buf error: ", err)
			return
		}
		fmt.Printf("server call back : %s, cnt = %d\n", buf,  cnt)
		time.Sleep(1 * time.Second)
	}
}

func TestServer(t *testing.T) {
	s := NewServer("[zinx v0.2]")

	go ClientTest()

	s.Serve()
}
