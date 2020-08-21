package ziface
/**
抽象层， 连接管理模块
 */

type IConnManager interface {
	// 添加连接
	AddConn(conn IConnection)
	// 删除连接
	RemoveConn(conn IConnection)
	// 根据connID 获取连接
	GetConn(connID uint32) (IConnection,error)
	// 当前连接总数
	Len() int
	// 删除所有的连接
	ClearAllConn()
}
