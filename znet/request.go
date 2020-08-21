package znet

import "github.com/HaHadaxigua/go-zrx/ziface"

type Request struct {
	// 建立好的连接
	conn ziface.IConnection
	// 客户端请求的数据
	data ziface.IMessage
}

/**
得到当前连接
 */
func(r *Request) Connection() ziface.IConnection{
	return r.conn
}

/**
的带当前的请求数据
 */
func(r *Request) Data() []byte{
	return r.data.MsgData()
}

/**
获取消息id
 */
func(r *Request) MsgID() uint32{
	return r.data.MsgID()
}