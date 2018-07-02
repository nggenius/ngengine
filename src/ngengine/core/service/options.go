package service

import (
	"encoding/json"
	"ngengine/share"
)

type Args map[string]interface{}

func (a Args) Has(arg string) bool {
	_, ok := a[arg]
	return ok
}

func (a Args) Int(arg string) int {
	if a, ok := a[arg]; ok {
		if val, ok := a.(int); ok {
			return val
		}
	}

	return 0
}

func (a Args) Bool(arg string) bool {
	if a, ok := a[arg]; ok {
		if val, ok := a.(bool); ok {
			return val
		}
	}

	return false
}

func (a Args) String(arg string) string {
	if a, ok := a[arg]; ok {
		if val, ok := a.(string); ok {
			return val
		}
	}

	return ""
}

func (a Args) Float64(arg string) float64 {
	if a, ok := a[arg]; ok {
		if val, ok := a.(float64); ok {
			return val
		}
	}

	return 0.0
}

// 服务配置选项
type CoreOption struct {
	ServId     share.ServiceId //服务id
	AdminAddr  string          //管理地址
	AdminPort  int             //管理端口
	ServType   string          //服务类型
	ServName   string          //服务名称
	ServAddr   string          //服务内部地址(内部通讯用)
	ServPort   int             //服务内部端口号
	Expose     bool            //是否启动外网连接
	OuterAddr  string          //外网连接地址
	HostAddr   string          //外网连接的监听地址
	HostPort   int             //外网连接端口号
	LogFile    string          //日志文件
	LogLevel   int             //日志等级
	Args       Args            //额外的启动参数
	MaxRpcCall int             //rpc缓冲区大小
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
