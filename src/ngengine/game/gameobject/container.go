package gameobject

import (
	"errors"
	"fmt"
)

var (
	ERR_CHILD_FULL    = errors.New("container is full")
	ERR_POS_NOT_EMPTY = errors.New("container pos is not empty")
	ERR_ADD_FAILED    = errors.New("add failed")
)

type Container struct {
	Cap    int
	Child  []GameObject `json:"-"`
	Childs int
}

func NewContainer(cap int) *Container {
	c := &Container{}
	if cap < 0 {
		panic("cap must >= 0")
	}
	c.Cap = cap
	if cap == 0 {
		cap = 16
	}
	c.Child = make([]GameObject, 0, cap)

	return c
}

// Resize 修改容量，如果新的容量小于原始容量，则容器必须为空
func (p *Container) Resize(newcap int) error {
	if newcap <= 0 {
		panic("cap must above 0")
	}

	if p.Cap == newcap {
		return nil
	}

	if (p.Cap == 0 || newcap < p.Cap) && p.Childs > 0 {
		return fmt.Errorf("container is not empty")
	}

	if cap(p.Child) >= newcap {
		p.Cap = newcap
		return nil
	}

	c := make([]GameObject, 0, newcap)
	copy(c, p.Child)
	p.Child = c
	p.Cap = newcap
	return nil
}

// freeIndex 获取一个空位置
func (p *Container) freeIndex() int {
	if p.Childs < cap(p.Child) {
		for i := range p.Child {
			if p.Child[i] == nil {
				return i
			}
		}
		i := len(p.Child)
		p.Child = append(p.Child, nil)
		return i
	}

	// 空间满了
	if p.Cap != 0 {
		return -1
	}

	// 无限容器
	i := len(p.Child)
	p.Child = append(p.Child, nil)
	return i
}

// Add 在Pos指定的位置，插入一个对象。pos如果小于0，查找一个空位插入。
func (p *Container) Add(pos int, g GameObject) error {
	if pos < 0 || pos >= len(p.Child) {
		return errors.New("index error")
	}
	if p.Child[pos] != nil {
		return errors.New("index error")
	}
	p.Child[pos] = g
	p.Childs++
	return nil
}

// Remove 移除一个指定位置的对象
func (p *Container) Remove(pos int, g GameObject) error {
	if pos < 0 || pos >= len(p.Child) {
		return errors.New("index error")
	}
	if p.Child[pos] == nil {
		return errors.New("index error")
	}

	if p.Child[pos] != g {
		return errors.New("object not equal")
	}

	p.Child[pos] = nil
	p.Childs--
	return nil
}

// ChildAt 获取指定位置的对象
func (p *Container) ChildAt(index int) GameObject {
	if index < 0 || index >= len(p.Child) {
		return nil
	}

	return p.Child[index]
}
