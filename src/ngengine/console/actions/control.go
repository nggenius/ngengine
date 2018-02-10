package actions

import (
	"encoding/json"
	"ngengine/console/logic"

	"github.com/lunny/tango"
	"github.com/tango-contrib/xsrf"
)

type Control struct {
	RenderBase
	xsrf.Checker
	tango.Json
}

func (c *Control) Post() interface{} {
	type jsondata struct {
		Cmd  string `json:"cmd"`
		Args string `json:"args"`
	}
	var data jsondata
	//err := c.DecodeJSON(&data)
	var cmd = c.Ctx.Context.Req().FormValue("cmd")
	err := json.Unmarshal([]byte(cmd), &data)

	if err != nil {
		return ReturnError(err.Error())
	}

	switch data.Cmd {
	case "list":
		jsonData, err := logic.GetServerList()
		if err != nil {
			return ReturnError(err.Error())
		}
		return ReturnResult(jsonData)
	case "start":
	case "close":
	case "restart":
	case "config":
	default:
		return ReturnError(err.Error())
	}

	return ReturnError("")
}

func ReturnError(err string) map[string]interface{} {
	return map[string]interface{}{
		"err:":   err,
		"status": 500,
	}
}

func ReturnResult(data interface{}) map[string]interface{} {
	return map[string]interface{}{
		"status": 200,
		"data":   data,
	}
}
