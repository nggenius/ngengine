// Package timeevent 时间事件
// 定时执行所注册的任务
package crontab

import (
	"time"
	"strings"
	"errors"
	"strconv"
	"fmt"
	"regexp"
	"ngengine/core/service"
)

// 检测间隔，先写死
const duration int64 = 60

var (
	// ErrCrontabNotInit 模块没有初始化
	ErrCrontabNotInit = errors.New("CrontabModule is not init")

	// ErrCbNil 回调方法为空
	ErrCbNil = errors.New("CallBack function is nil")

	// ErrArgNil 回调参数为空
	ErrArgNil = errors.New("CallBack args is nil")

	// ErrTimeStrFmtError 时间字符串格式错误
	ErrTimeStrFmtError = errors.New("time string must have five components like * * * * *")
)

// crontab 时间事件结构体
// ticks 时间通道 evts 事件集合
type crontab struct {
	lastTime 	int64
	evts		[]event
}

// CrontabModule 时间事件模块
type CrontabModule struct {
	core service.CoreApi
	crtab *crontab
}

// Name 模块名
func (m *CrontabModule) Name() string {
	return "CrontabModule"
}

// Init 模块初始化
func (m *CrontabModule) Init(core service.CoreApi) bool {
	m.core = core
	m.core.LogInfo("CrontabModule is init")
	m.crtab = new()
	return true
}

// Shut 模块关闭
func (m *CrontabModule) Shut() {

}

// OnUpdate 模块Update
func (m *CrontabModule) OnUpdate(t *service.Time) {
	m.check()
}

// OnMessage 模块消息
func (m *CrontabModule) OnMessage(id int, args ...interface{}) {
}

type callback func(args interface{})

// event 时间
type event struct {
	minute		map[int]struct{}
	hour		map[int]struct{}
	day			map[int]struct{}
	month		map[int]struct{}
	dayofweek	map[int]struct{}
	
	cb			callback
	args		interface{}
}

// tick 时间点机构体
type tick struct {
	minute 		int
	hour		int
	day			int
	month		int
	dayofweek	int
}

var (
	matchSpaces = regexp.MustCompile("\\s+")
	matchN      = regexp.MustCompile("(.*)/(\\d+)")
	matchRange  = regexp.MustCompile("^(\\d+)-(\\d+)$")
)

func new() *crontab {
	timeEvt := &crontab {
		lastTime : time.Now().Unix() - int64(time.Now().Second()),
	}

	return timeEvt
} 

// Check crontab插件的主调用方法
func (m *CrontabModule) check() {
	if m.crtab == nil {
		return 
	}
	now := time.Now().Unix()
	if now - m.crtab.lastTime >= duration {
		m.crtab.checkTriggerEvent(time.Now())
		m.crtab.lastTime = now - int64(time.Now().Second())
	}
}

// RegistEvent crontab插件事件注册接口
// 调用来注册时间事件
func (m *CrontabModule)RegistEvent(timeStr string, cb callback, args interface{}) error {
	if m.crtab == nil {
		return ErrCrontabNotInit
	}
	evt, err := parseEventTime(timeStr)
	if err != nil {
		return err	
	}

	if cb == nil {
		return ErrCbNil
	}

	if args == nil {
		return ErrArgNil
	}

	evt.cb = cb
	evt.args = args
	m.crtab.evts = append(m.crtab.evts, evt)

	return nil
}

// 检查时间判断是否触发事件
func (c * crontab)checkTriggerEvent(t time.Time) {
	tk := timeToTick(t)

	for _, evt := range c.evts {
		if canTriggerEvent(evt, tk) {
			triggerEvent(evt)
		}
	}
}

// 把一个时间字符串格式化到even里t
func parseEventTime(timeStr string) (evt event, err error) {
	timeStr = matchSpaces.ReplaceAllLiteralString(timeStr, " ")
	parts := strings.Split(timeStr, " ")
	if len(parts) != 5 {
		return event{}, ErrTimeStrFmtError
	}

	evt.minute, err = parsePartTimeStr(parts[0], 0, 59)
	if err != nil {
		return evt, err
	}

	evt.hour, err = parsePartTimeStr(parts[1], 0, 23)
	if err != nil {
		return evt, err
	}

	evt.day, err = parsePartTimeStr(parts[2], 1, 31)
	if err != nil {
		return evt, err
	}

	evt.month, err = parsePartTimeStr(parts[3], 1, 12)
	if err != nil {
		return evt, err
	}

	evt.dayofweek, err = parsePartTimeStr(parts[4], 0, 6)
	if err != nil {
		return evt, err
	}

	// 日期和星期之间的冲突解决
	switch {
	case len(evt.day) < 31 && len(evt.dayofweek) == 7: // 但是星期是一周的每一天都设置了执行，但是日期没有设置全部，则按照星期的来
		evt.dayofweek = make(map[int]struct{})
	case len(evt.dayofweek) < 7 && len(evt.day) == 31: // 如果日期设置了每一天都执行，但是星期没有设置全部，则按照日期的来
		evt.day = make(map[int]struct{})
	default:
		// both day and dayOfWeek are * or both are set, use combined
		// i.e. don't do anything here
	}

	return evt, nil
}

func parsePartTimeStr(s string, min, max int) (map[int]struct{}, error) {

	r := make(map[int]struct{}, 0)

	// 匹配星号
	if s == "*" {
		for i := min; i <= max; i++ {
			r[i] = struct{}{}
		}
		return r, nil
	}

	// 匹配 */5 或者 1-25/5
	if matches := matchN.FindStringSubmatch(s); matches != nil {
		localMin := min
		localMax := max
		if matches[1] != "" && matches[1] != "*" {
			if rng := matchRange.FindStringSubmatch(matches[1]); rng != nil {
				localMin, _ = strconv.Atoi(rng[1])
				localMax, _ = strconv.Atoi(rng[2])
				if localMin < min || localMax > max {
					return nil, fmt.Errorf("Out of range for %s in %s. %s must be in range %d-%d", rng[1], s, rng[1], min, max)
				}
			} else {
				return nil, fmt.Errorf("Unable to parse %s part in %s", matches[1], s)
			}
		}
		n, _ := strconv.Atoi(matches[2])
		for i := localMin; i <= localMax; i += n {
			r[i] = struct{}{}
		}
		return r, nil
	}

	// 匹配格式为 1,2,4  或者 1,2,10-15,20,30-45 等
	parts := strings.Split(s, ",")
	for _, x := range parts {
		if rng := matchRange.FindStringSubmatch(x); rng != nil {
			localMin, _ := strconv.Atoi(rng[1])
			localMax, _ := strconv.Atoi(rng[2])
			if localMin < min || localMax > max {
				return nil, fmt.Errorf("Out of range for %s in %s. %s must be in range %d-%d", x, s, x, min, max)
			}
			for i := localMin; i <= localMax; i++ {
				r[i] = struct{}{}
			}
		} else if i, err := strconv.Atoi(x); err == nil {
			if i < min || i > max {
				return nil, fmt.Errorf("Out of range for %d in %s. %d must be in range %d-%d", i, s, i, min, max)
			}
			r[i] = struct{}{}
		} else {
			return nil, fmt.Errorf("Unable to parse %s part in %s", x, s)
		}
	}

	if len(r) == 0 {
		return nil, fmt.Errorf("Unable to parse %s", s)
	}

	return r, nil
}


// 在某个时间点事件是否能触发
func canTriggerEvent(evt event, t tick) bool {
	// 判断分钟是否符合
	if _, ok := evt.minute[t.minute]; !ok {
		return false
	} 

	// 判断小时是否符合
	if _, ok := evt.hour[t.hour]; !ok{
		return false
	}

	// 判断天和星期，只要有一个能取出来就符合触发
	_, dayok := evt.day[t.day]
	_, dayofweekok := evt.dayofweek[t.dayofweek]
	if !dayok && !dayofweekok {
		return false
	}

	// 判断月是否符合
	if _, ok := evt.month[t.month]; !ok {
		return false
	}

	return true
}



func timeToTick(t time.Time) tick{
	return tick{
		minute:		t.Minute(),
		hour: 		t.Hour(),
		day:		t.Day(),
		month:		int(t.Month()),
		dayofweek:	int(t.Weekday()), 
	}
}

// 触发事件
func triggerEvent(evt event) {
	evt.cb(evt.args)
}
