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

func (m Mailbox) Flag() int8 {
	return int8((m >> 47) & 0x1)
}

func (m Mailbox) Id() uint64 {
	return uint64(m & ID_MAX)
}

func (m Mailbox) Uid() uint64 {
	return uint64(m)
}

// create uid
func generate(appid share.ServiceId, flag int8, id uint64) Mailbox {
	return Mailbox(((uint64(appid) << 48) & 0xFFFF000000000000) | ((uint64(flag) & 1) << 47) | (id & ID_MAX))
}

// 是否是一个客户端地址
func (m *Mailbox) IsClient() bool {
	return m.Flag() == share.MB_FLAG_CLIENT
}

// get object type
func (m Mailbox) ObjectType() int {
	return int((m >> 40) & 0x7F)
}

// get object index
func (m Mailbox) ObjectIndex() int {
	return int(m & 0xFFFFFFFF)
}

// 生成一个新的object id
func (m Mailbox) NewObjectId(otype, serial, index int) Mailbox {
	id := uint64((otype&0x7F)<<40) | uint64(serial&0xFF)<<32 | uint64(index&0xFFFFFFFF)
	m &= Mailbox(^uint64(ID_MAX))
	m |= Mailbox(id)
	return m
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
	mbox := Mailbox(val)
	return mbox
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
