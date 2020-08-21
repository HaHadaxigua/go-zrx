package znet

import "github.com/HaHadaxigua/go-zrx/ziface"

/**
具体的路由

BaseRouter具体的router:
- 处理业务之前的方法
- 处理业务的主方法
- 处理业务之后的方法

实现router时，先嵌入这个BaseRouter, 然后在根据需求再对这个基类中的方法进行重写

*/
type BaseRouter struct {
}

// 处理conn业务之前的钩子方法Hook
func (br *BaseRouter) PreHandle(r ziface.IRequest) {}

// 处理conn业务的主方法Hook
func (br *BaseRouter) Handle(r ziface.IRequest) {}

// 处理conn业务之后的钩子方法Hook
func (br *BaseRouter) AfterHandle(r ziface.IRequest) {}
