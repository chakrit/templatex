package main

import (
	"fmt"
	"os"
)

const BaseTemplateName = "templates/base.template"
const ChildTemplatePath = "templates/subfolder/main.template"

func main() {
	_ = os.Stdout

	if template, e := ParseFile("", ChildTemplatePath); e != nil {
		fmt.Println("ParseFile:", e.Error())

	} else if e := template.ExecuteTemplate(os.Stdout, BaseTemplateName, nil); e != nil {
		fmt.Println("ExecuteTemplate:", e.Error())
	}
}
