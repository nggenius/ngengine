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

func (s *Sql) Init(core service.CoreAPI) (err error) {
	opt := core.Option()

	if !opt.Args.Has("db") {
		return fmt.Errorf("db not define")
	}
	db := opt.Args.String("db")

	if !opt.Args.Has("datasource") {
		return fmt.Errorf("datasource not define")
	}

	ds := opt.Args.String("datasource")
	s.orm, err = xorm.NewEngine(db, ds)

	s.orm.ShowSQL(opt.Args.Bool("showsql"))

	return err
}

func (s *Sql) Sync(bean interface{}) error {
	if s.orm == nil {
		return fmt.Errorf("orm is nil")
	}

	return s.orm.Sync2(bean)
}
