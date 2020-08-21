package ziface

/**
封装的请求模块
- 连接信息
- 请求的数据
- 消息id
 */
type IRequest interface {
	Connection() IConnection
	Data() []byte
	MsgID() uint32
}

