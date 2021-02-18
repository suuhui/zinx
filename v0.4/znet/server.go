package znet

import (
	"fmt"
	"github.com/suuhui/v0.4/utils"
	"github.com/suuhui/v0.4/ziface"
	"net"
	"time"
)

type Server struct {
	Name string
	IPVersion string
	IP string
	Port int
	Router ziface.IRouter
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
			fmt.Printf("local addr: %s, remote addr: %s\n", conn.LocalAddr().String(), conn.RemoteAddr().String())
			if err != nil {
				fmt.Println("Accept err ", err)
				continue
			}

			//3.2 TODO Server.Start() 设置服务器最大连接控制,如果超过最大连接，那么则关闭此新的连接

			//3.3 TODO Server.Start() 处理该新连接请求的 业务 方法， 此时应该有 handler 和 conn是绑定的
			dealConn := NewConnection(conn, cid, s.Router)
			cid++

			go dealConn.Start()
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

func (s *Server) AddRouter(r ziface.IRouter) {
	s.Router = r
}

func NewServer() ziface.IServer {
	utils.GlobalObject.Reload()

	return &Server{
		Name:      utils.GlobalObject.Name,
		IPVersion: "tcp4",
		IP:       utils.GlobalObject.Host,
		Port:      utils.GlobalObject.TcpPort,
		Router: nil,
	}
}
