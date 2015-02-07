package main

import (
	"os"
)

const BaseTemplateName = "base.template"
const ChildTemplateName = "main.template"

func main() {
	if template, e := ParseFilesWithExtends(ChildTemplateName); e != nil {
		panic(e)
	} else if e := template.ExecuteTemplate(os.Stdout, BaseTemplateName, nil); e != nil {
		panic(e)
	}
}
