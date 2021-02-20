package znet

import (
	"fmt"
	"github.com/suuhui/v0.7/utils"
	"github.com/suuhui/v0.7/ziface"
	"strconv"
)

type MessageHandle struct {
	Apis map[uint32]ziface.IRouter
	WorkerPoolSize uint32 //worker数量
	TaskQueue []chan ziface.IRequest //worker任务队列
}

func NewMsgHandler() *MessageHandle {
	poolSize := utils.GlobalObject.WorkerPoolSize
	return &MessageHandle{
		Apis: make(map[uint32]ziface.IRouter),
		WorkerPoolSize: poolSize,
		TaskQueue: make([]chan ziface.IRequest, poolSize),
	}
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

func (mh *MessageHandle) StartWorkerPool() {
	for i := 0; i < int(mh.WorkerPoolSize); i++ {
		mh.TaskQueue[i] = make(chan ziface.IRequest, utils.GlobalObject.MaxWorkerTaskLen)
		go mh.startOneWorker(i, mh.TaskQueue[i])
	}
}

//将消息交给taskQueue，由worker处理
func (mh *MessageHandle) SendRequestToTaskQueue(request ziface.IRequest) {
	//根据ConnID决定该连接应该由哪个worker处理
	//轮询平均分配
	workerId := request.GetConnection().GetConnID() % mh.WorkerPoolSize
	fmt.Println("Add ConnID=", request.GetConnection().GetConnID()," request msgID=", request.GetMsgID(), "to workerID=", workerId)
	mh.TaskQueue[workerId] <- request
}

func (mh *MessageHandle) startOneWorker(i int, taskQueue chan ziface.IRequest) {
	fmt.Println("WorkerId = ", i, " has started")

	for {
		select {
		case request := <-taskQueue:
			mh.DoMessageHandler(request)
		}
	}
}