package main

import (
	"ngengine/core"

	"github.com/mysll/toolkit"
)

func main() {
	core.RegisterService("database", &Database{})
	core.RegisterService("login", &Login{})
	_, err := core.CreateService("database", startargs)
	if err != nil {
		panic(err)
	}
	_, err = core.CreateService("login", startargs2)
	if err != nil {
		panic(err)
	}
	core.RunAllService()

	toolkit.WaitForQuit()
	core.CloseAllService()
	core.Wait()
}
