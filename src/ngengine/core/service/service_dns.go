package service

import (
	"ngengine/common/event"
	"ngengine/core/rpc"
	"ngengine/protocol"
	. "ngengine/share"
	"sync"

	"github.com/mysll/toolkit"
)

// 服务的DNS
type SrvDNS struct {
	sync.RWMutex
	ctx         *context
	srvs        []*Srv
	idToIndex   map[ServiceId]int
	nameToIndex map[string]int
}

// 创建一个新的DNS
func NewSrvDNS(ctx *context) *SrvDNS {
	s := &SrvDNS{
		ctx:         ctx,
		srvs:        make([]*Srv, 0, 16),
		idToIndex:   make(map[ServiceId]int),
		nameToIndex: make(map[string]int),
	}
	return s
}

// 更新服务信息(由admin消息驱动)
func (s *SrvDNS) Update(srvs protocol.Services) {
	s.Lock()
	defer s.Unlock()
	switch srvs.OpType {
	case protocol.ST_SYNC:
		s.sync(srvs.All, srvs.Service)
	case protocol.ST_ADD:
		for _, v := range srvs.Service {
			s.addSrv(v)
		}
	case protocol.ST_DEL:
		for _, v := range srvs.Service {
			s.removeSrv(v.Id)
		}
	}
}

// 同步服务器信息
func (s *SrvDNS) sync(all bool, srvs []protocol.ServiceInfo) {
	del := make([]ServiceId, 0, 16)
	for _, v := range s.srvs {
		find := false
		for _, srv := range srvs {
			if v.Id == srv.Id {
				find = true
				break
			}
		}

		if !find {
			del = append(del, v.Id)
		}
	}

	for _, id := range del {
		s.removeSrv(id)
	}

	for _, srv := range srvs {
		_, find := s.idToIndex[srv.Id]
		if !find {
			s.addSrv(srv)
			continue
		}
		s.updateSrv(srv)
	}
}

// 增加一个服务
func (s *SrvDNS) addSrv(srv protocol.ServiceInfo) {
	newsrv := &Srv{
		SrvInfo: SrvInfo{
			Id:        srv.Id,
			Name:      srv.Name,
			Type:      srv.Type,
			Status:    srv.Status,
			Addr:      srv.Addr,
			Port:      srv.Port,
			OuterAddr: srv.OuterAddr,
			OuterPort: srv.OuterPort,
			Load:      srv.Load,
		},
		l: s.ctx.Core.Log,
	}

	mb := rpc.GetServiceMailbox(srv.Id)
	newsrv.mb = &mb
	index := -1
	for k, v := range s.srvs {
		if v == nil {
			s.srvs[k] = newsrv
			index = k
		}
	}
	if index == -1 {
		index = len(s.srvs)
		s.srvs = append(s.srvs, newsrv)
	}

	s.idToIndex[srv.Id] = index
	s.nameToIndex[srv.Name] = index
	s.ctx.Core.LogInfo("add ", newsrv)
	s.ctx.Core.Emitter.Fire(EVENT_READY, event.EventArgs{"id": srv.Id}, true)
}

func (s *SrvDNS) UpdateLoad(load protocol.LoadInfo) {
	s.Lock()
	defer s.Unlock()
	if index, find := s.idToIndex[load.Id]; find {
		s.srvs[index].Load = load.Load
	}
}

// 更新某个服务信息
func (s *SrvDNS) updateSrv(srv protocol.ServiceInfo) {
	if index, find := s.idToIndex[srv.Id]; find {
		old := s.srvs[index]
		if old.Name != srv.Name ||
			old.Type != srv.Type ||
			old.Status != srv.Status ||
			old.Addr != srv.Addr ||
			old.Port != srv.Port ||
			old.OuterAddr != srv.OuterAddr ||
			old.OuterPort != srv.OuterPort ||
			old.Load != srv.Load {
			s.ctx.Core.LogInfo("update service,", *old, srv)
			old.Name = srv.Name
			old.Type = srv.Type
			old.Status = srv.Status
			old.Load = srv.Load
			old.OuterAddr = srv.OuterAddr
			old.OuterPort = srv.OuterPort
			if old.Addr != srv.Addr || old.Port != srv.Port {
				old.Addr = srv.Addr
				old.Port = srv.Port
				old.Close()
			}

		}
	}

}

// 移除一个服务
func (s *SrvDNS) removeSrv(id ServiceId) {
	if index, find := s.idToIndex[id]; find {
		srv := s.srvs[index]
		srv.Close()

		delete(s.idToIndex, id)
		delete(s.nameToIndex, srv.Name)
		s.srvs[index] = nil
		s.ctx.Core.LogInfo("remove ", srv)
		s.ctx.Core.Emitter.Fire(EVENT_LOST, event.EventArgs{"id": srv.Id}, true)
	}
}

// 通过id查找服务
func (s *SrvDNS) Lookup(id ServiceId) *Srv {
	s.RLock()
	defer s.RUnlock()
	if index, find := s.idToIndex[id]; find {
		return s.srvs[index]
	}

	return nil
}

// 通过名称查找服务
func (s *SrvDNS) LookupByName(name string) *Srv {
	s.RLock()
	defer s.RUnlock()
	if index, find := s.nameToIndex[name]; find {
		return s.srvs[index]
	}

	return nil
}

// 获取某个类型的所有服务
func (s *SrvDNS) LookupByType(typ string) []*Srv {
	s.RLock()
	defer s.RUnlock()
	ret := make([]*Srv, 0, 8)
	for _, v := range s.srvs {
		if v.Type == typ {
			ret = append(ret, v)
		}
	}

	return ret
}

// 获取某个类型的一个服务
func (s *SrvDNS) LookupMinLoadByType(typ string) *Srv {
	s.RLock()
	defer s.RUnlock()
	var ret *Srv
	load := int32(0x7FFFFFFF)
	for _, v := range s.srvs {
		if v.Type == typ && v.Load < load {
			ret = v
			load = v.Load
		}
	}

	return ret
}

// 获取某个类型的一个服务
func (s *SrvDNS) LookupOneByType(typ string) *Srv {
	s.RLock()
	defer s.RUnlock()
	var ret *Srv
	for _, v := range s.srvs {
		if v.Type == typ {
			ret = v
			break
		}
	}

	return ret
}

// 随机获取某个类型的一个服务
func (s *SrvDNS) LookupRandByType(typ string) *Srv {
	s.RLock()
	defer s.RUnlock()
	ret := make([]*Srv, 0, 8)
	for _, v := range s.srvs {
		if v.Type == typ {
			ret = append(ret, v)
		}
	}

	if len(ret) > 0 {
		return ret[toolkit.RandRange(0, len(ret))]
	}
	return nil
}

// 通过Mailbox获取服务
func (s *SrvDNS) LookupByMailbox(mb rpc.Mailbox) *Srv {
	id := ServiceId(mb.Sid)
	return s.Lookup(id)
}

func (s *SrvDNS) Process() {
	s.RLock()
	defer s.RUnlock()
	for _, v := range s.srvs {
		if v != nil {
			v.Process()
		}
	}
}
