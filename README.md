# Light TCP Server Framework 
## Overview
It's a simple TCP Server Framework write in Golang.This project is organized by multiple modules.
(ex. server, router, mq, global configuration, work pool and unified message)

## Introduce
模块  | 说明
---- | ----
conf | 存放配置文件模块
setting | 调用配置文件模块, 从这个包中获取配置信息
ziface | 接口（逻辑实现）
znet  | 接口的实现
example | 一些demo实现


### ziface
- iserver: 服务器的接口
- iconnection: 自定义连接所需要实现的方法
- irequest: 封装的请求接口
- irouter: 抽象的路由接口
- imessage: 封装的消息接口
- idatapack: 自定义应用层协议的数据包接口
- imsgHandler: 消息管理模块的接口

### znet
- server: iserver的实现方法
- connection: iconnection的实现
- request: request接口的封装实现
- router: irouter接口的实现
- message: 封装的消息接口实现
- datapack: 数据包封包、拆包接口实现
- msgHandler: 消息管理模块的结构体

### conf
- app.ini 配置文件,存放一些配置信息

### setting
- setting.go: 配置类，存放着从配置文件中获取的数据

### example:
#### v0.1 第一个版本:
基本的go网络编程，实现了一个echo的网络服务器。
日后的升级功能都在此模块上叠加与优化。

- client_demo.go: client的demo实现
- server_demo.go: server的demo实现

#### v0.2 第二个版本(封装连接)：
封装了自己的一个connection模块，实现将业务处理函数与net.Conn进行绑定。
在获取封装的connection时，传入一个自定义处理函数，便能够使得这个连接按照自己所想要的业务进行处理。

#### v0.3 基础路由模块
1. request封装： 将连接和数据绑定在一起
2. 基础路由模块定义
3. 继承路由模块功能

=================================

irouter抽象的router:
- 处理业务之前的方法(相当于钩子)
- 处理业务的主方法
- 处理业务之后的方法

BaseRouter具体的router:
- 处理业务之前的方法
- 处理业务的主方法
- 处理业务之后的方法

路由可以提供一个指令，不同的指令有着不同的处理方式。将这些指令放在一起，就是路由。

tcp为什么需要路由？
不同的消息对应不同的消息处理方式

==================================

加入router模块
- IServer添加路由方法
- Server添加Router成员
- Connection绑定Router成员
- 在Connection中调用 绑定了的Router成员

#### v0.4 JSON形式的全局配置
暂时不用
#### v0.5 Message消息封装
定义一个消息的结构
- 消息id
- 消息长度
- 消息内容

方法
- get
- set

解决消息的粘包问题：经典的TLV序列化

TCP是以流的形式传播的，我们无法知道某个包的长度，接收的TCP包会跟TCP缓冲区大小有关。
粘包：多个包之前无法区分

解决办法：tlv序列化，一种自定义的应用层协议
- 头部：一段固定长度的数字，描述消息id以及消息长度
- 内容：具体内容部分

问题：如何实现自定义的应用层协议的封包以及拆包
- 针对Message进行TLV格式的封装
    - 写message 长度
    - 写message id
    - 写message 内容
- 针对Message进行TVL格式的拆包
    - 先读取固定长度的header
    - 再读取header中内容部分长度的数据

#### v0.6 多路由管理模块
- 消息管理模块，用来管理我们的路由

添加iMsgHandler

将server中的router属性， 变为msgHandler

#### v0.7 读写分离
1. 创建用于reader和writer之间通信的管道
2. 添加一个Goroutine writer用于向客户端写数据 
3. reader由之前直接发送给客户端 改为发送给通信channel
4. 启动reader和writer一同工作

#### v0.8 消息队列&多任务
1. 创建消息队列
```go
var TaskQueue []chan IRequest
```
2. 创建worker池
    - 根据workerPoolSize的数量去创建Worker
    - 每个worker都用一个Go去承载
        - 阻塞，等待与当前worker对应的chan的消息
        - 接收到消息，worker处理当前消息对应的业务，调用DoMsgHandler()
3. 将之前的发送数据，全部改成将消息发送给消息队列和worker池
    - 一个将消息发送给消息队列的方法
    - 需要对消息进行负载均衡(这里使用简单的平均分配)，让哪个woker处理，就发送给对应的queue
    
4. 集成消息队列和工作池到框架

#### v0.9 连接管理机制
用来管理当前所有的连接，例如：增加拒绝连接的方法

linux最大的io连接数 取决于内存大小。

1. 创建连接管理模块
2. 将连接管理集成到框架
3. 给框架提供对外暴露的钩子函数
4. 进行测试开发

连接管理模块：ConnectionManager
- 属性
    - 已经创建的connection集合:map
    - 针对map的保护锁(使用读写锁)
- 方法
    - 添加连接
    - 删除连接
    - 根据连接id 查找对应的连接
    - 总连接个数
    - 清理全部连接：防止程序内存泄漏
    
将连接管理器加入server,每次成功建立连接后，添加到连接管理器中，每次与客户端连接端开后，需要从连接管理器中删除

添加连接创建和销毁的钩子函数，能够在连接到来之后和销毁连接之前做一些事情。
需要在server中添加属性和方法。

#### v0.10 给连接绑定属性
开放一些可以配置连接属性的方法

给connection 添加一个属性集合 和保护这个属性集合的锁

新增的方法，设置连接属性、获取连接属性、删除连接属性


## Application---> MMO GAME
MsgID | Client | Server | 描述
----- | ------ | ------ | ----
1     | -      | SyncPid| 同步玩家本次登录的ID, 说明是server主动发送的
2     | Talk   | -      | 世界聊天
3     | movePackage | -     | 移动消息
200   | -       | Broadcast | 广播消息(To.1: 世界聊天; To.2: 坐标; To.3：动作; To.4: 坐标更新)
201   | -      | SyncPid | 广播消息 掉线
202   | -      | SyncPlayers | 同步周围人的位置信息 


MMOGAME需要实现的一些简单业务：

玩家业务：
- 玩家上线
- 世界聊天
- 上线位置的消息同步
- 移动位置与广播
- 玩家下线

### 1. 多人在线游戏的AOI算法
AOI： area of interest
这是游戏的基础核心，许多逻辑都是因为AOI进出事件驱动的，网络事件也是由于AOI进出事件而产生。

需要为每个玩家设置一个兴趣点，当一个对象状态发生改变时，需要广播给各个玩家，AOI覆盖到的玩家就会接收到消息

功能：
- 当服务器上的玩家或物体状态发生改变时，将消息广播给附近的玩家
- 当玩家进入NPC的警戒区时，AOI模块将消息发送给NPC，NPC再作出相应的AI反映

#### AOI算法实现：
将一个区域划分为多个块，每一个块就是一个AOI兴趣点。

需要两个结构体：
- AOI格子结构体
    - 属性
        - 格子ID
        - 格子左边界坐标
        - 格子右边界坐标
        - 格子上边界坐标
        - 格子下边界坐标
        - 当前格子内的成员集合
    - 方法
        - 初始化格子
        - 添加成员
        - 删除成员
        - 获取成员
        - 打印出格子的基本信息
- 地图结构体

