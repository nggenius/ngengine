package main

import (
	"errors"

	"github.com/tango"
)

type Action struct {
	tango.JSON
}

func (Action) Get() interface{} {
	if true {
		return map[string]string{
			"say": "Hello tango!",
		}
	}
	return errors.New("something error")
}

func main() {
	t := tango.Classic()
	t.Get("/", new(Action))
	t.Run()
}
