/**
获取配置文件信息
 */
package setting

import (
	"github.com/go-ini/ini"
	"log"
	"time"
)

var(
	Cfg *ini.File

	RunMode string				// 运行模式 dev | release
	Version string				// 读取版本信息

	Name string					// 服务器名
	IpVersion string			// ip版本
	IpAddress string			// ip地址
	Port int					// ip port
	ReadTimeout time.Duration	// 读超时设置
	WriteTimeout time.Duration  // 写超时设置
	MaxConn int					// 最大连接数
	MaxPackageSize int			// 数据包的最大大小
	WorkerPoolSize int			// 工作池中最多有多少个协程
	MaxTaskQueueLen int			// 任务队列的最大任务数
)

func init(){
	var err error
	Cfg, err = ini.Load("conf/app.ini")
	if err!=nil{
		log.Fatalf("Fail to parse 'conf/app.ini':%v", err)
	}
	LoadBase()
	LoadServer()
}

/**
加载基本模块
 */
func LoadBase(){
	RunMode = Cfg.Section("").Key("RUN_MODE").MustString("dev")
	Version = Cfg.Section("").Key("VERSION").MustString("1.0.0")
}

/**
加载服务模块
 */
func LoadServer(){
	sec, err := Cfg.GetSection("server")
	if err!=nil{
		log.Fatalf("Fail to get section 'server':%v", err)
	}
	Name = sec.Key("NAME").MustString("v1_ZRX")
	IpVersion = sec.Key("IP_VERSION").MustString("tcp4")
	IpAddress = sec.Key("IP_ADDRESS").MustString("127.0.0.1")
	Port = sec.Key("PORT").MustInt()
	ReadTimeout = time.Duration(sec.Key("READ_TIMEOUT").MustInt(60))*time.Second
	WriteTimeout = time.Duration(sec.Key("WRITE_TIMEOUT").MustInt(60))*time.Second
	MaxConn = sec.Key("MAX_CONN").MustInt()
	MaxPackageSize = sec.Key("MAX_PACKAGE_SIZE").MustInt()
	WorkerPoolSize = sec.Key("WORKER_POOL_SIZE").MustInt()
	MaxTaskQueueLen = sec.Key("MAX_TASK_QUEUE_LEN").MustInt()
}
