package main

import "os"
import "html/template"
import "text/template/parse"
import "sort"
import "fmt"

const BaseTemplateName = "base.template"
const ChildTemplateName = "main.template"

func main() {
	must := template.Must
	funcs := template.FuncMap{
		"extends": template_Extend,
	}

	t := must(template.New("").Funcs(funcs).ParseFiles(ChildTemplateName))
	t = must(preprocess(t))
	if e := t.ExecuteTemplate(os.Stdout, BaseTemplateName, nil); e != nil {
		panic(e)
	}
}

func template_Extend(name string) string {
	// only as a placeholder, so no-op
	return ""
}

func log(msg string, args ...interface{}) {
	fmt.Errorf(msg, args...)
}

func preprocess(t *template.Template) (*template.Template, error) {
	result, updated := t, false
	lookup := map[string]*template.Template{}
	extends := listExtends(t)

	lookup[result.Name()] = result
	for _, template := range result.Templates() {
		lookup[template.Name()] = template
	}

	for _, extension := range extends {
		if _, ok := lookup[extension]; !ok {
			log("extends:", extension)
			if update, e := t.ParseFiles(extension); e != nil {
				return result, e
			} else {
				result = update
				updated = true
			}
		}
	}

	if updated {
		return preprocess(result) // until we have no updates
	}

	return result, nil
}

func listExtends(t *template.Template) []string {
	log("processing:", t.Name())

	templates := []*template.Template{t}
	for _, template := range t.Templates() {
		templates = append(templates, template)
	}

	nodes := []parse.Node{}
	for _, template := range templates {
		if template.Tree != nil {
			log("root:", template.Name())
			nodes = append(nodes, template.Tree.Root)
		}
	}

	result := []string{}
	for _, node := range nodes {
		result = append(result, listNodeExtends(node)...)
	}

	sort.StringSlice(result).Sort()
	return result
}

func listNodeExtends(node parse.Node) []string {
	result := []string{}
	include := func(node parse.Node) {
		result = append(result, listNodeExtends(node)...)
	}

	switch n := node.(type) {
	case *parse.ListNode:
		for _, child := range n.Nodes {
			include(child)
		}
	case *parse.ActionNode:
		include(n.Pipe)
	case *parse.PipeNode:
		for _, child := range n.Cmds {
			include(child)
		}

	case *parse.CommandNode:
		if len(n.Args) == 2 {
			if ident, ok := n.Args[0].(*parse.IdentifierNode); ok && ident.Ident == "extends" {
				if extendee, ok := n.Args[1].(*parse.StringNode); ok && len(extendee.Text) > 0 {
					result = append(result, extendee.Text)
				}
			}
		}

		for _, child := range n.Args {
			include(child)
		}

	default:
		return []string{} // ignore all other node types.
	}

	return result
}
