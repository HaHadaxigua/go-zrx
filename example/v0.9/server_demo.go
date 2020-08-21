/**
基于ZRX开发的服务端应用程序
 */
package main

import (
	"fmt"
	"github.com/HaHadaxigua/go-zrx/ziface"
	"github.com/HaHadaxigua/go-zrx/znet"
)

type PingRouter struct{
	znet.BaseRouter
}

type PongRouter struct{
	znet.BaseRouter
}


// Test PreHandle
//func(p *PingRouter) PreHandle(r ziface.IRequest){
//}

func(p *PingRouter) Handle(r ziface.IRequest){
	fmt.Printf("Call router Handle\n")
	fmt.Printf("recv from client: msgID =%v, data=%v\n",r.MsgID(),string(r.Data()))

	err := r.Connection().SendMsg(200, []byte("ping..."))
	if err!=nil{
		fmt.Printf("send ping msg failed:%v", err)
	}
}

func(p *PongRouter) Handle(r ziface.IRequest){
	fmt.Printf("Call router Handle\n")
	fmt.Printf("recv from client: msgID =%v, data=%v\n",r.MsgID(),string(r.Data()))

	err := r.Connection().SendMsg(300, []byte("pong..."))
	if err!=nil{
		fmt.Printf("send pong msg failed:%v", err)
	}
}

// 创建建立连接之后的狗子函数
func DoConnBegin(connection ziface.IConnection){
	fmt.Printf("--> DoConnBegin...")
	err := connection.SendMsg(202, []byte("钩子函数开始\n"))
	if err!=nil{
		fmt.Printf("狗子函数调用失败")
	}

}

// 销毁连接之后的狗子函数
func DoConnAfter(connection ziface.IConnection){
	fmt.Printf("--> DoConnAfter...\n")
	// 用户下线 告知其他的客户端
	fmt.Printf("connID:%v is Lost\n", connection.GetConnID())
}

func main(){
	s := znet.NewDefaultServer()

	s.SetOnConnStart(DoConnBegin)
	s.SetOnConnStop(DoConnAfter)

	s.RegisterRouter(0, &PingRouter{})
	s.RegisterRouter(1, &PongRouter{})
	s.Run()
}
