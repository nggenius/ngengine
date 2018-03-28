package main

import (
	"ngengine/core"

	"github.com/mysll/toolkit"
)

func main() {
	core.RegisterService("object", &Object{})
	_, err := core.CreateService("object", objectargs)
	if err != nil {
		panic(err)
	}
	core.RunAllService()

	toolkit.WaitForQuit()
	core.CloseAllService()
	core.Wait()
}
