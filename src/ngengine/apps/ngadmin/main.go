package main

import (
	"flag"
	"fmt"
	"io/ioutil"
	"ngengine/ngadmin"
	"os"
	"runtime/debug"

	"github.com/mysll/toolkit"
)

var (
	configPath = flag.String("p", "../config/", "app config path")
)

func main() {
	defer func() {
		if x := recover(); x != nil {
			d := fmt.Sprintf("panic(%v)\n%s", x, debug.Stack())
			ioutil.WriteFile("dump.log", []byte(d), 0666)

			os.Exit(0)
		}
	}()

	flag.Parse()

	var config ngadmin.Options

	err := config.Load(*configPath+"app.cfg", *configPath+"servers.cfg")
	if err != nil {
		panic(err)
	}

	ngadmin := ngadmin.New(&config)
	ngadmin.Main()
	toolkit.WaitForQuit()

	ngadmin.Exit()
}
