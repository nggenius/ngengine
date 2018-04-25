package inner

import (
	"time"
)

type Role struct {
	Id          int64
	Index       int8
	Account     string `xorm:"varchar(128) index"`
	RoleName    string `xorm:"varchar(128) unique"`
	CreateTime  time.Time
	LastLogTime time.Time
	LastAddress string `xorm:"varchar(32)"`
	Status      int8
}

type RoleCreater struct {
}

func (c *RoleCreater) Create() interface{} {
	return &Role{}
}

func (c *RoleCreater) CreateSlice() interface{} {
	return &[]*Role{}
}
