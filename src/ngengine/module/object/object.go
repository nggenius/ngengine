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

type Object interface {
	Factory() *Factory
	SetFactory(f *Factory)
	GetAttrType(name string) string
	GetAttr(name string) interface{}
	SetAttr(name string, value interface{}) error
	Expose(name string) int
	AllAttr() []string
	AttrIndex(name string) int
	AddAttrObserver(name string, observer attrObserver) error
	AddTableObserver(name string, observer tableObserver) error
	UpdateAttr(name string, val interface{}, old interface{})
	UpdateTuple(name string, val interface{}, old interface{})
	AddTableRow(name string, row int)
	AddTableRowValue(name string, row int, val ...interface{})
	DelTableRow(name string, row int)
	ClearTable(name string)
	ChangeTable(name string, row, col int, val interface{})
}
