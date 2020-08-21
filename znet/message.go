package znet

type Message struct {
	Id      uint32
	DataLen uint32
	Data    []byte
}

// New方法
func NewMessage(id uint32, data[]byte) *Message{
	return &Message{
		id,
		uint32(len(data)),
		data,
	}
}

// get方法
func (msg *Message) MsgID() uint32 {
	return msg.Id
}
func (msg *Message) MsgLen() uint32  {
	return msg.DataLen
}
func (msg *Message) MsgData() []byte {
	return msg.Data
}

// set方法
func (msg *Message) SetMsgID(id uint32)  {
	msg.Id = id
}
func (msg *Message) SetMsgLen(len uint32) {
	msg.DataLen = len
}
func (msg *Message) SetData(data []byte)   {
	msg.Data = data
}
