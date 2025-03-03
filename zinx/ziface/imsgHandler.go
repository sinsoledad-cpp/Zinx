package ziface

/*
消息管理模块
*/
type IMsgHandler interface {
	//调赴/执行对应的Router消息处理方法
	DoMsgHandler(request IRequest)
	// 为消息添加具体的处理逻辑
	AddRouter(msgID uint32, router IRouter)
}
