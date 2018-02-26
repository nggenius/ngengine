package dbtable

import (
	"encoding/gob"
	"errors"
	"fmt"

	_ "github.com/go-sql-driver/mysql"
	"github.com/go-xorm/xorm"
)

// 当前链接的数据库
const (
	nxbase string = "nx_base"
)

// DbPtrMap 存放全部的连接
var DbPtrMap map[string]*DbPtrs

// DbPtrs 数据库指针
type DbPtrs struct {
	DbPtr *xorm.Engine
}

func init() {
	// 后面通过配置加载
	dbptr, err := InitDb("mysql", "root:@tcp(127.0.0.1:3306)/nx_base?charset=utf8")
	if err != nil {
		return
	}
	RegisterDb(nxbase, dbptr)
}

// IsDbConnect 检查是否有这个数据库连接
func IsDbConnect(ConnectName string) bool {

	if _, ok := DbPtrMap[ConnectName]; !ok {
		return false
	}
	if nil == DbPtrMap[ConnectName].DbPtr {
		return false
	}

	return true
}

// GetDbPtr 通过名字获取是否有这个数据库的连接
func GetDbPtr(dbName string) (*xorm.Engine, error) {
	if !IsDbConnect(dbName) {
		return nil, errors.New("dont have this connect")
	}

	return DbPtrMap[dbName].DbPtr, nil
}

// InitDb 初始化数据库指针
func InitDb(dbType string, dbpartem string) (*xorm.Engine, error) {
	engine, err := xorm.NewEngine(dbType, dbpartem)
	if err != nil {
		fmt.Println(err)
		return nil, err
	}
	return engine, nil
}

// RegisterDb 注册需要连接的数据库
func RegisterDb(dbName string, dbPtr *xorm.Engine) bool {
	if dbName == "" && dbPtr == nil {
		return false
	}

	if nil == DbPtrMap {
		DbPtrMap = make(map[string]*DbPtrs)
	}

	if _, ok := DbPtrMap[dbName]; ok {
		panic("This DbName had")
	}
	DbPtrMap[dbName] = &DbPtrs{dbPtr}

	return true
}

func RegisterTable(db *DbPtrs, tableStruct interface{}, tableStructArray interface{}) {
	if nil == db {
		return
	}

	db.DbPtr.Sync2(tableStruct)
	gob.Register(tableStruct)
	gob.Register(tableStructArray)
}
