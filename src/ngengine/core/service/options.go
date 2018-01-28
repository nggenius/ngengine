package service

import (
	"encoding/json"
	"ngengine/share"
)

// 服务配置选项
type CoreOption struct {
	ServId     share.ServiceId   //服务id
	AdminAddr  string            //管理地址
	AdminPort  int               //管理端口
	ServType   string            //服务类型
	ServName   string            //服务名称
	ServAddr   string            //服务内部地址(内部通讯用)
	ServPort   int               //服务内部端口号
	Expose     bool              //是否启动外网连接
	HostAddr   string            //外网连接地址
	HostPort   int               //外网连接端口号
	LogFile    string            //日志文件
	LogLevel   int               //日志等级
	Args       map[string]string //额外的启动参数
	MaxRpcCall int               //rpc缓冲区大小
}

// 从json文本中加载配置
func ParseOption(args string) (*CoreOption, error) {
	opt := &CoreOption{}
	if err := json.Unmarshal([]byte(args), opt); err != nil {
		return nil, err
	}
	if opt.MaxRpcCall <= 0 {
		opt.MaxRpcCall = 1024
	}
	return opt, nil
}
