package ziface

/**
将请求的消息封装到message中，定义抽象的接口
*/
type IMessage interface {
	// get方法
	MsgID() uint32
	MsgLen() uint32
	MsgData() []byte
	// set方法
	SetMsgID(uint32)
	SetMsgLen(uint32)
	SetData([]byte)
}
