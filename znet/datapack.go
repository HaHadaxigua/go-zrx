package znet

import (
	"bytes"
	"encoding/binary"
	"errors"
	"fmt"
	"github.com/HaHadaxigua/go-zrx/setting"
	"github.com/HaHadaxigua/go-zrx/ziface"
)

/**
封包、拆包 模块
解决TCP粘包问题
*/
type DataPack struct {
}

// 构造函数
func NewDataPack() *DataPack {
	return &DataPack{

	}
}

// 获取包的头的长度
func (dp *DataPack) GetHeaderLen() uint32 {
	// DataLen uint32(4 byte) + ID uint32(4 byte)
	return 8
}

/**
	封包
	contentLen + msgID | msgContent
 */
func (dp *DataPack) Pack(msg ziface.IMessage) ([]byte, error){
	// 创建一个存放bytes的缓冲
	dataBuff := bytes.NewBuffer([]byte{})

	// 将dataLen写入dataBuff
	if err := binary.Write(dataBuff, binary.LittleEndian, msg.MsgLen()); err!=nil{
		fmt.Printf("write dataLen failed:%v",err)
		return nil,err
	}

	// 将messageID写入dataBuff
	if err := binary.Write(dataBuff, binary.LittleEndian, msg.MsgID()); err!=nil{
		fmt.Printf("write msgId failed:%v", err)
		return nil,err
	}

	// 将messageContent写入dataBuff
	if err := binary.Write(dataBuff, binary.LittleEndian, msg.MsgData()); err!=nil{
		fmt.Printf("write msgData failed:%v", err)
		return nil,err
	}

	return dataBuff.Bytes(), nil
}

// 拆包
func (dp *DataPack) UnPack(data []byte) (ziface.IMessage, error) {
	// 创建一个从[]byte中读取数据的ioReader
	dataBuff := bytes.NewReader(data)
	// 先只读取header部分数据
	msg := &Message{}
	// 读datalen
	if err:=binary.Read(dataBuff, binary.LittleEndian, &msg.DataLen); err!=nil{
		fmt.Printf("read datalen failed:%v",err)
		return nil,err
	}
	// 读msgid
	if err := binary.Read(dataBuff, binary.LittleEndian, &msg.Id); err!=nil{
		fmt.Printf("read msgid failed:%v",err)
		return nil, err
	}

	// 如果发送的包超出了设置的最大数据包
	if setting.MaxPackageSize>0 && msg.DataLen>uint32(setting.MaxPackageSize) {
		return nil, errors.New("too large msg")
	}

	return msg, nil
}
