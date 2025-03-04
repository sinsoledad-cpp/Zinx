package ziface

/*
消息管理模块
*/
type IMsgHandler interface {
	//调赴/执行对应的Router消息处理方法
	DoMsgHandler(request IRequest)
	// 为消息添加具体的处理逻辑
	AddRouter(msgID uint32, router IRouter)
	// 启动一个Worker工作池(开启工作池的动作只能发生一次,一个zinx框架只能有一个worker工作池
	StartWorkerPool()
	//启动一个Worker工作流程
	StartOneWorker(workerID int, taskQueue chan IRequest)
	//将消息发送给消息任务队列处理
	SendMsgToTaskQueue(request IRequest)
}
