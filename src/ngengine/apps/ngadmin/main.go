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
	appPath = flag.String("ap", "./app.cfg", "app config path")

	appAresPath = flag.String("para", "./servers.cfg", "app parameter config path")
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
	err := config.Load(*appPath, *appAresPath)
	if err != nil {
		panic(err)
	}

	ngadmin := ngadmin.New(&config)
	ngadmin.Main()
	toolkit.WaitForQuit()
	ngadmin.Exit()
}
