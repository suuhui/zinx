package main

import (
	"fmt"
	"net"
	"time"
)

func main() {
	fmt.Println("Client starting...")
	time.Sleep(3 * time.Second)

	conn, err := net.Dial("tcp", "127.0.0.1:8888")
	if err != nil {
		fmt.Println("Client start err, exit!")
		return
	}

	for {
		_, err := conn.Write([]byte("client zinx v0.3"))
		if err != nil {
			fmt.Println("write error: ", err)
			return
		}
		buf := make([]byte, 512)
		cnt, err := conn.Read(buf)
		if err != nil {
			fmt.Println("read buf error ")
			return
		}
		buf = buf[:cnt]
		fmt.Printf("server call back: %s, cnt = %d\n", buf, cnt)
		time.Sleep(1 * time.Second)
	}
}
