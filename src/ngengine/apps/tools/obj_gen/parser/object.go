package parser

import (
	"bufio"
	"encoding/xml"
	"fmt"
	"os"
	"os/exec"
	"strings"
	"text/template"

	"github.com/mysll/toolkit"
)

type Table struct {
	MaxRows int   `xml:"maxrows,attr"`
	Cols    []Col `xml:"col"`
}

type Col struct {
	Name string `xml:"name,attr"`
	Type string `xml:"type,attr"`
	Save string `xml:"save,attr"`
	Desc string `xml:"desc,attr"`
}

type Element struct {
	Name string `xml:"name,attr"`
	Type string `xml:"type,attr"`
	Desc string `xml:"desc,attr"`
}

type Property struct {
	Name   string    `xml:"name,attr"`
	Type   string    `xml:"type,attr"`
	Size   int       `xml:"size,attr"`
	Save   string    `xml:"save,attr"`
	Expose string    `xml:"expose,attr"`
	Desc   string    `xml:"desc,attr"`
	Tuple  []Element `xml:"tuple"`
	Table  Table     `xml:"table"`
}

type Object struct {
	Package  string     `xml:"package"`
	Name     string     `xml:"name"`
	Type     string     `xml:"type"`
	Include  string     `xml:"include"`
	Archive  string     `xml:"archive"`
	Property []Property `xml:"propertys>property"`
}

func StringSize(p *Property) int {
	if p.Size == 0 {
		return 255
	}
	return p.Size
}

func OutputFile(tpl, path, outfile string, obj *Object) {
	t, err := template.New(tpl).Funcs(template.FuncMap{
		"tolower": strings.ToLower,
		"toupper": strings.ToUpper,
		"strsize": StringSize,
	}).ParseFiles(path + tpl)
	if err != nil {
		fmt.Println("writer", err)
		return
	}

	//save file
	file, err := os.Create(outfile)
	if err != nil {
		fmt.Println("writer", err)
		return
	}
	defer file.Close()

	writer := bufio.NewWriter(file)

	err = t.Execute(writer, obj)
	if err != nil {
		fmt.Println("writer", err)
	}

	writer.Flush()

	cmd := exec.Command("gofmt", "--w", outfile)
	cmd.Run()
}

func ParseFromXml(file, tpl, path, outfile string) {
	obj := &Object{}
	data, err := toolkit.ReadFile(file)
	if err != nil {
		panic(err)
	}
	xml.Unmarshal(data, obj)
	OutputFile(tpl, path, outfile, obj)
}
