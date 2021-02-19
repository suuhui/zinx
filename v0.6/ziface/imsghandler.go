package ziface

type IMsgHandler interface {
	DoMessageHandler(request IRequest) //处理消息
	AddRouter(msgId uint32, router IRouter) //为消息添加具体的处理逻辑
}
