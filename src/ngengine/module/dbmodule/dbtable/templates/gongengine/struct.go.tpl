package {{.Model}}

{{$module := .Model}}

{{$ilen := len .Imports}}
{{if gt $ilen 0}}
import (
	{{range .Imports}}"{{.}}"{{end}}
	"fmt"
	"ngengine/module/dbmodule/dbtable"
)
{{end}}


func init(){
	if _, ok := dbtable.DbPtrMap["{{$module}}"]; !ok {
		fmt.Print("{{$module}} init defeated")
	}
	db := dbtable.DbPtrMap["{{$module}}"]

	{{range .Tables}}
	dbtable.RegisterTable(db, &{{Mapper .Name}}{}, []{{Mapper .Name}}{})
	{{end}}
}

{{range .Tables}}
type {{Mapper .Name}} struct {
{{$table := .}}
{{range .ColumnsSeq}}{{$col := $table.GetColumn .}}	{{Mapper $col.Name}}	{{Type $col}} {{Tag $table $col}}
{{end}}
}

{{end}}