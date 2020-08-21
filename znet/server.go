package znet

import (
	"fmt"
	"github.com/HaHadaxigua/go-zrx/setting"
	"github.com/HaHadaxigua/go-zrx/ziface"
	"net"
)

/**
Iserver接口的实现
*/
type Server struct {
	// 服务器版本信息
	Version string
	// name
	Name string
	// ip_version
	IPVersion string
	// ip
	IpAddress string
	// port
	Port int
	// 当前server的消息管理模块
	MsgHandler ziface.IMsgHandler
	// 当前server的连接管理模块
	ConnMgr ziface.IConnManager
	// 创建连接之后会自动调用hock函数
	OnConnStart func(conn ziface.IConnection)
	// 销毁连接之后会自动销毁hock函数
	OnConnStop func(conn ziface.IConnection)

}

/**
启动基本服务器
1. 创建addr
2. 创建listener
3. 监听listener的连接
*/
func (s *Server) Start() {
	fmt.Printf("[Start] Server Listener at IP:%s, port:%d \n", s.IpAddress, s.Port)
	var cid uint32 = 0
	go func() {
		// 开启消息队列和worker工作池
		s.MsgHandler.StartWorkerPool()

		// 1. 获取一个tcp addr
		addr, err := net.ResolveTCPAddr(s.IPVersion, fmt.Sprintf("%s:%d", s.IpAddress, s.Port))
		if err != nil {
			fmt.Printf("resolve tcp addr error:%v\n", err)
			return
		}
		// 2. 根据addr 来创建一个listener
		listener, err := net.ListenTCP(s.IPVersion, addr)
		if err != nil {
			fmt.Printf("resolve tcp listener error:%v\n", err)
			return
		}
		// 3. listener:监听获取客户端的请求连接
		for {
			conn, err := listener.AcceptTCP()
			if err != nil {
				fmt.Printf("resolve tcp connection error:%v\n", err)
				continue
			}
			// 进行连接最大数量的判断，如果大于最大连接，则关闭此新连接
			if s.ConnMgr.Len() >= setting.MaxConn {
				// TODO 给客户端响应一个超出最大连接的错误消息
				fmt.Printf("Connection more than max value.")
				conn.Close()
				continue
			}
			delConn := NewConnection(s, conn, cid, s.MsgHandler)
			cid++

			go delConn.Start()
		}
	}()
}

func (s *Server) Run() {
	s.Start()
	// 做一些业务相关的事情
	select {}
	s.Stop()
}

func (s *Server) Stop() {
	defer fmt.Printf("[server stopper]\n")
	// 停止服务器、回收资源等
	fmt.Printf("[server is stopping...]\n")
	s.ConnMgr.ClearAllConn()
}

// 注册路由方法
func (s *Server) RegisterRouter(msgID uint32, router ziface.IRouter) {
	s.MsgHandler.RegisterRouter(msgID, router)
	fmt.Printf("register router:%v success", router)
}

// 获取当前server的连接管理器
func (s *Server) GetConnManager() ziface.IConnManager {
	return s.ConnMgr
}

/**
一个默认的Server模板
*/
func NewDefaultServer() ziface.Iserver {
	s := &Server{
		Version:    setting.Version,
		Name:       setting.Name,
		IPVersion:  setting.IpVersion,
		IpAddress:  setting.IpAddress,
		Port:       setting.Port,
		MsgHandler: NewMsgHandler(),
		ConnMgr:    NewConnManager(),
	}
	return s
}

// register onConnStart
func(s *Server)SetOnConnStart(hockFunc func(connection ziface.IConnection)){
	s.OnConnStart = hockFunc
}
// register onConnStop
func(s *Server)SetOnConnStop(hockFunc func(connection ziface.IConnection)){
	s.OnConnStop = hockFunc
}
// call hockFunc
func(s *Server)CallOnConnStart(connection ziface.IConnection){
	if s.OnConnStart!=nil{
		fmt.Printf("---> Call onConnStart")
		s.OnConnStart(connection)
	}

}
// call hockFunc
func(s *Server)CallOnConnStop(connection ziface.IConnection){
	if s.OnConnStop!=nil{
		fmt.Printf("---> Call onConnStop")
		s.OnConnStop(connection)
	}
}






