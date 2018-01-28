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
	Sid  share.ServiceId // 所在service的id
	Flag int8            // 见share.MB_FLAG_XXXX
	Id   uint64          // flag=share.MB_FLAG_CLIENT时，表示客户端的session.
	Uid  uint64          // 唯一id,由以上三个字段生成。主要用来传送
}

func (m Mailbox) String() string {
	return fmt.Sprintf("mailbox://%x", m.Uid)
}

func (m *Mailbox) ServiceId() share.ServiceId {
	return share.ServiceId(m.Sid)
}

func (m *Mailbox) IsClient() bool {
	return m.Flag == share.MB_FLAG_CLIENT
}

func (m *Mailbox) Generate() {
	m.Uid = ((uint64(m.Sid) << 48) & 0xFFFF000000000000) | ((uint64(m.Flag) & 1) << 47) | (m.Id & share.SESSION_MAX)
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
	mbox.Sid = share.ServiceId((mbox.Uid >> 48) & 0xFFFF)
	return mbox, nil
}

func NewMailboxFromUid(val uint64) Mailbox {
	mbox := Mailbox{}
	mbox.Uid = val
	mbox.Id = mbox.Uid & share.SESSION_MAX
	mbox.Flag = int8((mbox.Uid >> 47) & 1)
	mbox.Sid = share.ServiceId((mbox.Uid >> 48) & 0xFFFF)
	return mbox
}

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

func NewSessionMailbox(appid share.ServiceId, id uint64) Mailbox {
	if id > share.SESSION_MAX || appid > share.SID_MAX {
		panic("id is wrong")
	}
	m := Mailbox{}
	m.Sid = appid
	m.Flag = share.MB_FLAG_CLIENT
	m.Id = id
	m.Generate()
	return m
}

func NewMailbox(appid share.ServiceId, flag int8, id uint64) Mailbox {
	if id > share.SESSION_MAX || appid > share.SID_MAX {
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
	gob.Register(Mailbox{})
}
