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
	Conf *Configuration
)

type Configuration struct{
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
}

func init(){
	var err error
	Cfg, err = ini.Load("conf/app.ini")
	if err!=nil{
		log.Fatalf("Fail to parse 'conf/app.ini':%v", err)
	}
}

func NewConfiguration() *Configuration{
	sec, err := Cfg.GetSection("server")
	if err!=nil{
		log.Fatalf("Fail to get section 'server':%v", err)
	}
	return &Configuration{
		// load basic module
		RunMode: Cfg.Section("").Key("RUN_MODE").String(),
		Version: Cfg.Section("").Key("VERSION").String(),
		// load server module
		Name: sec.Key("NAME").String(),
		IpVersion: sec.Key("IP_VERSION").String(),
		IpAddress: sec.Key("IP_ADDRESS").String(),
		Port: sec.Key("PORT").MustInt(),
		ReadTimeout: time.Duration(sec.Key("READ_TIMEOUT").MustInt(60))*time.Second,
		WriteTimeout: time.Duration(sec.Key("WRITE_TIMEOUT").MustInt(60))*time.Second,
		MaxConn: sec.Key("MAX_CONN").MustInt(),
		MaxPackageSize: sec.Key("MAX_PACKAGE_SIZE").MustInt(),
		WorkerPoolSize: sec.Key("WORKER_POOL_SIZE").MustInt(),
		MaxTaskQueueLen: sec.Key("MAX_TASK_QUEUE_LEN").MustInt(),
	}
}


