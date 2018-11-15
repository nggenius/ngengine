package main

import (
	"github.com/mysll/toolkit"
	"github.com/nggenius/ngengine/ngadmin"
)

func main() {
	var config ngadmin.Options
	config.Load("./app.cfg", "./servers.cfg")
	ngadmin := ngadmin.New(&config)
	ngadmin.Main()
	toolkit.WaitForQuit()
	ngadmin.Exit()
}
