package ziface

/**
封包、拆包 模块
解决TCP粘包问题
*/

type IDataPack interface {
	// 获取包的头的长度
	GetHeaderLen() uint32
	// 封包
	Pack(IMessage)([]byte, error)
	// 拆包
	UnPack([]byte)(IMessage, error)
}

