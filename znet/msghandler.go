package znet

import (
	"fmt"
	"github.com/HaHadaxigua/go-zrx/setting"
	"github.com/HaHadaxigua/go-zrx/ziface"
	"strconv"
)

type MsgHandler struct {
	// 存放每个消息Id 所对应的路由
	Apis map[uint32]ziface.IRouter
	// worker池的大小
	WorkerPoolSize uint32
	// 负责Worker取任务的消息队列
	TaskQueue []chan ziface.IRequest
}

func NewMsgHandler() *MsgHandler {
	return &MsgHandler{
		Apis:           make(map[uint32]ziface.IRouter),
		WorkerPoolSize: uint32(setting.WorkerPoolSize),
		TaskQueue:      make([]chan ziface.IRequest, setting.WorkerPoolSize),
	}
}

func (this *MsgHandler) DoMsgHandler(request ziface.IRequest) {
	handler, ok := this.Apis[request.MsgID()]
	if !ok {
		fmt.Printf("api msgID=%v is not registered", request.MsgID())
	}
	// 根据msgID调用对应的router
	handler.PreHandle(request)
	handler.Handle(request)
	handler.AfterHandle(request)
}

func (this *MsgHandler) RegisterRouter(msgID uint32, router ziface.IRouter) {
	// 判断当前msgId 是否已经绑定路由
	if _, ok := this.Apis[msgID]; ok {
		// 如果没有 直接退出
		panic("repeat api, msgID=" + strconv.Itoa(int(msgID)))
		return
	}
	this.Apis[msgID] = router
	fmt.Printf("register router success")
}

/**
启动worker工作器
*/
func (this *MsgHandler) StartWorkerPool() {
	for i := 0; i < setting.WorkerPoolSize; i++ {
		// 每个channel中 最大能够存放的任务数量
		this.TaskQueue[i] = make(chan ziface.IRequest, setting.MaxTaskQueueLen)
		go this.startOneWorker(i, this.TaskQueue[i])
	}
}

/**
worker中处理具体任务
*/
func (this *MsgHandler) startOneWorker(id int, taskQueue chan ziface.IRequest) {
	fmt.Printf("[Start one worker], Id =%v is start...\n", id)
	/**
	死循环，等待消息队列中的消息
	如果有消息到来，他就是client发送的一个request,直接调用request绑定的业务
	*/
	for {
		select {
		case request := <-taskQueue:
			this.DoMsgHandler(request)
		}
	}
}

/**
将消息交给taskQueue,由worker处理
*/
func (this *MsgHandler) SendMsgToTaskQueue(request ziface.IRequest) {
	// 1. 将消息平均分配给不同的worker
	// 如何分配？
	workerID := request.Connection().GetConnID() % this.WorkerPoolSize
	fmt.Printf("Add ConnID =%v, request MsgID=%v, toWorkerID=%v\n",
		request.Connection().GetConnID(),
		request.MsgID(),
		workerID)
	// 2. 将消息发送给对应worker的taskQueue
	this.TaskQueue[workerID] <- request
}
