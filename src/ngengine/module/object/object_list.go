package object

import (
	"container/list"
	"errors"
	"fmt"
)

type ObjectSlot []interface{}
type ObjectList struct {
	slots    []ObjectSlot
	free     *list.List
	slotLens int
	max      int
	usedslot int
	count    int
	maxindex int
}

func NewObjectList(slotlens int, maxitems int) *ObjectList {
	o := &ObjectList{}
	slotcount := maxitems / slotlens
	o.max = maxitems
	o.slots = make([]ObjectSlot, 1, slotcount)
	o.free = list.New()
	o.slotLens = slotlens
	o.slots[0] = make(ObjectSlot, 1, slotlens)
	o.maxindex = 1 // 0号位不用
	return o
}

// 插入对象
func (l *ObjectList) Add(object interface{}) (int, error) {
	if object == nil {
		return -1, errors.New("object is nil")
	}

	if count := l.count + 1; count > l.max {
		return -1, errors.New("object list is full, too much objects")
	}

	if l.free.Len() != 0 {
		e := l.free.Front()
		fi := e.Value.(int)
		slot := fi / l.slotLens
		index := fi % l.slotLens
		l.slots[slot][index] = object
		l.count++
		l.free.Remove(e)
		return fi, nil
	}

	slot := l.maxindex / l.slotLens
	if slot >= len(l.slots) {
		l.slots = append(l.slots, make(ObjectSlot, 0, l.slotLens))
	}
	index := l.maxindex
	l.slots[slot] = append(l.slots[slot], object)
	l.maxindex++
	l.count++
	return index, nil
}

// 移除对象
func (l *ObjectList) Remove(index int, object interface{}) error {
	slot := index / l.slotLens
	slotindex := index % l.slotLens
	if slot >= len(l.slots) || slotindex >= len(l.slots[slot]) {
		return fmt.Errorf("remove object index error, %d", index)
	}

	if l.slots[slot][slotindex] == nil || l.slots[slot][slotindex] != object {
		return fmt.Errorf("remove object not equal")
	}

	l.slots[slot][slotindex] = nil
	l.count--
	l.free.PushBack(index)
	return nil
}

// 获取对象
func (l *ObjectList) Get(index int) (interface{}, error) {
	slot := index / l.slotLens
	slotindex := index % l.slotLens
	if slot >= len(l.slots) || slotindex >= len(l.slots[slot]) {
		return nil, fmt.Errorf("remove object index error, %d", index)
	}

	return l.slots[slot][slotindex], nil
}
