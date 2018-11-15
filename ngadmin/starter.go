package ngadmin

import (
	"github.com/nggenius/ngengine/share"
	"sync/atomic"
)

// StartService 启动进程
func (n *NGAdmin) StartService() {
	adminConfig := n.opts
	appConfig := n.opts.ServicesConfig

	for key := range appConfig.Services {
		for _, v := range appConfig.Services[key] {
			// 分配server_id
			atomic.AddInt32(&adminConfig.SeverCount, 1)
			if adminConfig.SeverCount&share.SID_MAX == 0 {
				// 服务器id上限了
				n.LogFatal("service too much")
				return
			}

			v.AdminId = adminConfig.Id
			v.AdminAddr = adminConfig.Host
			v.AdminPort = adminConfig.Port
			v.ServId = share.ServiceId(adminConfig.SeverCount)

			if v.IsRemoteStart {
				// 运程启动
			} else {
				// 本机启动
				startPath := adminConfig.StartPath(v.ServType)
				start(startPath, v, n.Log)
			}
		}
	}
}
