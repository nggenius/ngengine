package ngadmin

import (
	"encoding/json"
	"io/ioutil"
	"ngengine/core/service"
	"os"
)

type ServiceLink struct {
	ServType string //服务类型
	ExeName  string //启动文件名
}

// 启动配置文件
type Options struct {
	LogFile      string                 // 日志文件名
	LogLevel     int                    // 日志等级(DEBUG<INFO<WARN<ERR<FATAL)
	Host         string                 // 主服务地址
	Port         int                    // 主服务端口
	LocalAddr    string                 // 内网通讯地址
	OuterAddr    string                 // 外网通讯地址
	Master       bool                   // 是否是主控制器
	Exclusive    bool                   // 是否独占(不参与负载均衡)
	MinClusters  int                    // 最小集群数量(启动条件)
	HeartTimeout int                    // 心跳间隔时长
	ServiceDefs  map[string]ServiceLink // 服务定义
	Services     []service.CoreOption   // 要启动的服务
}

//从文件中加载配置
func (p *Options) LoadFromFile(file string) error {
	f, err := os.Open(file)
	if err != nil {
		return err
	}
	d, err := ioutil.ReadAll(f)
	if err != nil {
		return err
	}

	return json.Unmarshal(d, p)
}
