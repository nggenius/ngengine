package object

type ObjectWitness struct {
}

func (o *ObjectWitness) UpdateAttr(name, typ string, val interface{}, old interface{}) {
}

func (o *ObjectWitness) UpdateTuple(name string, val interface{}, old interface{}) {
}

func (o *ObjectWitness) AddTableRow(name string, row int) {

}

func (o *ObjectWitness) AddTableRowValue(name string, row int, val ...interface{}) {

}

func (o *ObjectWitness) DelTableRow(name string, row int) {

}

func (o *ObjectWitness) ClearTable(name string) {

}

func (o *ObjectWitness) ChangeTable(name string, row, col int, val interface{}) {

}
