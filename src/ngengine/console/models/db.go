package models

import (
	"errors"

	_ "github.com/go-sql-driver/mysql"
	"github.com/go-xorm/xorm"
)

var (
	errInitA = errors.New("Already Init Db")
	errInitF = errors.New("Init Db Failed")
)

var XormEngine *xorm.Engine

func InitDb() error {
	if XormEngine != nil {
		return errInitA
	}

	engine, err := xorm.NewEngine("mysql", "root:123456@tcp(127.0.0.1:3306)/ddd?charset=utf8")
	if err != nil {
		return err
	}
	XormEngine = engine

	isexist, err := XormEngine.IsTableExist(&NxConsole{})
	if err != nil {
		return err
	}

	if !isexist {
		err := XormEngine.CreateTables(&NxConsole{})
		if err != nil {
			return err
		}
	}

	XormEngine.ShowSQL(true)
	return err
}
