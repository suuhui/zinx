package znet

import (
	"fmt"
	"github.com/suuhui/zinx/v0.5/ziface"
	"strconv"
)

type MessageHandle struct {
	Apis map[uint32]ziface.IRouter
}

func NewMsgHandler() *MessageHandle {
	return &MessageHandle{Apis: make(map[uint32]ziface.IRouter)}
}

func (mh *MessageHandle) DoMessageHandler(request ziface.IRequest) {
	router, ok := mh.Apis[request.GetMsgID()]
	if !ok {
		fmt.Println("api msgId = " + strconv.Itoa(int(request.GetMsgID())) + " is not found!")
		return
	}
	router.PreHandle(request)
	router.Handle(request)
	router.PostHandle(request)
}

func (mh *MessageHandle) AddRouter(msgId uint32, router ziface.IRouter) {
	if _, ok := mh.Apis[msgId]; ok {
		panic("repeat api, msgId = " + strconv.Itoa(int(msgId)))
	}
	mh.Apis[msgId] = router
	fmt.Println("Add api msgId = ", msgId)
}

