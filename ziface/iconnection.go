package ziface

import "net"
/**
链接模块的抽象层

链接需要实现的方法：
- 启动链接
- 停止链接
- 获取当前链接的conn对象
- 获取链接id
- 获取客户端链接的信息
- 发送数据方法send
- 接收数据方法recv

链接具有的属性：
- socket
- 链接id
- 状态
- 所绑定的业务
- channel: 链接退出(关闭链接是异步的，告知当前链接退出的信息)

*/
type IConnection interface {
	// 启动链接(让当前链接准备开始工作)
	Start()
	// 停止链接
	Stop()
	// 获取当前链接绑定的conn对象
	GetTCPConnect() *net.TCPConn
	// 获取链接id
	GetConnID() uint32
	// 获取客户端链接信息
	GetRemoteAddr() net.Addr
	// 发送数据的方法
	SendMsg(msgId uint32 ,data []byte) error
	//添加带缓冲发送消息接口
	SendBuffMsg(msgId uint32, data []byte) error


	// 设置连接属性
	SetProperty(key string,value interface{})
	// 获取连接属性
	GetProperty(key string)(interface{}, error)
	// 删除连接属性
	DelProperty(key string)

}


// 与链接绑定的业务方法
// 处理链接业务的方法 （抽象函数类型）
// 处理链接的conn对象， 要处理的数据， 数据的长度， 返回错误
type HandleFunc func(*net.TCPConn, []byte,int) error