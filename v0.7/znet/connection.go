package znet

import (
	"errors"
	"fmt"
	"github.com/suuhui/zinx/v0.7/utils"
	"github.com/suuhui/zinx/v0.7/ziface"
	"io"
	"net"
	"sync"
)

type Connection struct {
	//当前conn属于哪个server
	TcpServer ziface.IServer
	//当前连接的socket tcp套接字
	Conn *net.TCPConn
	//当前连接的id 也可以称作sessionID，全局唯一
	ConnID uint32
	//当前连接的关闭状态
	isClosed bool
	//消息管理msgId和对应处理方法的消息管理模块
	MsgHandler ziface.IMsgHandler
	//告知该连接已经退出的channel
	ExitBuffChan chan bool
	//无缓冲管道，用于读写两个goroutine之间的消息通信
	msgChan chan []byte
	//有缓冲管道，用于读写两个goroutine之间的消息通信
	msgBufferChan chan []byte
	//链接属性
	property map[string]interface{}
	//保护链接属性修改的锁
	propertyLock sync.RWMutex
}

func NewConnection(server ziface.IServer, conn *net.TCPConn, connId uint32, msgHandler ziface.IMsgHandler) *Connection {
	c := &Connection{
		TcpServer:     server,
		Conn:          conn,
		ConnID:        connId,
		MsgHandler:    msgHandler,
		isClosed:      false,
		ExitBuffChan:  make(chan bool),
		msgChan:       make(chan []byte),
		msgBufferChan: make(chan []byte, utils.GlobalObject.MaxMsgChanLen),
		property: make(map[string]interface{}),
	}
	c.TcpServer.GetConnMgr().Add(c)
	return c
}

//处理conn读数据的goroutine
func (c *Connection) StartReader() {
	fmt.Println("Reader Goroutine is running")
	defer fmt.Println(c.RemoteAddr().String(), " conn reader exit!")
	defer c.Stop()

	for {
		dp := NewDataPack()
		//读取客户端的头信息,头信息中只包含两个uint32类型长度的字节
		headData := make([]byte, dp.GetHeadLen())
		//ReadFull在conn中读取数据，并且把读取指针移动读取的位数
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
		req := Request{conn: c, message: msg}
		//如果启动了工作池，则交由工作池处理
		if utils.GlobalObject.WorkerPoolSize > 0 {
			c.MsgHandler.SendRequestToTaskQueue(&req)
		} else {
			//否则直接使用绑定的处理方式处理
			go c.MsgHandler.DoMessageHandler(&req)
		}
	}
}

func (c *Connection) StartWriter() {
	fmt.Println("[Write goroutine is running]")
	defer fmt.Println(c.RemoteAddr().String(), "[conn writer exit!]")

	for {
		select {
		case data, _ := <-c.msgChan:
			//有数据要写给客户端
			if _, err := c.Conn.Write(data); err != nil {
				fmt.Println("Send Data error: ", err, " conn write exit")
				return
			}
		case data, ok := <-c.msgBufferChan:
			if ok {
				if _, err := c.Conn.Write(data); err != nil {
					fmt.Println("Send Data error: ", err, " conn write exit")
					return
				}
			} else {
				fmt.Println("msgBuffChan is Closed")
			}
		case <-c.ExitBuffChan:
			return
		}
	}
}

func (c *Connection) Start() {
	go c.StartReader()
	go c.StartWriter()
	c.TcpServer.CallOnConnStart(c)

	for {
		select {
		case <-c.ExitBuffChan:
			c.Stop()
			return
		}
	}
}

func (c *Connection) Stop() {
	if c.isClosed == true {
		return
	}
	c.isClosed = true

	//用户注册了该连接的关闭回调业务，那么应该显示调用
	c.TcpServer.CallOnConnStop(c)
	_ = c.Conn.Close()

	c.ExitBuffChan <- true
	c.TcpServer.GetConnMgr().Remove(c)
	close(c.ExitBuffChan)
	close(c.msgBufferChan)
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
	c.msgChan <- msg

	return nil
}

func (c *Connection) SendBufMsg(msgId uint32, data []byte) error {
	if c.isClosed {
		return errors.New("connection has closed")
	}
	dp := NewDataPack()
	msg, err := dp.Pack(NewMessage(msgId, data))
	if err != nil {
		fmt.Println("Pack error msg id=", msgId)
		return errors.New("pack error msg")
	}
	c.msgBufferChan <- msg
	return nil
}

func (c *Connection) SetProperty(key string, value interface{}) {
	c.propertyLock.Lock()
	defer c.propertyLock.Unlock()
	c.property[key] = value
}

func (c *Connection) GetProperty(key string) (interface{}, error) {
	c.propertyLock.RLock()
	defer c.propertyLock.RUnlock()

	if val, ok := c.property[key]; ok {
		return val, nil
	}
	return nil, errors.New("property not found")
}

func (c *Connection) RemoveProperty(key string) {
	c.propertyLock.Lock()
	defer c.propertyLock.Unlock()

	delete(c.property, key)
}
