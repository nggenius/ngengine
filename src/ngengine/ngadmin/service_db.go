package ngadmin

import (
	"container/list"
	"fmt"
	"ngengine/protocol"
	. "ngengine/share"
	"sync"
)

const (
	SRV_OPEN = iota
	SRV_CLOSE
)

type ServiceDB struct {
	sync.RWMutex
	ctx        *Context
	serviceMap map[ServiceId]*ServiceInfo
	watcher    map[string]*list.List
}

type PeerInfo struct {
	ServId     ServiceId //服务ID
	ServName   string    //服务名称
	ServType   string    //服务类型
	Status     int       //状态
	RemoteAddr string    //ip地址
	RemotePort int       //端口号
}

type ServiceInfo struct {
	PeerInfo *PeerInfo
	Client   *Client
}

func NewServ(peerinfo *PeerInfo, client *Client) *ServiceInfo {
	s := &ServiceInfo{
		PeerInfo: peerinfo,
		Client:   client,
	}
	return s
}

func NewServiceDB(context *Context) *ServiceDB {
	return &ServiceDB{
		ctx:        context,
		serviceMap: make(map[ServiceId]*ServiceInfo),
		watcher:    make(map[string]*list.List),
	}
}

// 增加一个服务
func (s *ServiceDB) AddService(id ServiceId, service *ServiceInfo) error {
	s.Lock()
	defer s.Unlock()
	srv, dup := s.serviceMap[id]
	if dup {
		if srv.PeerInfo.ServName == service.PeerInfo.ServName { //已经存在
			return fmt.Errorf("service dup, %v", *service.PeerInfo)
		}
		//id不一样,可能是服务重启了
		delete(s.serviceMap, id)
	}

	s.serviceMap[id] = service
	s.ctx.ngadmin.LogInfo("add service:", id, *service.PeerInfo)
	srvs := &protocol.Services{
		OpType:  protocol.ST_ADD,
		All:     false,
		Service: make([]protocol.ServiceInfo, 1),
	}

	srvs.Service[0].Id = service.PeerInfo.ServId
	srvs.Service[0].Name = service.PeerInfo.ServName
	srvs.Service[0].Type = service.PeerInfo.ServType
	srvs.Service[0].Status = service.PeerInfo.Status
	srvs.Service[0].Addr = service.PeerInfo.RemoteAddr
	srvs.Service[0].Port = service.PeerInfo.RemotePort

	for t, l := range s.watcher { //通知其它服务增加
		if t == "all" || t == service.PeerInfo.ServType {
			for ele := l.Front(); ele != nil; ele = ele.Next() {
				if srv, has := s.serviceMap[ele.Value.(ServiceId)]; has {
					if _, err := srv.Client.SendProtocol(protocol.A2S_SERVICES, srvs); err != nil {
						s.ctx.ngadmin.LogErr("sync service failed, ", err)
					}
				}
			}
		}
	}

	return nil
}

// 移除一个服务
func (s *ServiceDB) RemoveService(name string, id ServiceId) bool {
	s.Lock()
	defer s.Unlock()
	if srv, found := s.serviceMap[id]; found {
		if srv.PeerInfo.ServName != name {
			return false
		}
		delete(s.serviceMap, id)

		s.ctx.ngadmin.LogInfo("remove service:", *srv.PeerInfo)
		srvs := &protocol.Services{
			OpType:  protocol.ST_DEL,
			All:     false,
			Service: make([]protocol.ServiceInfo, 1),
		}

		srvs.Service[0].Id = srv.PeerInfo.ServId
		srvs.Service[0].Name = srv.PeerInfo.ServName
		srvs.Service[0].Type = srv.PeerInfo.ServType
		srvs.Service[0].Status = srv.PeerInfo.Status
		srvs.Service[0].Addr = srv.PeerInfo.RemoteAddr
		srvs.Service[0].Port = srv.PeerInfo.RemotePort

		for t, l := range s.watcher {
			ele := l.Front() //首先从关注列表移除
			for ele != nil {
				next := ele.Next()
				if ele.Value.(ServiceId) == id {
					l.Remove(ele)
					ele = next
					continue
				}

				if t == "all" || t == srv.PeerInfo.ServType { //通知其它服务删除
					if srv, has := s.serviceMap[ele.Value.(ServiceId)]; has {
						if _, err := srv.Client.SendProtocol(protocol.A2S_SERVICES, srvs); err != nil {
							s.ctx.ngadmin.LogErr("sync service failed, ", err)
						}
					}
				}
				ele = next
			}
		}

		return true
	}
	return false
}

// 查找服务
func (s *ServiceDB) LookupService(id ServiceId) *ServiceInfo {
	s.RLock()
	defer s.RUnlock()
	if srv, found := s.serviceMap[id]; found {
		return srv
	}
	return nil
}

// 通过类型查找服务信息
func (s *ServiceDB) LookupSrvByType(typ string) []*ServiceInfo {
	s.RLock()
	defer s.RUnlock()
	services := make([]*ServiceInfo, 0, len(s.serviceMap))
	for _, v := range s.serviceMap {
		if typ == "all" || v.PeerInfo.ServType == typ {
			services = append(services, v)
		}
	}
	return services
}

// 关注服务变动
func (s *ServiceDB) Watch(id ServiceId, typ []string) {
	self := s.LookupService(id)
	if self == nil {
		return
	}
	s.Lock()
	defer s.Unlock()
	if len(typ) == 0 {
		return
	}

	for _, v := range typ {
		if v == "all" {
			typ = typ[:1]
			typ[0] = "all"
		}
	}

	for k := range typ {
		t := typ[k] //关注的类型
		find := false
		if w, has := s.watcher[t]; has { //列表存在
			for ele := w.Front(); ele != nil; ele = ele.Next() {
				if ele.Value.(ServiceId) == id {
					find = true
					break
				}
			}
		} else {
			//不存在
			s.watcher[t] = list.New()
		}

		if !find {
			s.watcher[t].PushBack(id)
		}

		srvs := &protocol.Services{
			OpType:  protocol.ST_SYNC,
			All:     false,
			Service: make([]protocol.ServiceInfo, 0, len(s.serviceMap)),
		}

		if t == "all" {
			srvs.All = true
		}

		for _, v := range s.serviceMap {
			if t == "all" || v.PeerInfo.ServType == t {
				si := protocol.ServiceInfo{}
				si.Id = v.PeerInfo.ServId
				si.Name = v.PeerInfo.ServName
				si.Type = v.PeerInfo.ServType
				si.Status = v.PeerInfo.Status
				si.Addr = v.PeerInfo.RemoteAddr
				si.Port = v.PeerInfo.RemotePort
				srvs.Service = append(srvs.Service, si)
			}
		}

		if _, err := self.Client.SendProtocol(protocol.A2S_SERVICES, srvs); err != nil {
			s.ctx.ngadmin.LogErr("sync service failed, ", err)
		}
	}

}

//取消关注
func (s *ServiceDB) Unwatch(id ServiceId) {
	s.Lock()
	defer s.Unlock()

	for _, l := range s.watcher {
		for ele := l.Front(); ele != nil; ele = ele.Next() {
			if ele.Value.(ServiceId) == id {
				l.Remove(ele)
				break
			}
		}
	}
}
