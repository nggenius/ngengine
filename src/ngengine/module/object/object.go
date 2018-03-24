package object

const (
	EXPOSE_NONE  = 0
	EXPOSE_OWNER = 1
	EXPOSE_OTHER = 2
	EXPOSE_ALL   = EXPOSE_OWNER & EXPOSE_OTHER
)

// 对象创建接口
type ObjectCreate interface {
	Create() interface{}
}

type Object interface {
	// 沉默状态
	Silence() bool
	// 设置沉默状态
	SetSilence(bool)
	// 所属的工厂
	Factory() *Factory
	// 设置工厂，由工厂主动调用
	SetFactory(f *Factory)
	// 类型(对应xml里面的type)
	Type() string
	// entity 类型(对应xml里面的name)
	Entity() string
	// 获取某个属性的类型
	GetAttrType(name string) string
	// 获取属性
	GetAttr(name string) interface{}
	// 设置属性
	SetAttr(name string, value interface{}) error
	// 导出状态
	Expose(name string) int
	// 所有属性名
	AllAttr() []string
	// 属性的索引编号
	AttrIndex(name string) int
	// 增加一个属性观察者
	AddAttrObserver(name string, observer attrObserver) error
	// 增加表格观察者
	AddTableObserver(name string, observer tableObserver) error
	// 属性变动回调
	UpdateAttr(name string, val interface{}, old interface{})
	// tuple变动回调
	UpdateTuple(name string, val interface{}, old interface{})
	// 表格增加行回调
	AddTableRow(name string, row int)
	// 表格增加行并设置值回调
	AddTableRowValue(name string, row int, val ...interface{})
	// 设置表格行
	SetTableRowValue(name string, row int, val ...interface{})
	// 删除表格行
	DelTableRow(name string, row int)
	// 清除表格
	ClearTable(name string)
	// 表格单元格变动
	ChangeTable(name string, row, col int, val interface{})
}
