package main

import (
	"fmt"
	"github.com/suuhui/v0.6/utils"
	"github.com/suuhui/v0.6/znet"
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

func clientPackTest() {
	conn, err := net.Dial("tcp", "127.0.0.1:8888")
	if err != nil {
		fmt.Println("client dial err:", err)
		return
	}
	dp := znet.NewDataPack()
	msg1 := &znet.Message{
		Id:      0,
		DataLen: 5,
		Data:    []byte{'h', 'e', 'l', 'l', 'o'},
	}

	sendData1, err := dp.Pack(msg1)
	if err != nil {
		fmt.Println("client pack msg1 err:", err)
		return
	}

	msg2 := &znet.Message{
		Id:      1,
		DataLen: 7,
		Data:    []byte{'w', 'o', 'r', 'l', 'd', '!', '!'},
	}
	sendData2, err := dp.Pack(msg2)
	if err != nil {
		fmt.Println("client temp msg2 err:", err)
		return
	}

	//将sendData1，和 sendData2 拼接一起，组成粘包
	sendData1 = append(sendData1, sendData2...)

	//向服务器端写数据
	_, _ = conn.Write(sendData1)

	//客户端阻塞
	select {}
}
