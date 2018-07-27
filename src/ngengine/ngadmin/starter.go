package ngadmin

import (
	"ngengine/share"
	"sync/atomic"
)

func (n *NGAdmin) StartApp() {
	adminConfig := n.opts
	appConfig := n.opts.AppConfig

	for key := range appConfig.AppConfig {
		for _, v := range appConfig.AppConfig[key] {
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
