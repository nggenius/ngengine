package object

import (
	"ngengine/core/rpc"
)

// Locker 是否加锁
type Locker struct {
	Locker     rpc.Mailbox // 上锁的人
	LockIndex  uint32      // 上锁的索引
	IsSyncLock bool        // 是否是异步锁
}

// NewLocker 初始化一个上锁对象
func NewLocker(obj rpc.Mailbox, lockIndex uint32, isSyncLock bool) *Locker {
	l := &Locker{
		Locker:     obj,
		LockIndex:  lockIndex,
		IsSyncLock: isSyncLock,
	}
	return l
}

// LockObj 对对象上锁
func (o *ObjectWitness) LockObj(callback LockCallBack) {

	// 先加入回调
	o.LockCount++
	if nil != callback {
		o.LockCb[o.LockCount] = callback
	}

	// 需要操作远程对象
	if o.dummy {
		o.RemoteLockObj(o.LockCount)
		return
	}

	// 本地上锁
	o.AddLocker(o.objid, o.LockCount, false)
}

// AddLocker 增加上锁对象
func (o *ObjectWitness) AddLocker(locker rpc.Mailbox, lockindex uint32, isSynclock bool) {

	// 如果队列没人就直接通知上锁成功
	if 0 == o.LockerQueue.Len() {
		o.LockObjSuccess(locker, lockindex, isSynclock)
		return
	}

	// 如果正在上锁中就加入队列
	l := NewLocker(locker, lockindex, isSynclock)
	o.LockerQueue.PushBack(l)
}

// UnLockObj 解锁
func (o *ObjectWitness) UnLockObj(locker rpc.Mailbox, lockindex uint32, isSynclock bool) {

	// 检查是否是上锁者
	if !o.Islock || locker.ServiceId() != o.locker.Locker.ServiceId() || lockindex != o.locker.LockIndex {
		return
	}

	// 需要操作远程对象
	if o.dummy {
		o.RemoteUnLockObj(lockindex)
		return
	}

	// 本地解锁
	o.UnLockObjSuccess(isSynclock)
}

// LockObjSuccess 上锁成功
func (o *ObjectWitness) LockObjSuccess(locker rpc.Mailbox, lockindex uint32, isSynclock bool) {

	// 这里不管是远程还是本地都要把本地设置成锁定
	o.Islock = true
	o.locker = NewLocker(locker, lockindex, isSynclock)

	// 是否是远程操作上锁
	if isSynclock {
		o.RemoteLockObjSuccess(lockindex)
		return
	}

	// 本地的触发对应回调
	if cb, ok := o.LockCb[lockindex]; ok {
		cb()
	}

	// 调用解锁
	o.UnLockObj(locker, lockindex, isSynclock)
}

// UnLockObjSuccess 解锁成功
func (o *ObjectWitness) UnLockObjSuccess(isSynclock bool) {
	// 通知远程
	if isSynclock {
		o.RemoteUnLockObjSuccess()
	}

	o.Islock = false
	o.locker = nil
	// 对象是本地的查看列表还有没有任务
	if !o.dummy {
		o.ExecuteNextLock()
	}
}

// ExecuteNextLock 执行下一个锁请求
func (o *ObjectWitness) ExecuteNextLock() {
	if o.Islock || 0 == o.LockerQueue.Len() {
		return
	}

	for e := o.LockerQueue.Front(); e != nil; {
		if locker, ok := e.Value.(*Locker); ok {
			o.LockerQueue.Remove(e)
			o.LockObjSuccess(locker.Locker, locker.LockIndex, locker.IsSyncLock)
			break
		}

		d := e
		e = e.Next()
		o.LockerQueue.Remove(d)
	}
}
