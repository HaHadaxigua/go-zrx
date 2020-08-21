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
func(p *PingRouter) PreHandle(r ziface.IRequest){
	fmt.Printf("Call router PreHandle\n")
	_, err := r.Connection().GetTCPConnect().Write([]byte("before ping...\n"))
	if err!=nil{
		fmt.Printf("callback before ping error\n")
	}

}

func(p *PingRouter) Handle(r ziface.IRequest){
	fmt.Printf("Call router Handle\n")
	_, err := r.Connection().GetTCPConnect().Write([]byte(" ping...\n"))
	if err!=nil{
		fmt.Printf("callback ping error\n")
	}
}

func(p *PingRouter) AfterHandle(r ziface.IRequest){
	fmt.Printf("Call router AfterHandle\n")
	_, err := r.Connection().GetTCPConnect().Write([]byte("after ping..."))
	if err!=nil{
		fmt.Printf("callback after ping error")
	}
}

// Test Handle
// Test AfterHandle

func main(){
	s := znet.NewDefaultServer()
	s.RegisterRouter(&PingRouter{})
	s.Run()
}
