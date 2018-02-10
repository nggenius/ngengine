package actions

import (
	"github.com/tango-contrib/renders"
	"github.com/tango-contrib/xsrf"
)

type Index struct {
	RenderBase
	xsrf.Checker
}

func (i *Index) Get() error {
	return i.Render("index.html", renders.T{
		"XsrfFormHtml": i.XsrfFormHtml(),
	})
}
