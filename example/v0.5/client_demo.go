/**
模拟客户端
 */
package main

import (
	"fmt"
	"github.com/HaHadaxigua/go-zrx/setting"
	"github.com/HaHadaxigua/go-zrx/znet"
	"io"
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
		// 发送封包的message消息
		dp := znet.NewDataPack()
		binaryMsg, err := dp.Pack(znet.NewMessage(0, []byte("client test message")))
		if err!=nil{
			fmt.Printf("Pack client message failed:%v", err)
			return
		}
		if _, err := conn.Write(binaryMsg); err!=nil{
			fmt.Printf("write binaryMsg failed%v\n", err)
			return
		}

		// 服务器应该给我们返回一个数据
		// 1. 读取流中header部分 得到msgId和dataLen
		binaryHead := make([]byte, dp.GetHeaderLen())
		if _, err := io.ReadFull(conn, binaryHead); err!=nil{
			fmt.Printf("read head error%v",err)
			break
		}
		// 2. 将二进制的head拆包到msg结构体中
		msgHead, err := dp.UnPack(binaryHead)
		if err!=nil{
			fmt.Printf("client unpack msgHead failed %v",err)
			break
		}
		// 3. 读取msgLen长度的数据
		if msgHead.MsgLen()>0{
			msg := msgHead.(*znet.Message)		// 类型推断， 接口变为结构体
			msg.Data = make([]byte, msgHead.MsgLen())

			if _,err:=io.ReadFull(conn, msg.Data); err!=nil{
				fmt.Printf("read msg data error%v\n",err)
				return
			}
			fmt.Printf("receive server msgId:%v, msgLen:%v, msgData:%v, ",
				msg.Id, msg.DataLen, string(msg.Data))
		}


		time.Sleep(1*time.Second)
	}
}
