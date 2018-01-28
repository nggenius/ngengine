package main

import (
	"ngengine/ngadmin"

	"github.com/mysll/toolkit"
)

func main() {
	opts := &ngadmin.Options{}
	err := opts.LoadFromFile("service.cfg")
	if err != nil {
		panic(err)
	}
	ngadmin := ngadmin.New(opts)
	ngadmin.Main()
	toolkit.WaitForQuit()
	ngadmin.Exit()
}
