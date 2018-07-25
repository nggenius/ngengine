package space

import (
	"encoding/json"
	"fmt"
	"ngengine/common/event"
	"ngengine/core/rpc"
	"ngengine/share"
	"ngengine/utils"

	"github.com/mysll/toolkit"
)

const (
	REGION_CREATING = iota + 1
	REGION_RUNNING
	REGION_CLOSING
	REGION_DELETING
	REGION_FAILED
)

type RegionInfo struct {
	share.Region
	Where   rpc.Mailbox
	Dest    rpc.Mailbox
	Players int
	Status  int
}

const (
	RS_NONE = iota
	RS_QUERY
	RS_RUNNING
	RS_OFFLINE
)

type RegionState struct {
	mailbox rpc.Mailbox
	regions []int
	players int
	state   int
}

func NewRegionState(mb rpc.Mailbox) *RegionState {
	s := new(RegionState)
	s.mailbox = mb
	s.regions = make([]int, 0, 10)
	return s
}

// 负载量，每运行一个场景（即使没有玩家)折算成10个玩家+玩家总数
func (r RegionState) Capacity() int {
	return len(r.regions)*10 + r.players
}

func (r RegionState) HasRegion(id int) bool {
	for k := range r.regions {
		if r.regions[k] == id {
			return true
		}
	}
	return false
}

func (r *RegionState) AddRegion(id int) {
	if r.HasRegion(id) {
		return
	}

	r.regions = append(r.regions, id)
}

func (r *RegionState) RemoveRegion(id int) {
	for k := range r.regions {
		if r.regions[k] == id {
			copy(r.regions[k:], r.regions[k+1:])
			r.regions = r.regions[:len(r.regions)-1]
		}
	}
}

type SpaceManage struct {
	ctx         *WorldSpaceModule
	MinRegions  int
	regiondef   map[int]share.Region
	regionmap   map[int]*RegionInfo
	regionstate []*RegionState
}

func NewSpaceManage(ctx *WorldSpaceModule) *SpaceManage {
	s := new(SpaceManage)
	s.ctx = ctx
	s.regionmap = make(map[int]*RegionInfo)
	s.regiondef = make(map[int]share.Region)
	s.regionstate = make([]*RegionState, 0, 10)
	return s
}

func (s *SpaceManage) RegionState(id share.ServiceId) *RegionState {
	for k := range s.regionstate {
		if s.regionstate[k].mailbox.ServiceId() == id {
			return s.regionstate[k]
		}
	}

	return nil
}

func (s *SpaceManage) AddRegion(rs *RegionState) {
	s.regionstate = append(s.regionstate, rs)
}

func (s *SpaceManage) OnServiceReady(e string, args ...interface{}) {
	id := args[0].(event.EventArgs)["id"].(share.ServiceId)
	srv := s.ctx.Core.LookupService(id)
	if srv == nil {
		panic("service not found")
	}

	if srv.Type == "region" {
		rs := s.RegionState(id)
		if rs == nil {
			rs = NewRegionState(*srv.Mailbox())
			s.AddRegion(rs)
		}

		rs.state = RS_QUERY

		s.ctx.Core.MailtoAndCallback(nil, srv.Mailbox(), "Region.Query", s.OnRegionQuery)
	}
}

func (s *SpaceManage) OnRegionQuery(e *rpc.Error, ar *utils.LoadArchive) {
	var id share.ServiceId
	err := ar.Read(&id)
	if err != nil {
		s.ctx.Core.LogErr("read id error", err)
		return
	}

	rs := s.RegionState(id)
	if rs == nil {
		s.ctx.Core.LogWarn("region state not found")
		return
	}

	//TODO: 这里需要同步原来服务器的信息，主要是world异常关闭后进行重建
	rs.state = RS_RUNNING
	s.CreateRegion(1)
}

// 通过ID查找场景
func (s *SpaceManage) FindRegionById(id int) *RegionInfo {
	if r, has := s.regionmap[id]; has {
		return r
	}

	return nil
}

func (s *SpaceManage) findLowerLoadRegion() *RegionState {
	if len(s.regionstate) == 0 {
		return nil
	}
	low := s.regionstate[0].Capacity()
	rs := s.regionstate[0]
	for _, r := range s.regionstate {
		if r.Capacity() < low {
			rs = r
			low = r.Capacity()
		}
	}

	return rs
}

func (s *SpaceManage) CreateRegion(id int) error {
	if _, has := s.regionmap[id]; has {
		return fmt.Errorf("region already created")
	}

	def, has := s.regiondef[id]
	if !has {
		return fmt.Errorf("region def not find")
	}

	rs := s.findLowerLoadRegion()
	if rs == nil {
		return fmt.Errorf("region not found")
	}

	var r RegionInfo
	r.Id = id
	r.Region = def
	r.Status = REGION_CREATING
	r.Where = rs.mailbox

	s.regionmap[id] = &r
	rs.AddRegion(id)

	return s.ctx.Core.MailtoAndCallback(nil, &rs.mailbox, "Region.Create", s.OnCreateRegion, r.Region)
}

func (s *SpaceManage) OnCreateRegion(e *rpc.Error, ar *utils.LoadArchive) {
	var id int
	err := ar.Read(&id)
	if err != nil {
		s.ctx.Core.LogErr("get id error")
		return
	}

	ri := s.FindRegionById(id)
	if ri == nil {
		s.ctx.Core.LogErr("region not found")
		return
	}

	if e != nil {
		ri.Status = REGION_FAILED
		err, _ := ar.ReadString()
		s.ctx.Core.LogErr("region create failed", err)
		return
	}

	var mb rpc.Mailbox
	err = ar.Read(&mb)
	if err != nil {
		s.ctx.Core.LogErr("get mailbox error")
		return
	}
	ri.Dest = mb
	ri.Status = REGION_RUNNING

	s.ctx.Core.Mailto(nil, &mb, "GameScene.Test", "test")
	s.ctx.Core.LogInfo("region created,", ri)
}

func (s *SpaceManage) LoadResource(f string) bool {
	data, err := toolkit.ReadFile(f)
	if err != nil {
		return false
	}

	regions := make(map[string][]share.Region)
	err = json.Unmarshal(data, &regions)
	if err != nil {
		panic(err)
	}

	if r, ok := regions["Regions"]; ok {
		for k := range r {
			s.regiondef[r[k].Id] = r[k]
		}
		return true
	}

	return false
}
