package main

import (
	"flag"
	"fmt"
	"html/template"
	"ngengine/console/models"
	"ngengine/console/routes"

	"github.com/tango-contrib/xsrf"

	"github.com/lunny/tango"
	"github.com/mysll/toolkit"
	"github.com/tango-contrib/events"
	"github.com/tango-contrib/renders"
)

var (
	port = flag.Int("p", 7000, "No Port")
)

func InitDb() error {
	err := models.InitDb()
	if err != nil {
		fmt.Println(err)
		return err
	}
	return nil

	/*info := &models.NxConsole{}
	info.ServerIp = "1.1.1.1"

	if err := info.Insert(); err != nil {
		fmt.Println(err)
		return err
	}*/
}

func Serv() {
	if err := InitDb(); err != nil {
		return
	}

	t := tango.Classic()
	t.Use(
		xsrf.New(100),
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
	flag.Parse()
	go Serv()
	var url = fmt.Sprintf("http://127.0.0.1:%d/", *port)
	toolkit.OpenBrowser(url)
	fmt.Print("Open Url :", url)
	toolkit.WaitForQuit()
}
