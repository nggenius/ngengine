package main

import (
	"ngengine/apps/tools/obj_gen/parser"
)

func main() {
	parser.ParseFromXml("player.xml", "object.tpl", "./parser/", "../../../game/gameobject/entity/player.go")
}
