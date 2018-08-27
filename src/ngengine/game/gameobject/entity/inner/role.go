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

// set id
func (r *Role) SetId(val int64) {
	r.Id = val
}

// db id
func (r *Role) DBId() int64 {
	return r.Id
}
