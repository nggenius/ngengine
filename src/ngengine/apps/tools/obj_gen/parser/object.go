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
	Include  []string   `xml:"include"`
	Archive  string     `xml:"archive"`
	Property []Property `xml:"propertys>property"`
}

func (o *Object) HasProperty(name string) bool {
	for _, p := range o.Property {
		if p.Name == name {
			return true
		}
	}
	return false
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

func ParseObjectFromXml(file, tpl string) (*Object, error) {
	obj := &Object{}
	data, err := toolkit.ReadFile(file)
	if err != nil {
		return nil, err
	}
	err = xml.Unmarshal(data, obj)
	if err != nil {
		return nil, err
	}

	// 循环遍历所有包含的文件
	if len(obj.Include) > 0 {
		for _, i := range obj.Include {
			child, err := ParseObjectFromXml(i, tpl)
			if err != nil {
				return nil, err
			}

			for _, p := range child.Property {
				if obj.HasProperty(p.Name) { // 属性可以按层级履盖
					continue
				}
				obj.Property = append(obj.Property, p)
			}
		}
	}
	return obj, err
}

func ParseFromXml(file, tpl, tpl_path, outfile string) {
	obj, err := ParseObjectFromXml(file, tpl)
	if err != nil {
		panic(err)
	}

	OutputFile(tpl, tpl_path, outfile, obj)
}
