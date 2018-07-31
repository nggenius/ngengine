package gameobject

import (
	"errors"
)

var (
	ERR_CHILD_FULL    = errors.New("container is full")
	ERR_POS_NOT_EMPTY = errors.New("container pos is not empty")
)

type Container struct {
	Cap    int
	Child  []GameObject `json:"-"`
	Childs int
}

func NewContainer(cap int) *Container {
	c := &Container{}
	if cap == 0 {
		cap = 16
	}
	c.Child = make([]GameObject, 0, cap)
	c.Cap = cap
	return c
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

// AddChild 增加一个对象
func (p *Container) AddChild(pos int, g GameObject) error {
	if p.Childs >= p.Cap {
		return ERR_CHILD_FULL
	}

	if pos > 0 {
		if pos >= len(p.Child) && p.Child[pos] != nil {
			return ERR_POS_NOT_EMPTY
		}
	}

	if pos == 0 {
		pos = p.freeIndex()
		if pos == -1 {
			return ERR_CHILD_FULL
		}
	}

	p.Child[pos] = g
	g.SetContainerPos(pos)
	return nil
}

// RemoveChild 移除一个对象
func (p *Container) RemoveChild(pos int) error {
	return nil
}
