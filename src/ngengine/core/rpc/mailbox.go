package rpc

import (
	"encoding/gob"
	"errors"
	"fmt"
	"ngengine/share"
	"strconv"
	"strings"
)

type Mailbox struct {
	App  int32  //所在app的id
	Flag int8   //0:app 1:player
	Id   uint64 //flag=1时，表示客户端的session.
	Uid  uint64 //唯一id,由以上三个字段生成。主要用来传送
}

func (m Mailbox) String() string {
	return fmt.Sprintf("mailbox://%x", m.Uid)
}

func (m *Mailbox) ServiceId() share.ServiceId {
	return share.ServiceId(m.App)
}

func (m *Mailbox) IsClient() bool {
	return m.Flag == 1
}

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
	mbox.Id = mbox.Uid & share.SESSION_MAX
	mbox.Flag = int8((mbox.Uid >> 47) & 1)
	mbox.App = int32((mbox.Uid >> 48) & 0xFFFF)
	return mbox, nil
}

func NewMailboxFromUid(val uint64) Mailbox {
	mbox := Mailbox{}
	mbox.Uid = val
	mbox.Id = mbox.Uid & share.SESSION_MAX
	mbox.Flag = int8((mbox.Uid >> 47) & 1)
	mbox.App = int32((mbox.Uid >> 48) & 0xFFFF)
	return mbox
}

func GetServiceMailbox(appid share.ServiceId) Mailbox {
	if appid > 0xFFFF {
		panic("id is wrong")
	}
	m := Mailbox{}
	m.App = int32(appid)
	m.Flag = 0
	m.Id = 0
	m.Uid = ((uint64(appid) << 48) & 0xFFFF000000000000)
	return m
}

func NewMailbox(flag int8, id uint64, appid share.ServiceId) Mailbox {
	if id > share.SESSION_MAX || appid > 0xFFFF {
		panic("id is wrong")
	}
	m := Mailbox{}
	m.App = int32(appid)
	m.Flag = flag
	m.Id = id
	m.Uid = ((uint64(appid) << 48) & 0xFFFF000000000000) | ((uint64(flag) & 1) << 47) | (id & share.SESSION_MAX)
	return m
}

func init() {
	gob.Register(Mailbox{})
}
