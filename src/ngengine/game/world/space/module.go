package space

import (
	"ngengine/core/service"
	"ngengine/share"
)

type WorldSpaceModule struct {
	service.Module
	core        service.CoreAPI
	spaceManage *SpaceManage
}

func New() *WorldSpaceModule {
	r := &WorldSpaceModule{}
	r.spaceManage = NewSpaceManage(r)
	return r
}

func (r *WorldSpaceModule) Init(core service.CoreAPI) bool {
	r.core = core
	opt := core.Option()
	rf := opt.Args.String("Region")
	if !r.spaceManage.LoadResource(core.Option().ResRoot + rf) {
		return false
	}

	r.spaceManage.MinRegions = opt.Args.MustInt("MinRegions", 1)
	r.core.RegisterRemote("Space", NewSpace(r))

	r.core.Service().AddListener(share.EVENT_READY, r.spaceManage.OnServiceReady)

	return true
}

func (r *WorldSpaceModule) Name() string {
	return "WorldSpace"
}

func (r *WorldSpaceModule) Shut() {
	r.core.Service().RemoveListener(share.EVENT_READY, r.spaceManage.OnServiceReady)
}
