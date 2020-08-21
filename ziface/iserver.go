/**
存放接口相关
服务器需要实现的方法：
1. 启动服务器
2. 停止服务器
3. 运行服务
4. 初始化server
 */
package ziface

type Iserver interface {


	// start server
	Start()
	// run server
	Run()
	// stop server
	Stop()
	// register router
	RegisterRouter(msgID uint32, router IRouter)
	// Get connMgr
	GetConnManager() IConnManager

	// register onConnStart
	SetOnConnStart(func(connection IConnection))
	// register onConnStop
	SetOnConnStop(func(connection IConnection))
	// call
	CallOnConnStart(connection IConnection)
	// call
	CallOnConnStop(connection IConnection)
}