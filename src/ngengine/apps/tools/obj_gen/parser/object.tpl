// Code generated by data parser.
// DO NOT EDIT!
package {{.Package}}

import(
    "encoding/json"
	"encoding/gob"
    "fmt"
	"ngengine/module/object"

    "github.com/mysll/toolkit"
)

var _ = json.Marshal
var _ = toolkit.ParseNumber
{{range .Property}}
{{if eq .Type "tuple"}}
// tuple {{.Name}} {{.Desc}}
type {{$.Name}}{{.Name}}_t struct {
	root object.Object
    {{range .Tuple}}
	{{.Name}} {{.Type}} // {{.Desc}}{{end}}
}

// tuple {{.Name}} construct
func New{{$.Name}}{{.Name}}(root object.Object) *{{$.Name}}{{.Name}}_t {
	{{tolower .Name}} := &{{$.Name}}{{.Name}}_t{root:root}
	return {{tolower .Name}}
}

// tuple {{.Name}} equal other
func ({{$src := tolower .Name }}{{$src}} *{{$.Name}}{{.Name}}_t) Equal(other {{$.Name}}{{.Name}}_t) bool {
	if {{range $k, $t := .Tuple}} {{if ne $k 0}} && {{end}}({{$src}}.{{.Name}} == other.{{.Name}}) {{end}} {
		return true
	}
	return false
}

{{end}}

{{if eq .Type "table"}}
// record {{.Name}} row define
type {{$.Name}}{{.Name}}_c struct {
    {{range .Table.Cols}}
	{{.Name}} {{.Type}} // {{.Desc}}{{end}}
}

// record {{.Name}} {{.Desc}}
type {{$.Name}}{{.Name}}_r struct {
	root object.Object
	data    [{{.Table.MaxRows}}]*{{$.Name}}{{.Name}}_c
    Row     []*{{$.Name}}{{.Name}}_c
}

// record  {{.Name}}  serial
type {{$.Name}}{{.Name}}Json struct{
    ColName []string
    ColType []string
    Row [][]interface{}
}

// record {{.Name}} construct
func New{{$.Name}}{{.Name}}(root object.Object) *{{$.Name}}{{.Name}}_r {
	{{tolower .Name}} := &{{$.Name}}{{.Name}}_r{root:root}
	{{tolower .Name}}.Row = {{tolower .Name}}.data[:0]
	return {{tolower .Name}}
}


{{$expose := .Expose}}
{{$pname := .Name}}

{{range $index, $col := .Table.Cols}}{{with $col}}
// get {{.Name}}
func (r *{{$.Name}}{{$pname}}_r) {{.Name}}(rownum int) ({{.Type}}, error) {
	if rownum < 0 || rownum >= len(r.Row) {
        return {{if eq .Type "string"}}""{{else}}0{{end}}, fmt.Errorf("row num error")
	}
	return r.Row[rownum].{{.Name}}, nil
}

// set {{.Name}}
func (r *{{$.Name}}{{$pname}}_r) Set{{.Name}}(rownum int, {{tolower .Name}} {{.Type}}) error {
	if rownum < 0 || rownum >= len(r.Row) {
        return fmt.Errorf("row num error")
	}
	if r.Row[rownum].{{.Name}} != {{tolower .Name}} {
		r.Row[rownum].{{.Name}} = {{tolower .Name}}{{if ne $expose ""}}
		if r.root != nil {
			r.root.ChangeTable("{{$pname}}", rownum, {{$index}}, {{tolower .Name}})
		}{{end}}
	}
	return nil
}
{{end}}{{end}}

// set row value
func (r *{{$.Name}}{{.Name}}_r) SetRowValue(rownum int {{range .Table.Cols}}, {{tolower .Name}} {{.Type}} {{end}} ) error {
	if rownum < 0 || rownum >= len(r.Row) {
		return fmt.Errorf("row num error")
	}
    /*{{range $index, $col := .Table.Cols}}{{with $col}}
	if r.Row[rownum].{{.Name}} != {{tolower .Name}} {
		r.Row[rownum].{{.Name}} = {{tolower .Name}}{{if ne $expose ""}}
		if r.root != nil {
			r.root.ChangeTable("{{$pname}}", rownum, {{$index}}, {{tolower .Name}})
		} {{end}}{{end}}
	} {{end}}
	*/
	{{range $index, $col := .Table.Cols}}{{with $col}}
	r.Row[rownum].{{.Name}} = {{tolower .Name}}{{end}}{{end}}
	if r.root != nil {
		r.root.SetTableRowValue("{{$pname}}", rownum, 0)
	}
	return nil
}

// get row value
func (r *{{$.Name}}{{.Name}}_r) RowValue(rownum int) ({{range .Table.Cols}}{{.Type}},{{end}} error) {
	var row {{$.Name}}{{.Name}}_c
	if rownum < 0 || rownum >= len(r.Row) {
		return {{range .Table.Cols}}row.{{.Name}},{{end}} fmt.Errorf("row num error")
	}

	row = *r.Row[rownum]
	return {{range .Table.Cols}}row.{{.Name}},{{end}} nil
}

// add row
func (r *{{$.Name}}{{.Name}}_r) AddRow(rownum int) (int, error) {
	if len(r.Row) > cap(r.data) { // full
		return -1, fmt.Errorf("record {{$.Name}}{{.Name}} is full")
	}

	if rownum < -1 || rownum >= cap(r.data) { // out of range
		return -1, fmt.Errorf("record {{$.Name}}{{.Name}} row %d out of range", rownum)
	}

	size := len(r.Row)
	row := &{{$.Name}}{{.Name}}_c{}
	r.Row = r.data[:size+1]
	if rownum == -1 || rownum == size {
		r.Row[size] = row{{if ne $expose ""}}
		if r.root != nil {
			r.root.AddTableRow("{{$pname}}", rownum)
		} {{end}}
		return size, nil
	}
	copy(r.Row[rownum+1:], r.Row[rownum:])
	r.Row[rownum] = row	{{if ne $expose ""}}
	if r.root != nil {
		r.root.AddTableRow("{{$pname}}", rownum)
	} {{end}}
	return rownum, nil
}

// add row value
func (r *{{$.Name}}{{.Name}}_r) AddRowValue(rownum int {{range .Table.Cols}}, {{tolower .Name}} {{.Type}} {{end}} ) (int, error) {
	if len(r.Row) > cap(r.data) { // full
		return -1, fmt.Errorf("record {{$.Name}}{{.Name}} is full")
	}

	if rownum < -1 || rownum >= cap(r.data) { // out of range
		return -1, fmt.Errorf("record {{$.Name}}{{.Name}} row %d out of range", rownum)
	}

	size := len(r.Row)
	row := &{{$.Name}}{{.Name}}_c{ {{range $k, $v := .Table.Cols}} {{if ne $k 0}},{{end}} {{tolower .Name}}{{end}} }
	r.Row = r.data[:size+1]
	if rownum == -1 || rownum == size {
		r.Row[size] = row{{if ne $expose ""}}
		if r.root != nil {
			r.root.AddTableRowValue("{{$pname}}", rownum, {{range $k, $v := .Table.Cols}} {{if ne $k 0}},{{end}} {{tolower .Name}}{{end}} )
		} {{end}}
		return size, nil
	}
	copy(r.Row[rownum+1:], r.Row[rownum:])
	r.Row[rownum] = row	{{if ne $expose ""}}
	if r.root != nil {
		r.root.AddTableRowValue("{{$pname}}", rownum, {{range $k, $v := .Table.Cols}} {{if ne $k 0}},{{end}} {{tolower .Name}}{{end}} )
	} {{end}}
	return rownum, nil
}

// del row
func (r *{{$.Name}}{{.Name}}_r) Del(rownum int) error {
	if rownum < 0 || rownum >= len(r.Row) {
		return fmt.Errorf("row num error")
	}
	copy(r.Row[rownum:], r.Row[rownum+1:])
	r.Row = r.data[:len(r.Row)-1]{{if ne $expose ""}}
	if r.root != nil {
		r.root.DelTableRow("{{$pname}}", rownum )
	} {{end}}	
	return nil
}

// clear
func (r *{{$.Name}}{{.Name}}_r) Clear() {
	r.Row = r.data[:0]{{if ne $expose ""}}
	if r.root != nil {
		r.root.ClearTable("{{$pname}}")
	} {{end}}
}

// json encode interface
func (r *{{$.Name}}{{.Name}}_r) Marshal() ([]byte, error) {
    return r.pack()
}

// json decode interface
func (r *{{$.Name}}{{.Name}}_r) Unmarshal(data []byte) error {
    return r.unpack(data)
}

// xorm encode interface
func (r *{{$.Name}}{{.Name}}_r) ToDB() ([]byte, error) {
    return r.pack()
}

// xorm decode interface
func (r *{{$.Name}}{{.Name}}_r) FromDB(data []byte) error {
    return r.unpack(data)
}

// record {{.Name}} pack
func (r *{{$.Name}}{{.Name}}_r) pack() ([]byte, error) {
    j := &{{$.Name}}{{.Name}}Json{}
    j.ColName = make([]string, {{len .Table.Cols}})
    j.ColType = make([]string, {{len .Table.Cols}})
    {{range $k, $v := .Table.Cols}}
    j.ColName[{{$k}}] = "{{$v.Name}}"
    j.ColType[{{$k}}] = "{{$v.Type}}"{{end}}

    j.Row = make([][]interface{}, len(r.Row))
	for k, row := range r.Row {
		if row == nil {
			panic("row is nil")
		}
        j.Row[k] = make([]interface{}, 0, {{len .Table.Cols}})
		j.Row[k] = append(j.Row[k] {{range .Table.Cols}},row.{{.Name}}{{end}})
	}

    return json.Marshal(j)
}

// record {{.Name}} unpack
func (r *{{$.Name}}{{.Name}}_r) unpack(data []byte) error {
    r.Row = r.data[:0]
	j := &{{$.Name}}{{.Name}}Json{}
	err := json.Unmarshal(data, j)
	if err != nil {
		return err
	}

	for _, row := range j.Row {
		if len(r.Row) > cap(r.data) {
			break
		}
		{{tolower $pname}}row := &{{$.Name}}{{.Name}}_c{}
		for k, col := range row {
			switch j.ColName[k] { {{range .Table.Cols}}
			case "{{.Name}}":
				if j.ColType[k] == "{{.Type}}" { {{if eq .Type "string"}}
					{{tolower $pname}}row.{{.Name}} = col.({{.Type}}) {{else}}
                    toolkit.ParseNumber(col, &{{tolower $pname}}row.{{.Name}}) {{end}}
				}{{end}}
			}
		}
		r.Row = r.data[:len(r.Row)+1]
		r.Row[len(row)-1] = toolboxrow
	}
	return nil
}

{{end}}
{{end}}
// {{.Name}} archive
type {{.Name}}Archive struct {
	root object.Object `xorm:"-"`
    flag int `xorm:"-"`

    Id int64 {{range .Property}} {{if eq .Save "true"}}
	{{.Name}} {{if eq .Type "tuple"}}*{{$.Name}}{{.Name}}_t `xorm:"json"`{{else if eq .Type "table"}}*{{$.Name}}{{.Name}}_r `xorm:"json"`{{else}}{{.Type}} {{if eq .Type "string"}}`xorm:"varchar({{strsize .}})"`{{end}}{{end}}  // {{.Desc}}{{end}} {{end}}
}

// {{.Name}} archive construct
func New{{.Name}}Archive(root object.Object) *{{.Name}}Archive {
    archive := &{{.Name}}Archive{root:root}
    {{range .Property}}{{if eq .Save "true"}}
    {{if eq .Type "tuple"}}archive.{{.Name}} = New{{$.Name}}{{.Name}}(root){{else if eq .Type "table"}}archive.{{.Name}} = New{{$.Name}}{{.Name}}(root){{end}}{{end}}{{end}}
    return archive
}

// archive table name
func (a *{{.Name}}Archive) TableName() string {
    return "{{.Name}}"
}

// {{.Name}} attr
type {{.Name}}Attr struct{
	root object.Object

    {{range .Property}}{{if ne .Save "true"}}
    {{.Name}} {{if eq .Type "tuple"}}{{$.Name}}{{.Name}}_t{{else if eq .Type "table"}}{{$.Name}}{{.Name}}_r{{else}}{{.Type}}{{end}} // {{.Desc}}{{end}}{{end}}
}

// {{.Name}} attr construct
func New{{.Name}}Attr(root object.Object) *{{.Name}}Attr {
    attr := &{{.Name}}Attr{root:root} 
    {{range .Property}}{{if ne .Save "true"}}
    {{if eq .Type "tuple"}}attr.{{.Name}} = New{{$.Name}}{{.Name}}(root){{else if eq .Type "table"}}attr.{{.Name}} = New{{$.Name}}{{.Name}}(root){{end}}{{end}}{{end}}
    return attr
}

// {{.Name}}
type {{.Name}} struct{
	object.ObjectWitness
    archive *{{.Name}}Archive // archive
    attr *{{.Name}}Attr // attr
}

// {{.Name}} construct
func New{{.Name}}() *{{.Name}} {
    o := &{{.Name}}{}
    o.archive = New{{.Name}}Archive(o)
    o.attr = New{{.Name}}Attr(o)
	o.Witness(o)
    return o
}

// {{.Name}} store
func (o *{{.Name}}) Store() {
}

// {{.Name}} type
func (o *{{.Name}}) Type() string {
	return "{{.Type}}"
}

// {{.Name}} entity name
func (o *{{.Name}}) Entity() string {
	return "{{.Name}}"
}

// {{.Name}} load
func (o *{{.Name}}) Load() {
}

// get archive
func (o *{{.Name}}) Archive() *{{.Name}}Archive {
    return o.archive
}

// get attr
func (o *{{.Name}}) Attr() *{{.Name}}Attr {
    return o.attr
}

{{range .Property}}
// set {{.Name}} {{.Desc}}
func (o *{{$.Name}}) Set{{.Name}}( {{tolower .Name}} {{if eq .Type "tuple"}} {{$.Name}}{{.Name}}_t{{else if eq .Type "table"}} *{{$.Name}}{{.Name}}_r {{else}} {{.Type}} {{end}}){
    {{if eq .Type "table"}} panic("{{.Name}} can't set") {{else}} {{if eq .Save "true"}}{{if eq .Type "tuple"}}if o.archive.{{.Name}}.Equal({{tolower .Name}}) {
		return 
	} 
	old := *o.archive.{{.Name}}
	*o.archive.{{.Name}} = {{tolower .Name}} {{else}} if o.archive.{{.Name}} == {{tolower .Name}} {
		return
	} 
	old := o.archive.{{.Name}}
	o.archive.{{.Name}} = {{tolower .Name}} {{end}}	{{else}} {{if eq .Type "tuple"}}if o.attr.{{.Name}}.Equal({{tolower .Name}}) {
		return 
	} 
	old := *o.attr.{{.Name}}
	*o.attr.{{.Name}} = {{tolower .Name}} {{else}} if o.attr.{{.Name}} == {{tolower .Name}} {
		return
	} 
	old := o.attr.{{.Name}}
	o.attr.{{.Name}} = {{tolower .Name}}{{end}} {{end}} {{end}}  {{if ne .Type "table"}}{{if eq .Type "tuple"}}	
	o.UpdateTuple("{{.Name}}", {{tolower .Name}}, old) {{else}}
	o.UpdateAttr("{{.Name}}", {{tolower .Name}}, old) {{end}} {{end}}
}

// get {{.Name}} {{.Desc}}
func (o *{{$.Name}}) {{.Name}}() {{if eq .Type "tuple"}} *{{$.Name}}{{.Name}}_t{{else if eq .Type "table"}} *{{$.Name}}{{.Name}}_r {{else}} {{.Type}} {{end}} {
    {{if eq .Save "true"}}return o.archive.{{.Name}}{{else}}return o.attr.{{.Name}}{{end}}
}
{{end}}

// attr type
func  (o *{{$.Name}}) GetAttrType(name string) string {
	switch name { {{range .Property}}
	case "{{.Name}}":
		return "{{.Type}}" {{end}}
	default:
		return "unknown"
	}
}

// attr expose info 
func  (o *{{$.Name}}) Expose(name string) int {
	switch name { {{range .Property}}
	case "{{.Name}}":
		return object.EXPOSE_{{if eq .Expose ""}}NONE{{else}}{{toupper .Expose}}{{end}}{{end}}
	default:
		panic("unknown")
	}
}

// get all attr name
func (o *{{$.Name}}) AllAttr() []string {
	return []string{ {{range $k, $p := .Property}}{{if ne $k 0}},{{end}} {{with $p}}"{{.Name}}" {{end}}{{end}} }
}

// get attr index by name
func (o *{{$.Name}}) AttrIndex(name string) int {
	switch name { {{range $k, $p := .Property}}{{with $p}}
	case "{{.Name}}":
		return {{$k}} {{end}}{{end}}
	default:
		return -1
	}
}

// get attr value
func  (o *{{$.Name}}) GetAttr(name string) interface{} {
	switch name { {{range .Property}}
	case "{{.Name}}":
		return {{if eq .Save "true"}}o.archive.{{.Name}} {{else}}o.attr.{{.Name}} {{end}} {{end}}
	default:
		return nil
	}
}

// set attr value
func  (o *{{$.Name}}) SetAttr(name string, value interface{}) error {
	switch name { {{range .Property}}
	case "{{.Name}}": {{if eq .Type "tuple"}}
		if v, ok := value.({{$.Name}}{{.Name}}_t); ok {
			o.Set{{.Name}}(v)
			return nil
		} {{else if eq .Type "table"}}
		if v, ok := value.(*{{$.Name}}{{.Name}}_r); ok {
			o.Set{{.Name}}(v)
			return nil
		} {{else}}
		if v, ok := value.({{.Type}}); ok {
			o.Set{{.Name}}(v)
			return nil
		} {{end}}
		return  fmt.Errorf("attr {{.Name}} type not match") {{end}}
	default:
		return fmt.Errorf("attr %s not found", name)
	}
}

// gob register
func init() {
	gob.Register(&{{.Name}}{})
	gob.Register(&{{.Name}}Archive{})
	gob.Register([]*{{.Name}}{})
	gob.Register([]*{{.Name}}Archive{})
}