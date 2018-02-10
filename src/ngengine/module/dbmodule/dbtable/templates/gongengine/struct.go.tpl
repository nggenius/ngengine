package {{.Model}}

{{$module := .Model}}

{{$ilen := len .Imports}}
{{if gt $ilen 0}}
import (
	{{range .Imports}}"{{.}}"{{end}}
	"encoding/gob"
)
{{end}}


func init(){
	{{range .Tables}}
	gob.Register(&{{Mapper .Name}}{})
	gob.Register([]{{Mapper .Name}}{})
	{{end}}
}

{{range .Tables}}
type {{Mapper .Name}} struct {
{{$table := .}}
{{range .ColumnsSeq}}{{$col := $table.GetColumn .}}	{{Mapper $col.Name}}	{{Type $col}} {{Tag $table $col}}
{{end}}
}

{{end}}