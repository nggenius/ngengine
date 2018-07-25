package ngadmin

import (
	"encoding/json"
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
	AppConfig map[string][]*ServiceLink
}

// Options 启动配置文件
type Options struct {
	LogFile      string            //日志文件名
	LogLevel     int               //日志等级(DEBUG<INFO<WARN<ERR<FATAL)
	Host         string            //主服务地址
	Port         int               //主服务端口
	LocalAddr    string            //内网通讯地址
	OuterAddr    string            //外网通讯地址
	Master       bool              //是否是主控制器
	Exclusive    bool              //是否独占(不参与负载均衡)
	MinClusters  int               //最小集群数量(启动条件)
	HeartTimeout int               //心跳间隔时长
	SeverCount   int32             //启动服务器个数计数
	AppStartPath map[string]string //启动的路径
	AppConfig    ServiceLinks      //启动其他服的配置
}

// LoadingConfig 加载配置
func (s *Options) Load(appPath string, appParaPath string) error {
	s.AppStartPath = make(map[string]string)
	s.AppConfig = ServiceLinks{AppConfig: make(map[string][]*ServiceLink)}

	err := s.lodgingAppPath(appPath)
	if err != nil {
		return err
	}

	err = s.loadAppArgs(appParaPath)
	if err != nil {
		return err
	}

	return nil
}

// StartPath 获取启动的路径
func (s *Options) StartPath(startName string) string {
	if v, ok := s.AppStartPath[startName]; ok {
		return v
	}

	return ""
}

// LodgingAppPath 加载AppPath配置
func (s *Options) lodgingAppPath(path string) error {
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
		if key == "apps" {
			b, err := pathJSON.MarshalJSON()
			if err != nil {
				return err
			}

			err = json.Unmarshal(b, &s.AppStartPath)
			if err != nil {
				return err
			}

		}
	}

	return nil
}

// LoadAppArgs 添加参数
func (s *Options) loadAppArgs(path string) error {
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
			masterJSON, err := JSON.Get(key).MarshalJSON()
			if err != nil {
				return err
			}

			json.Unmarshal(masterJSON, s)
			JSON.Del(key)
			break
		}
	}

	appConfigJSON, err := JSON.MarshalJSON()
	if err != nil {
		return err
	}
	json.Unmarshal(appConfigJSON, &s.AppConfig.AppConfig)
	return nil
}
