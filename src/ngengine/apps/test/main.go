package main

import (
	"log"
	"ngengine/core"
	"os"
	"runtime/debug"

	_ "github.com/go-sql-driver/mysql"
	"github.com/mysll/toolkit"
)

func main() {
	defer func() {
		if e := recover(); e != nil {
			logFile, _ := os.Create("panic.log")
			defer logFile.Close()
			l := log.New(logFile, "", log.LstdFlags)
			l.Println(e)
			l.Println(string(debug.Stack()))
			os.Exit(1)
		}
	}()

	core.RegisterService("database", &Database{})
	core.RegisterService("object", &Object{})

	_, err := core.CreateService("object", objectargs)
	if err != nil {
		panic(err)
	}
	_, err = core.CreateService("database", dbargs)
	if err != nil {
		panic(err)
	}
	core.RunAllService()

	toolkit.WaitForQuit()
	core.CloseAllService()
	core.Wait()
}
