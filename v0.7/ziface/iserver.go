package ziface

type IServer interface {
	//启动服务器
	Start()
	//停止服务器
	Stop()
	//开启业务服务
	Serve()
	//路由功能：给当前服务注册一个路由业务方法，供客户端连接处理使用
	AddRouter(msgId uint32, router IRouter)
	//获取连接管理
	GetConnMgr() IConnManager
	//执行创建连接后的hook函数
	CallOnConnStart(conn IConnection)
	//执行连接断开时的hook函数
	CallOnConnStop(conn IConnection)
}
