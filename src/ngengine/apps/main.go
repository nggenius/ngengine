package main

import (
	"bytes"
	"encoding/gob"
	"fmt"
	"ngengine/game/gameobject/entity"
)

func main() {
	p1 := entity.NewPlayer()
	p2 := entity.NewPlayer()
	p1.SetName("sll")
	p1.SetPosXYZ(1, 2, 3)
	p1.SetOrient(5.4)
	p1.SetGroupId(1)
	p1.SetVisualRange(200)
	p1.Toolbox().AddRowValue(-1, 1, 1)
	p1.Toolbox().AddRowValue(-1, 2, 1)
	w := bytes.NewBuffer(nil)
	enc := gob.NewEncoder(w)
	enc.Encode(p1)

	p2.Toolbox().AddRowValue(-1, 3, 1)
	r := bytes.NewBuffer(w.Bytes())
	dec := gob.NewDecoder(r)
	fmt.Println(p2, p2.Toolbox())
	dec.Decode(p2)
	fmt.Println(p2, p2.Toolbox())
}
