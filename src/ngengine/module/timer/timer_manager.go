package timer

type TimerCallBack func(id int64, count int, args interface{})

// 定时器类型
const (
	NONE = iota
	REPEAT
	COUNT
)

type TimerTask struct {
	id      int64 // 任务id
	kind    int   // 类型
	delta   int64 // 心跳时间
	count   int   // 次数
	amount  int   // 总数
	args    interface{}
	cb      TimerCallBack
	timer   *TimerQueue
	manager *TimerManager
}

func (t *TimerTask) TaskCallBack(id int64) {
	t.count++
	if t.kind == REPEAT {
		t.timer.Schedule(t.delta, t)
		t.cb(t.id, t.count, t.args)
	} else if t.kind == COUNT {
		if t.count < t.amount {
			t.timer.Schedule(t.delta, t)
			t.cb(t.id, t.count, t.args)
		} else if t.count == t.amount {
			t.cb(t.id, t.count, t.args)
		}
	}
}

type TimerManager struct {
	genID   int64
	taskMap map[int64]*TimerTask
	timer   *TimerQueue
}

func NewManager() *TimerManager {
	manager := &TimerManager{
		genID:   0,
		taskMap: make(map[int64]*TimerTask),
		timer:   New(),
	}
	return manager
}

func (t *TimerManager) Run() {
	t.timer.Run()
}

func (t *TimerManager) GenerateID() int64 {
	t.genID++
	return t.genID
}

func (t *TimerManager) AddTimer(delta int64, args interface{}, cb TimerCallBack) (id int64) {
	task := &TimerTask{
		id:      t.GenerateID(),
		kind:    REPEAT,
		delta:   delta,
		count:   0,
		amount:  0,
		cb:      cb,
		args:    args,
		timer:   t.timer,
		manager: t,
	}
	t.taskMap[task.id] = task
	t.timer.Schedule(delta, task)
	return task.id
}

func (t *TimerManager) AddCountTimer(amount int, delta int64, args interface{}, cb TimerCallBack) (id int64) {
	task := &TimerTask{
		id:      t.GenerateID(),
		kind:    COUNT,
		delta:   delta,
		count:   0,
		amount:  amount,
		cb:      cb,
		args:    args,
		timer:   t.timer,
		manager: t,
	}
	t.taskMap[task.id] = task
	t.timer.Schedule(delta, task)
	return task.id
}

func (t *TimerManager) RemoveTimer(id int64) bool {
	_, ok := t.taskMap[id]
	if ok {
		t.taskMap[id] = nil
		delete(t.taskMap, id)
		return true
	}
	return false
}

func (t *TimerManager) FindTimer(id int64) (bool, int) {
	task, ok := t.taskMap[id]
	if ok {
		return true, task.kind
	}
	return false, NONE
}

func (t *TimerManager) GetTimerDelta(id int64) int64 {
	task, ok := t.taskMap[id]
	if ok {
		return task.delta
	}
	return 0
}
