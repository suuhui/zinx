package ziface

type IMsgHandler interface {
	DoMessageHandler(request IRequest)       //处理消息
	AddRouter(msgId uint32, router IRouter)  //为消息添加具体的处理逻辑
	StartWorkerPool()                        //启动worker工作池
	SendRequestToTaskQueue(request IRequest) //将消息交给消息队列，由worker处理
}
