package space

import (
	"ngengine/common/event"
	"ngengine/core/service"
	"ngengine/share"
	"time"
)

type WorldSpaceModule struct {
	service.Module
	sm *SpaceManage
	rl *event.EventListener
}

func New() *WorldSpaceModule {
	w := &WorldSpaceModule{}
	w.sm = NewSpaceManage(w)
	return w
}

func (w *WorldSpaceModule) Init() bool {
	opt := w.Core.Option()
	rf := opt.Args.String("Region")
	if !w.sm.LoadResource(w.Core.Option().ResRoot + rf) {
		return false
	}

	w.sm.MinRegions = opt.Args.MustInt("MinRegions", 1)
	w.Core.RegisterRemote("Space", NewSpace(w))

	w.rl = w.Core.Service().AddListener(share.EVENT_MUST_SERVICE_READY, w.sm.OnServiceReady)

	w.AddPeriod(time.Second)
	w.AddCallback(time.Second, w.PerSecondCheck)
	return true
}

func (w *WorldSpaceModule) Name() string {
	return "WorldSpace"
}

func (w *WorldSpaceModule) PerSecondCheck(t time.Duration) {
	w.sm.fsm.Dispatch(ETIMER, nil)
}

func (w *WorldSpaceModule) OnUpdate(t *service.Time) {
	w.Module.Update(t)
}

func (w *WorldSpaceModule) Shut() {
	if w.rl != nil {
		w.Core.Service().RemoveListener(share.EVENT_MUST_SERVICE_READY, w.rl)
	}
}
