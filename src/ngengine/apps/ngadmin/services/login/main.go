package main

import (
	"flag"
	"fmt"
	"io/ioutil"
	"ngengine/core"
	"ngengine/game/login"
	"os"
	"runtime/debug"
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
		panic("login parameter is empty")
	}

	core.RegisterService("login", new(login.Login))

	_, err := core.CreateService("login", *startPara)
	if err != nil {
		panic(err)
	}
	core.RunAllService()

	core.Wait()
	return
}
