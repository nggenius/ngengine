package service

import (
	"time"
)

var (
	BEAT_TIME = time.Millisecond * 100 //100ms
)

// 服务的时间集合
type Time struct {
	time        time.Time     // 启动时间
	beatTime    time.Time     // 心跳时间
	beatTrigger int           // 心跳触发器
	updateTime  time.Time     // 帧更新时间
	deltaTime   time.Duration // 帧间隔时间
	frameCount  int           // 帧数
}

func NewTime() *Time {
	t := &Time{}
	t.time = time.Now()
	t.updateTime = t.time
	t.beatTime = t.time
	return t
}

func (t *Time) FrameCount() int {
	return t.frameCount
}

func (t *Time) DeltaTime() time.Duration {
	return t.deltaTime
}

// 更新所有时间
func (t *Time) Update(now time.Time) {
	t.beatTrigger = 0
	t.deltaTime = now.Sub(t.updateTime)
	t.updateTime = now
	t.frameCount++
	if duration := now.Sub(t.beatTime); duration > BEAT_TIME {
		t.beatTrigger = int(duration / BEAT_TIME)
		t.beatTime = now
	}
}

// 检查心跳
func (t *Time) CheckBeat() int {
	return t.beatTrigger
}

// 获取服务运行时间
func (t *Time) Time() time.Duration {
	return time.Now().Sub(t.time)
}
