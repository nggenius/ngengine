package service

import (
	"ngengine/common/event"
	"ngengine/core/rpc"
	"ngengine/share"
	"time"
)

var (
	warninglvl = 10 * time.Millisecond
)

//进程rpc处理
func (c *Core) doRPCProcess(ch chan *rpc.RpcCall) {
	var start_time time.Time
	var delay time.Duration
	for {
		select {
		case call := <-ch:
			if call.IsThreadWork() {
				c.busy = true
			} else {
				c.LogDebug(call.GetSrc(), " call ", call.GetMethod())
				start_time = time.Now()
				err := call.Call()
				if err != nil {
					c.LogErr(err)
				}
				delay = time.Now().Sub(start_time)
				if delay > warninglvl {
					c.LogWarn("call ", call.GetMethod(), " delay:", delay.Nanoseconds()/1000000, "ms")
				}
				err = call.Done()
				if err != nil {
					c.LogErr(err)
				}
				call.Free()
				c.busy = true
			}

		default:
			return
		}
	}
}

//RpcResponseProcess rpc回调处理
func (c *Core) doRPCResponseProcess() {
	c.dns.Process()
}

//DoEvent 事件执行
func (c *Core) doEvent(e *event.Event) {

	switch e.Typ {
	case share.EVENT_USER_CONNECT:
		id := e.Args["id"].(uint64)
		c.service.OnConnect(id)
	case share.EVENT_USER_LOST:
		id := e.Args["id"].(uint64)
		c.service.OnDisconnect(id)
		c.clientDB.RemoveClient(id)
	case share.EVENT_SHUTDOWN:
		c.Close()
	}

	// 对消息进行分发
	c.service.DispatchEvent(e.Typ, e.Args)
}

//EventProcess 事件遍历
func (c *Core) eventProcess(e *event.AsyncEvent) {
	var start_time time.Time
	var delay time.Duration
	for {
		evt := e.Capture()
		if evt == nil {
			break
		}
		start_time = c.time.updateTime
		c.doEvent(evt)
		delay = time.Now().Sub(start_time)
		if delay > warninglvl {
			c.LogWarn("DoEvent delay:", delay.Nanoseconds()/1000000, "ms")
		}
		c.busy = true
		evt.Free()
	}

}

// 主循环
func (c *Core) run() {
	c.time = NewTime()
	var now time.Time
	for c.closeState != CS_SHUT {
		c.busy = false
		now = time.Now()
		c.time.Update(now)
		// 处理事件
		c.eventProcess(c.Emitter)
		// 处理rpc
		c.doRPCProcess(c.rpcCh)
		// 处理rpc响应
		c.doRPCResponseProcess()

		if c.time.CheckBeat() > 0 { //大于零，表示过了几个心跳周期
			//c.LogDebug("beat")
		}

		// 运行模块
		for _, m := range c.modules.modules {
			m.OnUpdate(c.time)
		}

		if !c.busy {
			time.Sleep(time.Millisecond)
		}

		if c.closeState == CS_QUIT {
			select {
			case <-c.coreClose:
				c.release()
				c.closeState = CS_SHUT
			default:
				break
			}
		}
	}

	c.LogInfo("core is shutdown")
	close(c.coreQuit)
}
