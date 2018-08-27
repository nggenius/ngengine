package main

import (
	"net/http"
	"ngengine/core"
	"ngengine/game/login"
	"ngengine/game/nest"
	"ngengine/game/region"
	"ngengine/game/store"
	"ngengine/game/world"

	_ "net/http/pprof"
)

var startnest = `{
	"ServId":5,
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
	"LogFile":"log/nest.log",
	"Args": {
		"MainEntity":"entity.GamePlayer"
	}
}`

var startlogin = `{
	"ServId":4,
	"ServType": "login",
	"AdminAddr":"127.0.0.1",
	"AdminPort":12500,
	"ServName": "login_1",
	"ServAddr": "127.0.0.1",
	"ServPort": 0,
	"Expose": true,
	"OuterAddr":"192.168.1.12",
	"HostAddr": "0.0.0.0",
	"HostPort": 4000,
	"LogFile":"log/login.log",
	"Args": {}
}`

var startworld = `{
	"ServId":3,
	"ServType": "world",
	"AdminAddr":"127.0.0.1",
	"AdminPort":12500,
	"ServName": "world_1",
	"ServAddr": "127.0.0.1",
	"ServPort": 0,
	"Expose": false,
	"LogFile":"log/world.log",
	"ResRoot":"D:/home/work/github/ngengine/res/",
	"Args": {
		"Region":"region.json",
		"MinRegions":1
	}
}`

var startregion = `{
	"ServId":2,
	"ServType": "region",
	"AdminAddr":"127.0.0.1",
	"AdminPort":12500,
	"ServName": "region_1",
	"ServAddr": "127.0.0.1",
	"ServPort": 0,
	"Expose": false,
	"LogFile":"log/region.log",
	"ResRoot":"D:/home/work/github/ngengine/res/",
	"Args": {}
}`

var dbargs = `{
	"ServId":1,
	"ServType": "store",
	"AdminAddr":"127.0.0.1",
	"AdminPort":12500,
	"ServName": "db_1",
	"ServAddr": "127.0.0.1",
	"ServPort": 0,
	"Expose": false,
	"HostAddr": "",
	"HostPort": 0,
	"LogFile":"log/db.log",
	"Args": {
		"db":"mysql",
		"datasource":"sa:abc@tcp(192.168.1.52:3306)/ngengine?charset=utf8",
		"showsql":false
	}
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

	go http.ListenAndServe(":9600", nil)
	core.Wait()
}
