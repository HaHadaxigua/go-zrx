package znet

import (
	"errors"
	"fmt"
	"github.com/HaHadaxigua/go-zrx/ziface"
	"sync"
)

/**
	连接管理模块
 */
type ConnManager struct{
	connections map[uint32] ziface.IConnection	// 管理的连接集合
	connLock sync.RWMutex						// 保护集合的读写锁
}

func NewConnManager() *ConnManager{
	return &ConnManager{
		connections: make(map[uint32]ziface.IConnection),
	}
}

func (this *ConnManager) AddConn(conn ziface.IConnection) {
	// 添加加写锁
	this.connLock.Lock()
	defer this.connLock.Unlock()
	// 将conn加入到connManager中
	this.connections[conn.GetConnID()] = conn
	fmt.Printf("connection:%v add to ConnManager successful\n", conn.GetConnID())
}

func (this *ConnManager) RemoveConn(conn ziface.IConnection) {
	// 删除加写锁
	this.connLock.Lock()
	defer this.connLock.Unlock()
	// 删除connection
	delete(this.connections, conn.GetConnID())

	fmt.Printf("connection:%v remove from ConnManage\n", conn.GetConnID())
}

func (this *ConnManager) GetConn(connID uint32) (ziface.IConnection, error) {
	// 加读锁
	this.connLock.RLock()
	defer this.connLock.RUnlock()
	if conn,ok := this.connections[connID]; ok{
		return conn, nil
	}else{
		return nil, errors.New("connection not found")
	}

}

func (this *ConnManager) Len() int {
	return len(this.connections)
}

func (this *ConnManager) ClearAllConn() {
	// 删除加写锁
	this.connLock.Lock()
	defer this.connLock.Unlock()

	for connID, conn := range this.connections{
		conn.Stop()
		delete(this.connections, connID)
	}

	fmt.Printf("clear all connections\n")
}

