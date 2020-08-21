package ziface
/**
	消息管理模块，用来控制我们的路由
 */
type IMsgHandler interface {
	/**
	根据请求的msgID执行调用相应的router执行方法
	 */
	DoMsgHandler(request IRequest)

	/**
	添加路由方法到map集合中
	 */
	RegisterRouter(msgID uint32, router IRouter)

	/**
	启动worker工作池
	 */
	StartWorkerPool()

	SendMsgToTaskQueue(request IRequest)
}
