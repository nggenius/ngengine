package {{.Model}}

{{$ilen := len .Imports}}
{{if gt $ilen 0}}
import (
	{{range .Imports}}"{{.}}"{{end}}
	"ngengine/module/db/dbtable"
	"reflect"
)
{{end}}

{{range .Tables}}

func init(){
	dbtable.DbtableMap["{{Mapper .Name}}"] = reflect.TypeOf({{Mapper .Name}}{})
}

type {{Mapper .Name}} struct {
{{$table := .}}
{{range .ColumnsSeq}}{{$col := $table.GetColumn .}}	{{Mapper $col.Name}}	{{Type $col}} {{Tag $table $col}}
{{end}}
}

{{end}}