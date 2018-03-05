package entity

type Witness interface {
	GetAttrWitness() AttrWitness
	GetTupleWitness() TupleWitness
	GetTableWitness() TableWitness
}

type AttrWitness interface {
	Update(name, typ string, val interface{})
}

type TupleWitness interface {
	Update(name string, val interface{})
}

type TableWitness interface {
	AddRow(name string, row int)
	AddRowValue(name string, row int, val ...interface{})
	DelRow(name string, row int)
	Clear(name string)
	Change(name string, row, col int, val interface{})
}
