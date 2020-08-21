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

// Test PreHandle
//func(p *PingRouter) PreHandle(r ziface.IRequest){
//}

func(p *PingRouter) Handle(r ziface.IRequest){
	fmt.Printf("Call router Handle\n")
	fmt.Printf("recv from client: msgID =%v, data=%v",r.MsgID(),string(r.Data()))

	err := r.Connection().SendMsg(1, []byte("ping..."))
	if err!=nil{
		fmt.Printf("send ping msg failed:%v", err)
	}

}

//func(p *PingRouter) AfterHandle(r ziface.IRequest){
//}

// Test Handle
// Test AfterHandle

func main(){
	s := znet.NewDefaultServer()
	s.RegisterRouter(&PingRouter{})
	s.Run()
}
