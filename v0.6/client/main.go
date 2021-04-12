package main

import (
	"fmt"
	"github.com/suuhui/zinx/v0.6/utils"
	"github.com/suuhui/zinx/v0.6/znet"
	"io"
	"net"
	"strconv"
	"time"
)

func main() {
	fmt.Println("Client starting...")
	time.Sleep(1 * time.Second)

	address := utils.GlobalObject.Host + ":" + strconv.Itoa(utils.GlobalObject.TcpPort)
	conn, err := net.Dial("tcp", address)
	if err != nil {
		fmt.Println("Client start err, exit!", err)
		return
	}

	i := 0
	for {
		dp := znet.NewDataPack()
		msgId := uint32(i % 2)
		i++
		msg, _ := dp.Pack(znet.NewMessage(msgId, []byte("client-" + strconv.Itoa(int(msgId)) + " hello")))
		_, err := conn.Write(msg)
		if err != nil {
			fmt.Println("write error: ", err)
			return
		}
		buf := make([]byte, dp.GetHeadLen())
		_, err = io.ReadFull(conn, buf)
		if err != nil {
			fmt.Println("read buf error ")
			return
		}
		msgHead, err := dp.Unpack(buf)
		if err != nil {
			fmt.Println("server unpack err: ", err)
			return
		}
		if msgHead.GetDataLen() > 0 {
			msg := msgHead.(*znet.Message)
			msg.Data = make([]byte, msg.GetDataLen())

			_, err := io.ReadFull(conn, msg.Data)
			if err != nil  {
				fmt.Println("client unpack data err: ", err)
				return
			}
			fmt.Printf("==> Recv Msg: ID=%d, len=%d, data=%s\n",
				msg.GetMsgId(), msg.GetDataLen(), msg.GetData())
		}
		time.Sleep(1 * time.Second)
	}
	//clientPackTest()
}
