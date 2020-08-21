package znet

import (
	"errors"
	"fmt"
	"github.com/HaHadaxigua/go-zrx/setting"
	"github.com/HaHadaxigua/go-zrx/ziface"
	"io"
	"net"
	"sync"
)

/**
链接模块
*/
type Connection struct {
	// 当前connection属于哪个server
	TcpServer ziface.Iserver
	// 当前链接的socket
	Conn *net.TCPConn
	// 链接的id
	ConnID uint32
	// 链接状态
	isClosed bool
	// 告知当前链接已经退出、停止的channel
	exitChan chan bool
	// 一个无缓冲管道，用于读、写goroutine之间的通信
	msgChan chan []byte
	//有缓冲管道，用于读、写两个goroutine之间的消息通信
	msgBuffChan chan []byte
	// 消息的管理msgID 和对应的处理业务API关系
	MsgHandler ziface.IMsgHandler
	// 连接属性集合
	property map[string]interface{}
	// 保护连接属性的锁
	propertyLock sync.RWMutex
}

// 初始化链接模块的方法
func NewConnection(server ziface.Iserver, conn *net.TCPConn, connID uint32, handler ziface.IMsgHandler) *Connection {
	c := &Connection{
		TcpServer:  server,
		Conn:       conn,
		ConnID:     connID,
		MsgHandler: handler,
		isClosed:   false,
		exitChan:   make(chan bool, 1),
		msgChan:    make(chan []byte),
		msgBuffChan: make(chan[] byte, setting.MaxTaskQueueLen),
		property:   make(map[string]interface{}),
	}

	// 将Conn加入到connManager
	c.TcpServer.GetConnManager().AddConn(c)
	return c
}

func (c *Connection) StartReader() {
	fmt.Printf("[Reader Goroutine ]start...ConnID=%v\n", c.ConnID)
	defer fmt.Printf("[Reader is exit!]%v Stop...ConnID=%v\n", c.ConnID, c.GetRemoteAddr().String())
	defer c.Stop()

	for {
		// 创建拆包装包对象
		dp := NewDataPack()
		// 1. 读取客户端的msgHead
		headData := make([]byte, dp.GetHeaderLen())
		if _, err := io.ReadFull(c.GetTCPConnect(), headData); err != nil {
			fmt.Printf("read msg head failed:%v\n", err)
			break
		}
		// 2. 拆包，得到msgId和msgDataLen放在msg消息中
		msg, err := dp.UnPack(headData)
		if err != nil {
			fmt.Printf("unpack msgHead error:%v\n", err)
			break
		}
		// 3. 根绝dataLen再次读取
		var data []byte
		if msg.MsgLen() > 0 {
			data = make([]byte, msg.MsgLen())
			_, err := io.ReadFull(c.GetTCPConnect(), data)
			if err != nil {
				fmt.Printf("read msg data failed:%v\n", err)
				break
			}
		}
		msg.SetData(data)
		// 4. 得到request请求数据
		req := Request{
			c,
			msg,
		}
		// 5. 将req交给工作池处理
		if setting.WorkerPoolSize > 0 {
			// 工作池已经开启
			c.MsgHandler.SendMsgToTaskQueue(&req)
		} else {
			go c.MsgHandler.DoMsgHandler(&req)
		}
	}
}

/**
写消息Goroutine，专门发送给客户端消息的模块
*/
func (c *Connection) StartWriter() {
	fmt.Printf("[Writer Goroutine] start...ConnID=%v\n", c.ConnID)
	defer fmt.Printf("[Writer exit] ConnID=%v\n", c.GetRemoteAddr().String())
	// 不断的阻塞等待channel的消息 ，再写给客户端
	for {
		select {
		case data := <-c.msgChan:
			if _, err := c.Conn.Write(data); err != nil {
				fmt.Printf("Send data error:%v\n", err)
				return
			}
		case <-c.exitChan:
			// reader 退出， writer也要退出
			return
		}
	}
}

// 启动链接(让当前链接准备开始工作)
func (c *Connection) Start() {
	fmt.Printf("Connection start()...ConnID=%v\n", c.ConnID)
	go c.StartReader()
	go c.StartWriter()

	// 执行Hock函数
	c.TcpServer.CallOnConnStart(c)
}

// 停止链接
func (c *Connection) Stop() {
	fmt.Println("Conn Stop()...ConnID = ", c.ConnID)
	//如果当前链接已经关闭
	if c.isClosed == true {
		return
	}
	c.isClosed = true

	//如果用户注册了该链接的关闭回调业务，那么在此刻应该显示调用
	c.TcpServer.CallOnConnStop(c)

	// 关闭socket链接
	c.Conn.Close()
	//关闭Writer
	c.exitChan <- true

	//将链接从连接管理器中删除
	c.TcpServer.GetConnManager().RemoveConn(c)

	//关闭该链接全部管道
	close(c.exitChan)
	close(c.msgBuffChan)

}

// 获取当前链接绑定的conn对象
func (c *Connection) GetTCPConnect() *net.TCPConn {
	return c.Conn
}

// 获取链接id
func (c *Connection) GetConnID() uint32 {
	return c.ConnID
}

// 获取客户端链接信息
func (c *Connection) GetRemoteAddr() net.Addr {
	return c.Conn.RemoteAddr()
}

// 发送数据的方法
func (c *Connection) SendMsg(msgId uint32, data []byte) error {
	if c.isClosed {
		return errors.New("Connection closed when send data\n")
	}
	// 将data进行封包
	dp := NewDataPack()
	binaryMsg, err := dp.Pack(NewMessage(msgId, data))
	if err != nil {
		fmt.Printf("pack err msg id =%v\n", msgId)
		return errors.New("pack error msg\n")
	}
	// 将数据发送给writer
	c.msgChan <- binaryMsg

	return nil
}

//添加带缓冲发送消息接口
func (c *Connection) SendBuffMsg(msgId uint32, data []byte) error {
	if c.isClosed == true {
		return errors.New("Connection closed when send buff msg")
	}
	//将data封包，并且发送
	dp := NewDataPack()
	msg, err := dp.Pack(NewMessage(msgId, data))
	if err != nil {
		fmt.Println("Pack error msg id = ", msgId)
		return errors.New("Pack error msg ")
	}

	//写回客户端
	c.msgBuffChan <- msg

	return nil
}

// 设置连接属性
func (c *Connection) SetProperty(key string, value interface{}) {
	c.propertyLock.Lock()
	defer c.propertyLock.Unlock()
	c.property[key] = value
}

// 获取连接属性
func (c *Connection) GetProperty(key string) (interface{}, error) {
	c.propertyLock.RLock()
	defer c.propertyLock.RUnlock()
	if e, ok := c.property[key]; ok {
		return e, nil
	} else {
		return nil, errors.New("GetProperty filed")
	}
}

// 删除连接属性
func (c *Connection) DelProperty(key string) {
	c.propertyLock.Lock()
	defer c.propertyLock.Unlock()
	delete(c.property, key)
}
