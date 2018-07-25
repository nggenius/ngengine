package space

import (
	"ngengine/core/service"
	"ngengine/share"
)

type WorldSpaceModule struct {
	service.Module
	spaceManage *SpaceManage
}

func New() *WorldSpaceModule {
	r := &WorldSpaceModule{}
	r.spaceManage = NewSpaceManage(r)
	return r
}

func (r *WorldSpaceModule) Init() bool {
	opt := r.Core.Option()
	rf := opt.Args.String("Region")
	if !r.spaceManage.LoadResource(r.Core.Option().ResRoot + rf) {
		return false
	}

	r.spaceManage.MinRegions = opt.Args.MustInt("MinRegions", 1)
	r.Core.RegisterRemote("Space", NewSpace(r))

	r.Core.Service().AddListener(share.EVENT_READY, r.spaceManage.OnServiceReady)

	return true
}

func (r *WorldSpaceModule) Name() string {
	return "WorldSpace"
}

func (r *WorldSpaceModule) Shut() {
	r.Core.Service().RemoveListener(share.EVENT_READY, r.spaceManage.OnServiceReady)
}
