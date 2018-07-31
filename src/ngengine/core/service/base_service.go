package service

import (
	"ngengine/common/event"
)

// Service 服务接口，由各个服务自己实现
type Service interface {
	event.Dispatcher
	// 服务前期准备
	Prepare(CoreAPI) error
	// 初始化
	Init(*CoreOption) error
	// 开始运行
	Start() error
	// 服务已经就绪
	Ready()
	// 关闭服务，返回false由服务自行关闭(主动调用Shut函数)
	Close() bool
	// 客户端连接回调
	OnConnect(id uint64)
	// 客户端断开连接
	OnDisconnect(id uint64)
}

type BaseService struct {
	event.EventDispatch
	CoreAPI
}

// Prepare 预处理，在这个回调里进行预加载
func (b *BaseService) Prepare(CoreAPI) error {
	return nil
}

// Init 初始化操作，这里可以获取到服务配置
func (b *BaseService) Init(o *CoreOption) error {
	return nil
}

// Start 服务器启动
func (b *BaseService) Start() error {
	b.CoreAPI.Watch("all")
	return nil
}

// Ready 服务就绪
func (b *BaseService) Ready() {
	b.CoreAPI.SendReady()
}

// Close 服务关闭，如果返回false，则服务自主处理关闭
func (b *BaseService) Close() bool {
	return true
}

// OnConnect 新的客户端连接
func (b *BaseService) OnConnect(id uint64) {

}

// OnDisconnect 客户端断开连接
func (b *BaseService) OnDisconnect(id uint64) {

}
