/**
模拟客户端
 */
package main

import (
	"fmt"
	"github.com/HaHadaxigua/go-zrx/setting"
	"net"
	"time"
)

func main(){
	fmt.Println("start client_demo ...")
	time.Sleep(1*time.Second)
	// 1. 直接连接到远程服务器 拿到conn
	conn, err := net.Dial("tcp", fmt.Sprintf("%v:%v",setting.IpAddress,setting.Port))
	if err!=nil{
		fmt.Println("connect failed exit")
		return
	}
	// 2. 调用writer写数据到conn中
	for{
		 _,err :=conn.Write([]byte("hello, i'm client"))
		if err!=nil{
			fmt.Println("write failed...")
			return
		}
		buf := make([]byte, 128)
		count, err := conn.Read(buf)
		if err!=nil{
			fmt.Println("read failed...")
			return
		}

		fmt.Printf("server callback:%s, count=%d\n", buf, count)
		time.Sleep(1*time.Second)
	}
}
