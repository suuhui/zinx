package znet

import "github.com/suuhui/v0.7/ziface"

type Request struct {
	conn    ziface.IConnection
	message ziface.IMessage
}

func (r *Request) GetConnection() ziface.IConnection {
	return r.conn
}

func (r *Request) GetData() []byte {
	return r.message.GetData()
}

func (r *Request) GetMsgID() uint32 {
	return r.message.GetMsgId()
}
