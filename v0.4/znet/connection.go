package znet

import (
	"errors"
	"fmt"
	"github.com/suuhui/v0.4/ziface"
	"io"
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
		dp := NewDataPack()
		//读取客户端的头信息
		headData := make([]byte, dp.GetHeadLen())
		if _, err := io.ReadFull(c.GetTCPConnection(), headData); err != nil {
			fmt.Println("read msg head error ", err)
			c.ExitBuffChan <- true
			continue
		}

		//拆包，得到msgId和dataLen放在msg中
		msg, err := dp.Unpack(headData)
		if err != nil {
			fmt.Println("unpack error ", err)
			c.ExitBuffChan <- true
			continue
		}

		//根据dataLen读取data
		var data []byte
		if msg.GetDataLen() > 0 {
			data = make([]byte, msg.GetDataLen())
			if _, err := io.ReadFull(c.GetTCPConnection(), data); err != nil {
				fmt.Print("read msg data error ", err)
				c.ExitBuffChan <- true
				continue
			}
		}
		msg.SetData(data)
		req := Request{conn:c, message: msg}
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

func (c *Connection) SendMsg(msgId uint32, data []byte) error {
	if c.isClosed {
		return errors.New("connection has closed")
	}
	dp := NewDataPack()
	msg, err := dp.Pack(NewMessage(msgId, data))
	if err != nil {
		fmt.Println("Pack error msg id=", msgId)
		return errors.New("pack error msg")
	}
	if  _, err := c.Conn.Write(msg); err != nil {
		fmt.Println("Write msg id ", msgId, " error ")
		c.ExitBuffChan <- true
		return errors.New("conn write error")
	}

	return nil
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
