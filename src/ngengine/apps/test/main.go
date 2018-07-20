package main

import (
	"ngengine/core"
	"ngengine/game/login"
	"ngengine/game/nest"
	"ngengine/game/region"
	"ngengine/game/store"
	"ngengine/game/world"

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

var startworld = `{
	"ServId":5,
	"ServType": "world",
	"AdminAddr":"127.0.0.1",
	"AdminPort":12500,
	"ServName": "world_1",
	"ServAddr": "127.0.0.1",
	"ServPort": 0,
	"Expose": false,
	"LogFile":"world.log",
	"ResRoot":"D:/home/work/github/ngengine/res/",
	"Args": {
		"Region":"region.json",
		"MinRegions":1
	}
}`

var startregion = `{
	"ServId":6,
	"ServType": "region",
	"AdminAddr":"127.0.0.1",
	"AdminPort":12500,
	"ServName": "region_1",
	"ServAddr": "127.0.0.1",
	"ServPort": 0,
	"Expose": false,
	"LogFile":"region.log",
	"ResRoot":"D:/home/work/github/ngengine/res/",
	"Args": {}
}`

func main() {
	// 捕获异常
	core.RegisterService("store", new(store.Store))
	core.RegisterService("login", new(login.Login))
	core.RegisterService("nest", new(nest.Nest))
	core.RegisterService("world", new(world.World))
	core.RegisterService("region", new(region.Region))
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

	_, err = core.CreateService("world", startworld)
	if err != nil {
		panic(err)
	}

	_, err = core.CreateService("region", startregion)
	if err != nil {
		panic(err)
	}

	core.RunAllService()

	toolkit.WaitForQuit()
	core.CloseAllService()
	core.Wait()
}
