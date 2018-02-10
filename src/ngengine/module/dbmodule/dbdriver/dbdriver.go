package dbdriver

import (
	"database/sql"
	"errors"
	"fmt"
	"ngengine/module/dbmodule/dbtable"

	"reflect"

	_ "github.com/go-sql-driver/mysql"
)

// Get dbname: 数据的名字 arg:查询数据对应的结构
func Get(dbname string, arg interface{}) (interface{}, error) {
	dbptr, err := dbtable.GetDbPtr(dbname)
	if err != nil {
		fmt.Println(err)
		return nil, err
	}
	has, err := dbptr.Get(arg)
	if !has {
		fmt.Println(err)
		return nil, err
	}

	return arg, nil
}

// SQLQuery dbname: 数据的名字 sql:sql文
func SQLQuery(dbname string, sql string) ([]map[string][]byte, error) {
	if sql == "" {
		return nil, errors.New("sql is nil")
	}

	dbptr, err := dbtable.GetDbPtr(dbname)
	if err != nil {
		fmt.Println(err)
		return nil, err
	}

	reslut, err1 := dbptr.Query(sql)
	if err1 != nil {
		return nil, err1
	}
	return reslut, nil
}

// SQLExec dbname:数据库名字 sql:sql文 args:sql文的参数
func SQLExec(dbname string, sql string, args ...interface{}) (sql.Result, error) {
	dbptr, err := dbtable.GetDbPtr(dbname)
	if err != nil {
		fmt.Println(err)
		return nil, err
	}
	return dbptr.Exec(sql, args...)
}

// Insert dbname:数据库名字 sct:插入的对应数据结构和数据
func Insert(dbname string, sct interface{}) error {
	dbptr, err := dbtable.GetDbPtr(dbname)
	if err != nil {
		fmt.Println(err)
		return err
	}
	sctType := reflect.TypeOf(sct)
	if sctType.Elem().Kind() != reflect.Struct {
		return errors.New("parameter must be a struct")
	}

	insertRow, err := dbptr.Insert(sct)
	if err != nil {
		return err
	}
	fmt.Println(insertRow)
	return nil
}

// Update dbname:数据库名字 sct:更新的对应数据结构和数据 oldstc:需要更新的结构参数
func Update(dbname string, stc, oldstc interface{}) error {
	dbptr, err := dbtable.GetDbPtr(dbname)
	if err != nil {
		fmt.Println(err)
		return err
	}
	dbptr.Update(stc, oldstc)
	return nil
}

// Delete dbname:数据库名字 stc:删除的结构参数
func Delete(dbname string, stc interface{}) error {
	dbptr, err := dbtable.GetDbPtr(dbname)
	if err != nil {
		fmt.Println(err)
		return err
	}

	dbptr.Delete(stc)
	return nil
}

// Find dbname:数据库名字 stc:查询条件
func Find(dbname string, stc interface{}) (interface{}, error) {
	dbptr, err := dbtable.GetDbPtr(dbname)
	if err != nil {
		fmt.Println(err)
		return nil, err
	}

	newSlice := reflect.SliceOf(reflect.TypeOf(stc).Elem())

	newValue := reflect.New(newSlice)
	k := newValue.Interface()
	err = dbptr.Find(k, stc)
	if err != nil {
		fmt.Println(err)
		return nil, err
	}

	return k, nil
}
