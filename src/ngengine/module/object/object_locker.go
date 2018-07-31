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

	if !o.dummy {
		// 先加入队列
		l := NewLocker(locker, lockindex, isSynclock)
		o.LockerQueue.PushBack(l)
	}

	// 如果队列没人就直接通知上锁成功
	if 1 == o.LockerQueue.Len() {
		o.LockObjSuccess(locker, lockindex, isSynclock)
		return
	}
}

// UnLockObj 解锁
func (o *ObjectWitness) UnLockObj(locker rpc.Mailbox, lockindex uint32, isSynclock bool) {

	// 检查是否是上锁者
	if !o.Islock || locker.ServiceId() != o.locker.Locker.ServiceId() || lockindex != o.locker.LockIndex {
		return
	}
	// 需要操作远程对象
	if o.dummy {
		err := o.RemoteUnLockObj(lockindex)
		if err != nil {
			// 没有发送成功
		}
		return
	}

	// 本地解锁
	o.UnLockObjSuccess(isSynclock)
}

// LockObjSuccess 上锁成功(实际上锁的地方)
func (o *ObjectWitness) LockObjSuccess(locker rpc.Mailbox, lockindex uint32, isSynclock bool) {

	// 这里不管是远程还是本地都要把本地设置成锁定
	o.Islock = true
	o.locker = NewLocker(locker, lockindex, isSynclock)

	// 是否是远程操作上锁
	if isSynclock {
		err := o.RemoteLockObjSuccess(lockindex)
		if err != nil {
			// 远程上锁的mailbox已经没有了,本地把锁解开
			o.UnLockObjSuccess(false)
		}
		return
	}

	// 本地的触发对应回调
	if cb, ok := o.LockCb[lockindex]; ok {
		cb()
	}

	// 调用解锁
	o.UnLockObj(locker, lockindex, isSynclock)
}

// UnLockObjSuccess 解锁成功(实际解锁的地方)
func (o *ObjectWitness) UnLockObjSuccess(isSynclock bool) {
	// 通知远程
	if isSynclock {
		// 这里就算没有通知成功也不用做处理本地的解开就没问题
		o.RemoteUnLockObjSuccess()
	}

	if !o.dummy {
		// 移除队列
		e := o.LockerQueue.Front()
		o.LockerQueue.Remove(e)
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
			o.LockObjSuccess(locker.Locker, locker.LockIndex, locker.IsSyncLock)
			break
		}
		// 如果不是这个就是放入错误，直接干掉
		d := e
		e = e.Next()
		o.LockerQueue.Remove(d)
	}
}
