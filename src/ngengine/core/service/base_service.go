package service

import (
	"ngengine/common/event"
)

// 服务接口，由各个服务自己实现
type Service interface {
	// 服务前期准备
	Prepare(CoreApi) error
	// 初始化
	Init(*CoreOption) error
	// 开始运行
	Start() error
	// 服务已经就绪
	Ready()
	// 关闭服务，返回false由服务自行关闭(主动调用Shut函数)
	Close() bool
	// 收到事件
	OnEvent(string, event.EventArgs)
	// 客户端连接回调
	OnConnect(id uint64)
	// 客户端断开连接
	OnDisconnect(id uint64)
}

type BaseService struct {
	CoreApi
}

func (b *BaseService) Prepare(CoreApi) error {
	return nil
}

func (b *BaseService) Init(o *CoreOption) error {
	return nil
}

func (b *BaseService) Start() error {
	b.CoreApi.Watch("all")
	return nil
}

func (b *BaseService) Ready() {

}

func (b *BaseService) Close() bool {
	return true
}

func (b *BaseService) OnEvent(string, event.EventArgs) {

}

func (b *BaseService) OnConnect(id uint64) {

}

func (b *BaseService) OnDisconnect(id uint64) {

}
