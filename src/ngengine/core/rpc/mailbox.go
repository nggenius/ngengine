package rpc

import (
	"errors"
	"fmt"
	"ngengine/share"
	"strconv"
	"strings"
)

const (
	ID_MAX      = 0x7FFFFFFFFFFF // id最大值
	FLAG_MASK   = 0x800000000000 // flag mask
	NullMailbox = Mailbox(0)
)

type Mailbox uint64

func (m Mailbox) String() string {
	return fmt.Sprintf("mailbox://%x", uint64(m))
}

// 是否为空
func (m Mailbox) IsNil() bool {
	return m == 0
}

// 获取服务编号
func (m Mailbox) ServiceId() share.ServiceId {
	return share.ServiceId((m & 0xFFFF000000000000) >> 48)
}

// 获取标志位
func (m Mailbox) Flag() int8 {
	return int8((m >> 47) & 0x1)
}

// 获取id
func (m Mailbox) Id() uint64 {
	return uint64(m & ID_MAX)
}

// 获取uid
func (m Mailbox) Uid() uint64 {
	return uint64(m)
}

// 是否是一个客户端地址
func (m Mailbox) IsClient() bool {
	return m&FLAG_MASK > 0
}

// 是否是对象
func (m Mailbox) IsObject() bool {
	return (m&FLAG_MASK == 0) && ((m & ID_MAX) > 0)
}

// create uid
func generate(appid share.ServiceId, flag int8, id uint64) Mailbox {
	return Mailbox(((uint64(appid) << 48) & 0xFFFF000000000000) | ((uint64(flag) & 1) << 47) | (id & ID_MAX))
}

// 通过字符串生成mailbox
func NewMailboxFromStr(mb string) (Mailbox, error) {
	mbox := Mailbox(0)
	if !strings.HasPrefix(mb, "mailbox://") {
		return mbox, errors.New("mailbox string error")
	}
	vals := strings.Split(mb, "/")
	if len(vals) != 3 {
		return mbox, errors.New("mailbox string error")
	}

	var val uint64
	var err error

	val, err = strconv.ParseUint(vals[2], 16, 64)
	if err != nil {
		return mbox, err
	}
	mbox = Mailbox(val)
	return mbox, nil
}

// 通过uid生成mailbox
func NewMailboxFromUid(val uint64) Mailbox {
	return Mailbox(val)
}

// 通过服务编号获取mailbox
func GetServiceMailbox(appid share.ServiceId) Mailbox {
	if appid > 0xFFFF {
		panic("id is wrong")
	}
	m := generate(appid, 0, 0)
	return m
}

// 生成一个新的客户端mailbox
func NewSessionMailbox(appid share.ServiceId, id uint64) Mailbox {
	if id > ID_MAX || appid > share.SID_MAX {
		panic("id is wrong")
	}
	m := generate(appid, share.MB_FLAG_CLIENT, id)
	return m
}

// 生成mailbox
func NewMailbox(appid share.ServiceId, flag int8, id uint64) Mailbox {
	if id > ID_MAX || appid > share.SID_MAX {
		panic("id is wrong")
	}
	m := generate(appid, flag, id)
	return m
}

// get object identity
func (m Mailbox) Identity() int {
	return int((m >> 32) & share.OBJECT_TYPE_MAX)
}

// get object index
func (m Mailbox) ObjectIndex() int {
	return int(m & share.OBJECT_MAX)
}

// 生成一个新的object id
func (m Mailbox) NewObjectId(identity, serial, index int) Mailbox {
	id := uint64((identity&share.OBJECT_TYPE_MAX)<<32) | uint64(serial&0xFF)<<24 | uint64(index&share.OBJECT_MAX)
	return generate(m.ServiceId(), 0, id)
}
