package db

import (
	"encoding/gob"
	"encoding/json"
	"fmt"
	"ngengine/module/db/dbtable"
	"strings"

	_ "github.com/go-sql-driver/mysql"
	"github.com/go-xorm/xorm"
)

type Dbptr struct {
	xormptr *xorm.Engine
}

// InitDb初始化数据库指针
func InitDb(dbType string, dbpartem string) (*Dbptr, error) {
	engine, err := xorm.NewEngine(dbType, dbpartem)
	if err != nil {
		return nil, err
	}
	ptr := &Dbptr{}
	ptr.xormptr = engine
	return ptr, nil
}

func (p *Dbptr) GetForJson(tableName string, arg []byte) (interface{}, error) {
	sct, err := dbtable.DbstructPool(tableName)
	if err != nil {
		fmt.Println(err)
		return nil, err
	}

	if arg != nil {
		err1 := json.Unmarshal(arg, sct)
		if err1 != nil {
			fmt.Println(err1)
			return nil, err1
		}
	}

	has, err := p.xormptr.Get(sct)
	if !has {
		fmt.Println(err)
		return nil, err
	}

	return sct, nil
}

func (p *Dbptr) GetForGob(tableName string, arg []byte) (interface{}, error) {
	sct, err := dbtable.DbstructPool(tableName)
	if err != nil {
		fmt.Println(err)
		return nil, err
	}

	if arg != nil {
		f := string(arg)
		gob := gob.NewDecoder(strings.NewReader(f))
		gob.Decode(sct)
	}

	has, err := p.xormptr.Get(sct)
	if !has {
		fmt.Println(err)
		return nil, err
	}

	return sct, nil
}
