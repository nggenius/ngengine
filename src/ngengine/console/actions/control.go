package actions

import (
	"github.com/lunny/tango"
)

type Control struct {
	RenderBase
	tango.Json
}

func (c *Control) Post() interface{} {
	type jsondata struct {
		Cmd  string `json:"cmd"`
		Args string `json:"args"`
	}
	var data jsondata
	err := c.DecodeJSON(&data)
	if err != nil {
		return map[string]interface{}{
			"err:":   err.Error(),
			"status": 500,
		}
	}

	return map[string]interface{}{
		"status": 200,
	}
}
