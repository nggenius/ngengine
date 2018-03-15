package object

const (
	EXPOSE_NONE  = 0
	EXPOSE_OWNER = 1
	EXPOSE_OTHER = 2
	EXPOSE_ALL   = EXPOSE_OWNER & EXPOSE_OTHER
)

type ObjectCreate interface {
	Create() interface{}
}

type AttrWitness interface {
	UpdateAttr(name, typ string, val interface{}, old interface{})
}

type TupleWitness interface {
	UpdateTuple(name string, val interface{}, old interface{})
}

type TableWitness interface {
	AddTableRow(name string, row int)
	AddTableRowValue(name string, row int, val ...interface{})
	DelTableRow(name string, row int)
	ClearTable(name string)
	ChangeTable(name string, row, col int, val interface{})
}

type Object interface {
	GetAttrType(name string) string
	GetAttr(name string) interface{}
	SetAttr(name string, value interface{}) error
	Expose(name string) int

	UpdateAttr(name, typ string, val interface{}, old interface{})
	UpdateTuple(name string, val interface{}, old interface{})
	AddTableRow(name string, row int)
	AddTableRowValue(name string, row int, val ...interface{})
	DelTableRow(name string, row int)
	ClearTable(name string)
	ChangeTable(name string, row, col int, val interface{})
}
