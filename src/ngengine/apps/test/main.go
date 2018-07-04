package main

import (
	"ngengine/core"
	"ngengine/game/login"
	"ngengine/game/nest"
	"ngengine/game/store"

	_ "github.com/go-sql-driver/mysql"
	"github.com/mysll/toolkit"
)

var startnest = `{
	"ServId":4,
	"ServType": "nest",
	"AdminAddr":"127.0.0.1",
	"AdminPort":12500,
	"ServName": "nest_1",
	"ServAddr": "127.0.0.1",
	"ServPort": 0,
	"Expose": true,
	"OuterAddr":"192.168.1.12",
	"HostAddr": "0.0.0.0",
	"HostPort": 0,
	"LogFile":"nest.log",
	"Args": {
		"MainEntity":"entity.Player"
	}
}`

var startlogin = `{
	"ServId":3,
	"ServType": "login",
	"AdminAddr":"127.0.0.1",
	"AdminPort":12500,
	"ServName": "login_1",
	"ServAddr": "127.0.0.1",
	"ServPort": 0,
	"Expose": true,
	"OuterAddr":"192.168.1.12",
	"HostAddr": "0.0.0.0",
	"HostPort": 2002,
	"LogFile":"login.log",
	"Args": {}
}`

func main() {
	// 捕获异常
	core.RegisterService("store", &store.Store{})
	core.RegisterService("login", &login.Login{})
	core.RegisterService("nest", &nest.Nest{})

	_, err := core.CreateService("login", startlogin)
	if err != nil {
		panic(err)
	}

	_, err = core.CreateService("nest", startnest)
	if err != nil {
		panic(err)
	}

	_, err = core.CreateService("store", dbargs)
	if err != nil {
		panic(err)
	}
	core.RunAllService()

	toolkit.WaitForQuit()
	core.CloseAllService()
	core.Wait()
}
