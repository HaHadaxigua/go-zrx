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

//func(p *PingRouter) AfterHandle(r ziface.IRequest){
//}

// Test Handle
// Test AfterHandle

func main(){
	s := znet.NewDefaultServer()
	s.RegisterRouter(0, &PingRouter{})
	s.RegisterRouter(1, &PongRouter{})
	s.Run()
}
