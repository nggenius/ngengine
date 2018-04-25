package main

import (
	"ngengine/core"

	_ "github.com/go-sql-driver/mysql"
	"github.com/mysll/toolkit"
)

func main() {
	// 捕获异常
	core.RegisterService("database", &Database{})
	core.RegisterService("login", &Login{})

	_, err := core.CreateService("login", startlogin)
	if err != nil {
		panic(err)
	}

	_, err = core.CreateService("database", dbargs)
	if err != nil {
		panic(err)
	}
	core.RunAllService()

	toolkit.WaitForQuit()
	core.CloseAllService()
	core.Wait()
}
