package ngadmin

import (
	"ngengine/share"
	"sync/atomic"
)

func (n *NGAdmin) StartApp() {
	admin_config := n.opts
	appconfig := admin_config.AppConfig

	for key := range appconfig.AppConfig {
		for _, v := range appconfig.AppConfig[key] {
			// 分配server_id
			atomic.AddInt32(&admin_config.SeverCount, 1)
			if admin_config.SeverCount&share.SID_MAX == 0 {
				// 服务器id上限了
				n.LogFatal("severId full")
				return
			}

			v.AdminAddr = admin_config.Host
			v.AdminPort = admin_config.Port
			v.ServId = share.ServiceId(admin_config.SeverCount)

			if v.IsRemoteStart {
				// 运程启动
			} else {
				// 本机启动
				startPath := admin_config.StartPath(v.ServType)
				start(startPath, v, n.Log)
			}
		}
	}
}
