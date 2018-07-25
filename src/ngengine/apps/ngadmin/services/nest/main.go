package main

import (
	"flag"
	"fmt"
	"io/ioutil"
	"ngengine/core"
	"ngengine/game/nest"
	"os"
	"runtime/debug"

	"github.com/mysll/toolkit"
)

var startPara = flag.String("p", "", "startPara")

func main() {
	defer func() {
		if x := recover(); x != nil {
			d := fmt.Sprintf("panic(%v)\n%s", x, debug.Stack())
			ioutil.WriteFile("dump.log", []byte(d), 0666)

			os.Exit(0)
		}
	}()
	flag.Parse()

	if *startPara == "" {
		flag.PrintDefaults()
		panic("nest parameter is empty")
	}

	core.RegisterService("nest", new(nest.Nest))

	_, err := core.CreateService("nest", *startPara)
	if err != nil {
		panic(err)
	}
	core.RunAllService()

	toolkit.WaitForQuit()
	core.CloseAllService()
	core.Wait()
	return
}
