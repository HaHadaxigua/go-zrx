/**
基于ZRX开发的服务端应用程序
 */
package main

import "github.com/HaHadaxigua/go-zrx/znet"

func main(){
	s := znet.NewDefaultServer()
	s.Run()
}
