package rpc

import (
	"encoding/gob"
	"errors"
	"fmt"
	"ngengine/share"
	"strconv"
	"strings"
)

const (
	ID_MAX = 0x7FFFFFFFFFFF // id最大值
)

type Mailbox struct {
	Sid  share.ServiceId // 所在service的id
	Flag int8            // 见share.MB_FLAG_XXXX
	Id   uint64          // flag=share.MB_FLAG_CLIENT时，表示客户端的session.
	Uid  uint64          // 唯一id,由以上三个字段生成。主要用来传送
}

func (m Mailbox) String() string {
	return fmt.Sprintf("mailbox://%x", m.Uid)
}

// 是否为空
func (m Mailbox) IsNil() bool {
	if m.Sid == 0 && m.Flag == 0 && m.Id == 0 {
		return true
	}
	return false
}

// 获取服务编号
func (m *Mailbox) ServiceId() share.ServiceId {
	return share.ServiceId(m.Sid)
}

// 是否是一个客户端地址
func (m *Mailbox) IsClient() bool {
	return m.Flag == share.MB_FLAG_CLIENT
}

// create uid
func (m *Mailbox) Generate() {
	m.Uid = ((uint64(m.Sid) << 48) & 0xFFFF000000000000) | ((uint64(m.Flag) & 1) << 47) | (m.Id & ID_MAX)
}

// get object type
func (m Mailbox) ObjectType() int {
	return int((m.Id >> 40) & 0x7F)
}

// get object index
func (m Mailbox) ObjectIndex() int {
	return int(m.Id & 0xFFFFFFFF)
}

// 生成一个新的object id
func (m Mailbox) NewObjectId(otype, serial, index int) Mailbox {
	id := uint64((otype&0x7F)<<40) | uint64(serial&0xFF)<<32 | uint64(index&0xFFFFFFFF)
	mb := Mailbox{m.Sid, 0, id, 0}
	mb.Generate()
	return mb
}

// 通过字符串生成mailbox
func NewMailboxFromStr(mb string) (Mailbox, error) {
	mbox := Mailbox{}
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
	mbox.Uid = val
	mbox.Id = mbox.Uid & ID_MAX
	mbox.Flag = int8((mbox.Uid >> 47) & 1)
	mbox.Sid = share.ServiceId((mbox.Uid >> 48) & 0xFFFF)
	return mbox, nil
}

// 通过uid生成mailbox
func NewMailboxFromUid(val uint64) Mailbox {
	mbox := Mailbox{}
	mbox.Uid = val
	mbox.Id = mbox.Uid & ID_MAX
	mbox.Flag = int8((mbox.Uid >> 47) & 1)
	mbox.Sid = share.ServiceId((mbox.Uid >> 48) & 0xFFFF)
	return mbox
}

// 通过服务编号获取mailbox
func GetServiceMailbox(appid share.ServiceId) Mailbox {
	if appid > 0xFFFF {
		panic("id is wrong")
	}
	m := Mailbox{}
	m.Sid = appid
	m.Flag = share.MB_FLAG_APP
	m.Id = 0
	m.Generate()
	return m
}

// 生成一个新的客户端mailbox
func NewSessionMailbox(appid share.ServiceId, id uint64) Mailbox {
	if id > ID_MAX || appid > share.SID_MAX {
		panic("id is wrong")
	}
	m := Mailbox{}
	m.Sid = appid
	m.Flag = share.MB_FLAG_CLIENT
	m.Id = id
	m.Generate()
	return m
}

// 生成mailbox
func NewMailbox(appid share.ServiceId, flag int8, id uint64) Mailbox {
	if id > ID_MAX || appid > share.SID_MAX {
		panic("id is wrong")
	}
	m := Mailbox{}
	m.Sid = appid
	m.Flag = flag
	m.Id = id
	m.Generate()
	return m
}

func init() {
	gob.Register([]Mailbox{})
	gob.Register(Mailbox{})
}
