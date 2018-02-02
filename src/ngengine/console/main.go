package main

import (
	"flag"
	"fmt"
	"html/template"
	"ngengine/console/routes"

	"github.com/lunny/tango"
	"github.com/mysll/toolkit"
	"github.com/tango-contrib/events"
	"github.com/tango-contrib/renders"
)

var (
	port = flag.Int("p", 7000, "No Port")
)

func Serv() {
	t := tango.Classic()
	t.Use(
		events.Events(),
		tango.Static(tango.StaticOptions{
			RootPath: "./view/statics/assets",
			Prefix:   "assets",
		}),
		renders.New(renders.Options{
			Reload:    true,
			Directory: "./view/templates",
			Funcs:     template.FuncMap{},
		}),
	)
	routes.SetRoutes(t)
	t.Run(fmt.Sprintf("127.0.0.1:%d", *port))
}

func main() {
	go Serv()
	var url = fmt.Sprintf("http://127.0.0.1:%d/", *port)
	toolkit.OpenBrowser(url)
	fmt.Print("Open Url :", url)
	toolkit.WaitForQuit()
}
