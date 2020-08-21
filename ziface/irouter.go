package ziface
/**
irouter抽象的router:
- 处理业务之前的方法(相当于钩子)
- 处理业务的主方法
- 处理业务之后的方法


路由可以提供一个指令，不同的指令有着不同的处理方式。将这些指令放在一起，就是路由。

tcp为什么需要路由？
不同的消息对应不同的消息处理方式

路由中的数据都是IRequest

服务端应用可以给框架配置当前链接的处理业务方法
 */
type IRouter interface {
	// 处理conn业务之前的钩子方法Hook
	PreHandle(r IRequest)
	// 处理conn业务的主方法Hook
	Handle(r IRequest)
	// 处理conn业务之后的钩子方法Hook
	AfterHandle(r IRequest)
}
