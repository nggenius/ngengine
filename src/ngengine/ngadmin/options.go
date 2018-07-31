package ngadmin

import (
	"encoding/json"
	"fmt"
	"ngengine/core/service"

	simplejson "github.com/bitly/go-simplejson"
	"github.com/mysll/toolkit"
)

// ServiceLink 启动其他地方的app
type ServiceLink struct {
	IsRemoteStart bool // 是否是远程启动

	service.CoreOption
}

type ServiceLinks struct {
	Services map[string][]*ServiceLink
}

// Options 启动配置文件
type Options struct {
	Id             int               // id
	LogFile        string            // 日志文件名
	LogLevel       int               // 日志等级(DEBUG<INFO<WARN<ERR<FATAL)
	Host           string            // 主服务地址
	Port           int               // 主服务端口
	LocalAddr      string            // 内网通讯地址
	OuterAddr      string            // 外网通讯地址
	Master         bool              // 是否是主控制器
	Exclusive      bool              // 是否独占(不参与负载均衡)
	HeartTimeout   int               // 心跳间隔时长
	SeverCount     int32             // 启动服务器个数计数
	MinClusters    int               // 最小的admin个数（最小集群数量(启动条件)
	ServicePath    map[string]string // 启动的路径
	MustServices   []string          // 必须要启动的app
	ServicesConfig ServiceLinks      // 启动其他服的配置
	DebugMode      bool              // 调试模式，只启动admin,不启动其它服务
	MessageLog     bool              // 消息日志
}

// Load 加载配置
func (s *Options) Load(servicePath string, serviceDefPath string) error {
	s.ServicePath = make(map[string]string)
	s.ServicesConfig = ServiceLinks{Services: make(map[string][]*ServiceLink)}

	err := s.loadServicePath(servicePath)
	if err != nil {
		return err
	}

	err = s.loadServiceArgs(serviceDefPath)
	if err != nil {
		return err
	}

	// 不能没有id
	if 0 == s.Id {
		return fmt.Errorf("id cont 0")
	}
	return nil
}

// StartPath 获取启动的路径
func (s *Options) StartPath(startName string) string {
	if v, ok := s.ServicePath[startName]; ok {
		return v
	}

	return ""
}

// loadServicePath 加载ServicePath配置
func (s *Options) loadServicePath(path string) error {
	b, err := toolkit.ReadFile(path)
	if err != nil {
		panic(false)
	}

	JSON, err := simplejson.NewJson(b)
	if err != nil {
		return err
	}

	def, err := JSON.Map()
	if err != nil {
		return err
	}
	for key := range def {
		pathJSON := JSON.Get(key)
		if key == "Services" {
			b, err := pathJSON.MarshalJSON()
			if err != nil {
				return err
			}

			err = json.Unmarshal(b, &s.ServicePath)
			if err != nil {
				return err
			}
		} else if key == "MustServices" {
			b, err := pathJSON.MarshalJSON()
			if err != nil {
				return err
			}
			err = json.Unmarshal(b, &s.MustServices)
			if err != nil {
				return err
			}
		}
	}

	return nil
}

// loadServiceArgs 添加参数
func (s *Options) loadServiceArgs(path string) error {
	b, err := toolkit.ReadFile(path)
	if err != nil {
		panic(false)
	}

	JSON, err := simplejson.NewJson(b)
	if err != nil {
		return err
	}

	def, err := JSON.Map()
	if err != nil {
		return err
	}
	for key := range def {
		if key == "Admin" {
			adminJSON, err := JSON.Get(key).MarshalJSON()
			if err != nil {
				return err
			}

			json.Unmarshal(adminJSON, s)
			JSON.Del(key)
			break
		}
	}

	serviceJSON, err := JSON.MarshalJSON()
	if err != nil {
		return err
	}
	json.Unmarshal(serviceJSON, &s.ServicesConfig.Services)
	return nil
}
