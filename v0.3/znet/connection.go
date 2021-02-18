package znet

import (
	"fmt"
	"github.com/suuhui/v0.3/ziface"
	"net"
)

type Connection struct {
	//当前连接的socket tcp套接字
	Conn *net.TCPConn
	//当前连接的id 也可以称作sessionID，全局唯一
	ConnID uint32
	//当前连接的关闭状态
	isClosed bool
	//该连接的处理方法router
	Router ziface.IRouter
	//告知该连接已经退出的channel
	ExitBuffChan chan bool
}

//处理conn读数据的goroutine
func (c *Connection) StartReader() {
	fmt.Println("Reader Goroutine is running")
	defer fmt.Println(c.RemoteAddr().String(), " conn reader exit!")

	defer c.Stop()

	for {
		buf := make([]byte, 512)
		_, err := c.Conn.Read(buf)
		if err != nil {
			fmt.Println("recv buf err ", err)
			c.ExitBuffChan <- true
			continue
		}

		//得到当前客户端请求的Request数据
		req := Request{data:buf, conn:c}
		go func(request ziface.IRequest) {
			c.Router.PreHandle(request)
			c.Router.Handle(request)
			c.Router.PostHandle(request)
		}(&req)
	}
}

func (c *Connection) Start() {
	go c.StartReader()

	for {
		select {
		case <- c.ExitBuffChan:
			return
		}
	}
}

func (c *Connection) Stop() {
	if c.isClosed == true {
		return
	}
	c.isClosed = true

	//TODO Connection.Stop()如果用户注册了该连接的关闭回调业务，那么刺客应该显示调用
	_ = c.Conn.Close()

	c.ExitBuffChan <- true
	close(c.ExitBuffChan)
}

func (c *Connection) GetTCPConnection() *net.TCPConn {
	return c.Conn
}

func (c *Connection) GetConnID() uint32 {
	return c.ConnID
}

func (c *Connection) RemoteAddr() net.Addr {
	return c.Conn.RemoteAddr()
}

func NewConnection(conn *net.TCPConn, connId uint32, router ziface.IRouter) *Connection {
	return &Connection{
		Conn:         conn,
		ConnID:       connId,
		Router:    router,
		isClosed:     false,
		ExitBuffChan: make(chan bool),
	}
}
