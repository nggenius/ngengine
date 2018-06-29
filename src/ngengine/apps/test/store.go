package main

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
	"LogFile":"db.log",
	"Args": {
		"db":"mysql",
		"datasource":"sa:abc@tcp(192.168.1.52:3306)/ngengine?charset=utf8"
	}
}`
