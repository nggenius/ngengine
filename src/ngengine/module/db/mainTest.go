package main

import (
	"bytes"
	"encoding/gob"
	"encoding/json"
	"fmt"
	"ngengine/module/db/dbdrive"
	"ngengine/module/db/dbtable/dbmodel"
)

func main() {
	aa, err := db.InitDb("mysql", "root:@tcp(127.0.0.1:3306)/nx_base?charset=utf8")
	if err != nil {
		fmt.Println(err)
		return
	}
	g := &dbmodel.NxChangename{
		Name: "1234",
	}

	var fout bytes.Buffer
	enc := gob.NewEncoder(&fout)
	enc.Encode(g)
	aa.Get1("NxChangename", fout.Bytes())

	if Data, err := json.Marshal(g); err == nil {
		reslut, _ := aa.Get("NxChangename", Data)
		fmt.Println(reslut)
		return
	}

	return
}
