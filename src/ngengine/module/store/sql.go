package store

import (
	"fmt"
	"ngengine/core/service"

	"github.com/go-xorm/xorm"
)

type Sql struct {
	orm *xorm.Engine
}

func newSql() *Sql {
	s := &Sql{}
	return s
}

func (s *Sql) Init(core service.CoreApi) (err error) {
	var db, ds string
	var has bool
	opt := core.Option()
	if db, has = opt.Args["db"]; !has {
		return fmt.Errorf("db not define")
	}
	if ds, has = opt.Args["datasource"]; !has {
		return fmt.Errorf("datasource not define")
	}
	s.orm, err = xorm.NewEngine(db, ds)
	return err
}

func (s *Sql) Sync(bean interface{}) error {
	if s.orm == nil {
		return fmt.Errorf("orm is nil")
	}

	return s.orm.Sync2(bean)
}
