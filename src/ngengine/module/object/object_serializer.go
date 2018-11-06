package object

import (
	"fmt"
	"ngengine/core/rpc"
	"ngengine/protocol"
	"ngengine/share"
	"ngengine/utils"
)

type Container interface {
	Cap() int
	SetCap(cap int) error
	ChildCount() int
	ParentIndex() int
	ChildAtIf(pos int) interface{}
	FirstChildIf() (int, interface{})
	NextChildIf(index int) (int, interface{})
	AddChildIf(pos int, g interface{}) (int, error)
}

func (f *Factory) Encode(o interface{}) ([]byte, error) {
	buf := protocol.NewMessage(share.INNER_MESSAGE_BUF_LEN)
	ar := utils.NewStoreArchiver(buf.Body)
	err := encode(ar, o)
	if err != nil {
		return nil, err
	}
	return ar.Data(), nil
}

func encode(ar *utils.StoreArchive, o interface{}) (err error) {
	if oc, ok := o.(ObjectCreate); ok {
		typ := oc.EntityType()
		if typ == "" {
			return fmt.Errorf("encode object type is nil")
		}
		err = ar.Put(typ)
		if err != nil {
			return err
		}
	} else {
		return fmt.Errorf("object not implement ObjectCreate")
	}

	pos := 0
	caps := 0
	childs := 0
	var c Container
	var ok bool
	c, ok = o.(Container)
	if ok {
		pos = c.ParentIndex()
		caps = c.Cap()
		childs = c.ChildCount()
	}

	err = ar.Put(caps) // 容量
	if err != nil {
		return err
	}
	err = ar.Put(pos) // 位置
	if err != nil {
		return err
	}
	err = ar.Put(childs)
	if err != nil {
		return err
	}

	if obj, ok := o.(Object); ok {
		if obj.Original() != nil {
			err = ar.Put(*obj.Original())
			if err != nil {
				return err
			}
		} else {
			err = ar.Put(obj.ObjId())
			if err != nil {
				return err
			}
		}
	} else {
		return fmt.Errorf("object is not implement Object")
	}

	err = ar.Put(o)
	if err != nil {
		return err
	}

	if ok {
		count := 0
		i, o := c.FirstChildIf()
		for o != nil {
			err = encode(ar, o)
			if err != nil {
				return err
			}
			i, o = c.NextChildIf(i)
			count++
		}
		if childs != count {
			return fmt.Errorf("child count not match")
		}
	}

	return nil
}

func (f *Factory) Decode(b []byte) (interface{}, error) {
	ar := utils.NewLoadArchiver(b)
	o, err := f.decodeObj(nil, ar)
	return o, err
}

func (f *Factory) decodeObj(parent Container, ar *utils.LoadArchive) (interface{}, error) {
	var typ string
	err := ar.Read(&typ)
	if err != nil {
		return nil, err
	}

	pos := 0
	caps := 0
	childs := 0

	err = ar.Read(&pos)
	if err != nil {
		return nil, err
	}

	err = ar.Read(&caps)
	if err != nil {
		return nil, err
	}

	err = ar.Read(&childs)
	if err != nil {
		return nil, err
	}

	var origin rpc.Mailbox
	err = ar.Read(&origin)

	o, err := f.Create(typ)
	if err != nil {
		return nil, err
	}

	err = ar.Read(o)
	if err != nil {
		f.Destroy(o)
		return nil, err
	}

	obj := o.(Object)
	obj.SetDummy(true)
	obj.SetOriginal(&origin)

	c, ok := o.(Container)
	if ok {
		c.SetCap(caps)

		if parent != nil {
			parent.AddChildIf(pos, o)
		}

		for i := 0; i < childs; i++ {
			_, err := f.decodeObj(c, ar)
			if err != nil {
				f.Destroy(o)
				return nil, err
			}
		}
	}

	return o, nil
}
