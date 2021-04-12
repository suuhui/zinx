package znet

import (
	"fmt"
	"github.com/suuhui/zinx/v0.2/ziface"
	"net"
)

type Connection struct {
	//当前连接的socket tcp套接字
	Conn *net.TCPConn
	//当前连接的id 也可以称作sessionID，全局唯一
	ConnID uint32
	//当前连接的关闭状态
	isClosed bool
	//该连接的处理方法api
	handleAPI ziface.HandleFunc
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
		cnt, err := c.Conn.Read(buf)
		if err != nil {
			fmt.Println("recv buf err ", err)
			c.ExitBuffChan <- true
			continue
		}
		if err := c.handleAPI(c.Conn, buf, cnt); err != nil {
			fmt.Println("connID ", c.ConnID, " handle is error")
			c.ExitBuffChan <- true
			return
		}
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

func NewConnection(conn *net.TCPConn, connId uint32, handleAPI ziface.HandleFunc) *Connection {
	return &Connection{
		Conn:         conn,
		ConnID:       connId,
		handleAPI:    handleAPI,
		isClosed:     false,
		ExitBuffChan: make(chan bool),
	}
}
