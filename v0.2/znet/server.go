package znet

import (
	"errors"
	"fmt"
	"github.com/suuhui/zinx/v0.2/ziface"
	"net"
	"time"
)

type Server struct {
	Name string
	IPVersion string
	IP string
	Port int
}

func (s *Server) Start() {
	fmt.Printf("[START] Server listenner at IP: %s, Port %d, is starting\n", s.IP, s.Port)

	//开启一个goroutine做服务端Listener业务
	go func() {
		//1 获取一个TCP的addr
		addr, err := net.ResolveTCPAddr(s.IPVersion, fmt.Sprintf("%s:%d", s.IP, s.Port))
		if err != nil {
			fmt.Println("resolve tcp addr err: ", err)
			return
		}

		//2 监听服务器地址
		listener, err := net.ListenTCP(s.IPVersion, addr)
		if err != nil {
			fmt.Println("listen", s.IPVersion, " err ", err)
			return
		}
		fmt.Println("start Zinx server ", s.Name, " success, now listening")

		//TODO server.go应该有个自动生成ID的方法
		var cid uint32
		cid = 0

		//3 启动server网络连接业务
		for {
			//3.1 阻塞等待客户端建立连接请求
			conn, err := listener.AcceptTCP()
			if err != nil {
				fmt.Println("Accept err ", err)
				continue
			}

			//3.2 TODO Server.Start() 设置服务器最大连接控制,如果超过最大连接，那么则关闭此新的连接

			//3.3 TODO Server.Start() 处理该新连接请求的 业务 方法， 此时应该有 handler 和 conn是绑定的
			dealConn := NewConnection(conn, cid, CallbackToClient)
			cid++

			go dealConn.Start()
			//暂时做一个最大512字节的回显服务
			//go func() {
			//	for {
			//		buf := make([]byte, 512)
			//		cnt, err := conn.Read(buf)
			//		if err != nil {
			//			fmt.Println("recv buf err ", err)
			//			continue
			//		}
			//
			//		//回显
			//		if _, err := conn.Write(buf[:cnt]); err != nil {
			//			fmt.Println("write back buf err ", err)
			//			continue
			//		}
			//	}
			//}()
		}
	}()
}

func (s *Server) Stop() {
	fmt.Println("[STOP] Zinx server, name: ", s.Name)
	//TODO Server.Stop() 将其他需要清理的连接信息或者其他信息也一并停止或者清理
}

func (s *Server) Serve() {
	s.Start()
	//TODO Server.Serve()是否在启动的时候还要处理其他事情，可以在这里添加

	// 阻塞，否则main goroutine退出，listener的goroutine将退出
	for {
		time.Sleep(10 * time.Second)
	}
}

func NewServer(name string) ziface.IServer {
	return &Server{
		Name:      name,
		IPVersion: "tcp4",
		IP:        "0.0.0.0",
		Port:      8888,
	}
}

func CallbackToClient(conn *net.TCPConn, data []byte, cnt int) error {
	fmt.Println("[Conn Handle] CallbackToClient...")
	if _, err := conn.Write(data[:cnt]); err != nil {
		fmt.Println("write back buf err ", err)
		return errors.New("CallbackToClient error")
	}
	return nil
}
