package ngadmin

import (
	"container/list"
	"fmt"
	"ngengine/protocol"
	. "ngengine/share"
	"sync"
	"time"
)

const (
	SRV_OPEN = iota
	SRV_CLOSE
)

type ServiceDB struct {
	sync.RWMutex
	ctx        *Context
	serviceMap map[ServiceId]*ServiceInfo // 当前已经连接的服务
	watcher    map[string]*list.List      // 服务的关注列表，[服务类型,all关注所有类型]：[关注者id]
	Ready      bool
}

type PeerInfo struct {
	ServId     ServiceId // 服务ID
	ServName   string    // 服务名称
	ServType   string    // 服务类型
	Status     int8      // 状态
	RemoteAddr string    // ip地址
	RemotePort int       // 端口号
	OuterAddr  string    // 外网端口
	OuterPort  int       // 端口号
	Load       int32     // 负载情况
}

type ServiceInfo struct {
	AdminId  int // admin id
	PeerInfo *PeerInfo
	Client   *Client
}

func (s *ServiceInfo) String() string {
	return fmt.Sprintf("Service{Id:%d,Name:%s,Type:%s,Status:%d,Addr:%s,Port:%d,Outer Addr:%s,Outer Port:%d,Load:%d}",
		s.PeerInfo.ServId,
		s.PeerInfo.ServName,
		s.PeerInfo.ServType,
		s.PeerInfo.Status,
		s.PeerInfo.RemoteAddr,
		s.PeerInfo.RemotePort,
		s.PeerInfo.OuterAddr,
		s.PeerInfo.OuterPort,
		s.PeerInfo.Load,
	)
}

// NewServ new serverinfo
func NewServ(adminid int, peerinfo *PeerInfo, client *Client) *ServiceInfo {
	s := &ServiceInfo{
		AdminId:  adminid,
		PeerInfo: peerinfo,
		Client:   client,
	}
	return s
}

// NewServiceDB new service db
func NewServiceDB(context *Context) *ServiceDB {
	return &ServiceDB{
		ctx:        context,
		serviceMap: make(map[ServiceId]*ServiceInfo),
		watcher:    make(map[string]*list.List),
	}
}

// AddService 增加一个服务
func (s *ServiceDB) AddService(id ServiceId, service *ServiceInfo) error {
	s.Lock()
	defer s.Unlock()
	srv, dup := s.serviceMap[id]
	if dup {
		if srv.PeerInfo.ServName == service.PeerInfo.ServName && srv.AdminId == service.AdminId { //已经存在
			return fmt.Errorf("service dup, %v", service)
		}
		//id不一样,可能是服务重启了
		delete(s.serviceMap, id)
	}

	s.serviceMap[id] = service
	s.ctx.ngadmin.LogInfo("add ", service)
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
	srvs.Service[0].OuterAddr = service.PeerInfo.OuterAddr
	srvs.Service[0].OuterPort = service.PeerInfo.OuterPort
	srvs.Service[0].Load = service.PeerInfo.Load

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

// CheckReady 检查关键服务是否就绪了
func (s *ServiceDB) CheckReady(si *ServiceInfo) {
	s.Lock()
	defer s.Unlock()

	if !s.Ready {
		for src_typ, num := range s.ctx.ngadmin.opts.MustServices {
			// 查看是否有这么多个Ready好的服务
			srvList := s.lookupReadySrvByType(src_typ)
			if num != len(srvList) {
				return
			}
		}

		// 给所有的服务发送Ready
		s.Broadcast(protocol.A2S_ALL_READY, nil)
		s.Ready = true

	} else { // 发给给当前服务
		si.Client.SendProtocol(protocol.A2S_ALL_READY, nil)
	}
}

// RemoveService 移除一个服务
func (s *ServiceDB) RemoveService(name string, id ServiceId) bool {
	s.Lock()
	defer s.Unlock()
	if srv, found := s.serviceMap[id]; found {
		if srv.PeerInfo.ServName != name {
			return false
		}
		delete(s.serviceMap, id)

		s.ctx.ngadmin.LogInfo("remove ", srv)
		srvs := &protocol.Services{
			OpType:  protocol.ST_DEL,
			All:     false,
			Service: make([]protocol.ServiceInfo, 1),
		}

		srvs.Service[0].Id = srv.PeerInfo.ServId

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

// LookupService 查找服务
func (s *ServiceDB) LookupService(id ServiceId) *ServiceInfo {
	s.RLock()
	defer s.RUnlock()
	if srv, found := s.serviceMap[id]; found {
		return srv
	}
	return nil
}

// LookupSrvByType 通过类型查找服务信息
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

// lookupReadySrvByType 通过类型查找已准备好的服务信息(无锁的)
func (s *ServiceDB) lookupReadySrvByType(typ string) []*ServiceInfo {
	services := make([]*ServiceInfo, 0, len(s.serviceMap))
	for _, v := range s.serviceMap {
		if (typ == "all" || v.PeerInfo.ServType == typ) && v.PeerInfo.Status == 1 {
			services = append(services, v)
		}
	}
	return services
}

// Watch 关注服务变动
func (s *ServiceDB) Watch(id ServiceId, typ []string) {
	s.Lock()
	defer s.Unlock()
	if len(typ) == 0 {
		return
	}

	var self *ServiceInfo
	found := false
	if self, found = s.serviceMap[id]; !found {
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
				si.OuterAddr = v.PeerInfo.OuterAddr
				si.OuterPort = v.PeerInfo.OuterPort
				si.Load = v.PeerInfo.Load
				srvs.Service = append(srvs.Service, si)
			}
		}

		if _, err := self.Client.SendProtocol(protocol.A2S_SERVICES, srvs); err != nil {
			s.ctx.ngadmin.LogErr("sync service failed, ", err)
		}
	}

}

// UpdateLoad 关注服务变动
func (s *ServiceDB) UpdateLoad(id ServiceId, load int32) {
	s.Lock()
	defer s.Unlock()
	var srv *ServiceInfo
	found := false
	if srv, found = s.serviceMap[id]; !found {
		return
	}
	// 更新负载时间
	srv.PeerInfo.Load = load
	update := &protocol.LoadInfo{
		Id:   id,
		Load: load,
	}

	//通知其它服务更新
	for t, l := range s.watcher {
		for ele := l.Front(); ele != nil; ele = ele.Next() {
			if t == "all" || t == srv.PeerInfo.ServType {
				if srv, has := s.serviceMap[ele.Value.(ServiceId)]; has {
					if _, err := srv.Client.SendProtocol(protocol.A2S_LOAD, update); err != nil {
						s.ctx.ngadmin.LogErr("update service load failed, ", err)
					}
				}
			}
		}

	}
}

// Unwatch 取消关注
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

// CloseAll 关闭所有服务
func (s *ServiceDB) CloseAll() {
	s.Broadcast(protocol.A2S_SERVICE_CLOSE, nil)
}

// Broadcast 广播给所有人
func (s *ServiceDB) Broadcast(msgid uint16, msg interface{}) {
	m, err := PackMessage(msgid, msg)
	if err != nil {
		s.ctx.ngadmin.LogErr("pack message failed, ", err)
		return
	}
	for _, v := range s.serviceMap {
		v.Client.SendMessage(m)
	}
	m.Free()
}

// Done 等待所有服务退出
func (s *ServiceDB) Done() chan struct{} {
	ch := make(chan struct{})
	t := time.NewTicker(time.Second)
	go func() {
		for {
			select {
			case <-t.C:
				if len(s.serviceMap) == 0 {
					t.Stop()
					ch <- struct{}{}
					return
				}
			}
		}
	}()
	return ch
}
