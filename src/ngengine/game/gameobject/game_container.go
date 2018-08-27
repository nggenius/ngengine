package gameobject

import (
	"fmt"
)

// SetCap 设置容量
func (b *BaseObject) SetCap(cap int) error {
	if b.c == nil {
		b.c = NewContainer(cap)
		return nil
	}

	return b.c.Resize(cap)
}

// Cap 获取容量
func (b *BaseObject) Cap() int {
	if b.c == nil {
		return 0
	}
	return b.c.Cap
}

// CanAdd 是否可以增加子对象
func (b *BaseObject) CanAdd(pos int, g GameObject) bool {
	return true
}

// AddChild 增加一个对象
func (b *BaseObject) AddChild(pos int, g GameObject) (int, error) {
	if b.c.Childs >= b.c.Cap {
		return -1, ERR_CHILD_FULL
	}

	if pos >= 0 {
		if pos >= len(b.c.Child) && b.c.Child[pos] != nil {
			return -1, ERR_POS_NOT_EMPTY
		}
	}

	if pos < 0 {
		pos = b.c.freeIndex()
		if pos == -1 {
			return -1, ERR_CHILD_FULL
		}
	}

	if !b.CanAdd(pos, g) {
		return -1, ERR_ADD_FAILED
	}

	if err := b.c.Add(pos, g); err != nil {
		return -1, err
	}

	g.SetParent(b.gameObject)
	g.SetParentIndex(pos)
	return pos, nil
}

// RemoveChild 移除一个对象
func (b *BaseObject) RemoveChild(pos int, g GameObject) error {

	if g.Parent() != b.gameObject {
		return fmt.Errorf("parent not equal")
	}

	if g.ParentIndex() != pos {
		return fmt.Errorf("container pos not equal")
	}

	// TODO: 事件回调 移除前
	if err := b.c.Remove(pos, g); err != nil {
		return err
	}

	g.SetParent(nil)
	g.SetParentIndex(-1)
	// TODO: 事件回调 移除后

	return nil
}

// ChildAt 取子对象
func (b *BaseObject) ChildAt(pos int) GameObject {
	return b.c.ChildAt(pos)
}

// ChildAtIf 取子对象接口
func (b *BaseObject) ChildAtIf(pos int) interface{} {
	return b.c.ChildAt(pos)
}
