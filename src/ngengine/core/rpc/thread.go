package rpc

import (
	"fmt"
	"ngengine/logger"
	"time"

	"github.com/mysll/toolkit"
)

type Threader interface {
	NewJob(*RpcCall) bool
}

type Thread struct {
	TAG        string
	NumProcess int32
	Queue      []chan *RpcCall
	Quit       bool
	Pools      int
	wg         toolkit.WaitGroupWrapper
	run        bool
	log        *logger.Log
}

func NewThread(tag string, pools int, queuelen int, log *logger.Log) *Thread {
	if pools < 1 || queuelen < 2 {
		return nil
	}
	t := &Thread{}
	t.TAG = tag
	t.Pools = pools
	t.Queue = make([]chan *RpcCall, pools)
	for i := 0; i < pools; i++ {
		t.Queue[i] = make(chan *RpcCall, queuelen)
	}
	t.log = log
	return t
}

func (t *Thread) Start() error {
	if t.run {
		return fmt.Errorf(t.TAG, " thread already run")
	}
	t.log.LogInfo(t.TAG, " start thread, total:", t.Pools)
	for i := 0; i < t.Pools; i++ {
		id := i
		t.wg.Wrap(func() { t.work(id) })
	}

	t.run = true
	return nil
}

func (t *Thread) Wait() {
	t.wg.Wait()
}

func (t *Thread) NewJob(r *RpcCall) bool {
	t.Queue[int(r.GetSrc().Uid)%t.Pools] <- r
	return true
}

func (t *Thread) work(id int) {
	t.log.LogInfo(t.TAG, " thread work, id:", id)
	var start_time time.Time
	var delay time.Duration
	warninglvl := 50 * time.Millisecond
	for {
		select {
		case rpc := <-t.Queue[id]:
			t.log.LogInfo(t.TAG, " thread:", id, rpc.GetSrc(), " call:", rpc.GetMethod())
			start_time = time.Now()
			err := rpc.Call()
			if err != nil {
				t.log.LogErr("rpc error:", err)
			}
			delay = time.Now().Sub(start_time)
			if delay > warninglvl {
				t.log.LogWarn("rpc call ", rpc.GetMethod(), " delay:", delay.Nanoseconds()/1000000, "ms")
			}
			err = rpc.Done()
			if err != nil {
				t.log.LogErr("rpc error:", err)
			}
			rpc.Free()
			break
		default:
			if t.Quit {
				t.log.LogInfo(t.TAG, " thread ", id, " quit")
				return
			}
			time.Sleep(time.Millisecond)
		}
	}
}
