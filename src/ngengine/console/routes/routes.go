package routes

import (
	"ngengine/console/actions"

	"github.com/lunny/tango"
)

func SetRoutes(t *tango.Tango) {
	t.Get("/", new(actions.Index))
}
