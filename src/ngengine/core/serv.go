package core

import (
	"fmt"
	"math/rand"
	"ngengine/core/service"
	"sort"
	"time"

	"github.com/mysll/toolkit"
)

var (
	services = make(map[string]service.Service)
	srvinst  = make(map[string]*service.Core)
	wg       = toolkit.WaitGroupWrapper{}
)

// 注册服务
func RegisterService(name string, service service.Service) {
	if service == nil {
		panic("serv: Register service is nil" + name)
	}

	if _, dup := services[name]; dup {
		panic("serv: Register twice for service" + name)
	}

	services[name] = service
}

func CreateService(name string, args string) (*service.Core, error) {
	if _, dup := srvinst[name]; dup {
		return nil, fmt.Errorf("service (%s) already created", name)
	}

	var srv service.Service
	has := false
	if srv, has = services[name]; !has {
		return nil, fmt.Errorf("service (%s) not found", name)
	}
	_service := service.CreateService(srv)
	if err := _service.Init(args); err != nil {
		return nil, err
	}
	srvinst[name] = _service
	return _service, nil
}

func RunAllService() {
	for k := range srvinst {
		srv := srvinst[k]
		wg.Wrap(func() {
			srv.Serv()
		})
	}
}

func CloseAllService() {
	for _, srv := range srvinst {
		srv.Close()
	}
}

func Wait() {
	wg.Wait()
}

// 返回所有服务名字按字符顺序排序
func Services() []string {
	var lists []string
	for name := range services {
		lists = append(lists, name)
	}
	sort.Strings(lists)
	return lists
}

func init() {
	rand.Seed(time.Now().Unix())
}
