package main

import (
	"container/list"
	"html/template"
	// "text/template" // also works
	"text/template/parse"
)

func ParseFilesWithExtends(filenames ...string) (result *template.Template, e error) {
	funcs := template.FuncMap{
		"extends": func(name string) string { return "" }, // dummy
	}

	if result, e = template.New("").Funcs(funcs).ParseFiles(filenames...); e != nil {
		return nil, e
	} else if result, e = preprocess(result); e != nil {
		return nil, e
	}

	return result, nil
}

func preprocess(t *template.Template) (*template.Template, error) {
	result, extends := t, list.New()

	findExtends(extends, result)
	for elem := extends.Front(); elem != nil; elem = extends.Front() {
		extends.Remove(elem)

		extension := elem.Value.(string)
		if result.Lookup(extension) != nil {
			continue
		}

		if update, e := t.ParseFiles(extension); e != nil {
			return result, e

		} else {
			result = update
			findExtends(extends, result.Lookup(extension))
		}
	}

	return result, nil
}

func findExtends(result *list.List, t *template.Template) {
	templates := []*template.Template{t}
	for _, template := range t.Templates() {
		templates = append(templates, template)
	}

	nodes := []parse.Node{}
	for _, template := range templates {
		if template.Tree != nil {
			nodes = append(nodes, template.Tree.Root)
		}
	}

	for _, node := range nodes {
		findNodeExtends(result, node)
	}
}

func findNodeExtends(result *list.List, node parse.Node) {
	switch n := node.(type) {
	case *parse.ListNode:
		for _, child := range n.Nodes {
			findNodeExtends(result, child)
		}
	case *parse.ActionNode:
		findNodeExtends(result, n.Pipe)
	case *parse.PipeNode:
		for _, child := range n.Cmds {
			findNodeExtends(result, child)
		}

	case *parse.CommandNode:
		if len(n.Args) == 2 {
			if ident, ok := n.Args[0].(*parse.IdentifierNode); ok && ident.Ident == "extends" {
				if extendee, ok := n.Args[1].(*parse.StringNode); ok && len(extendee.Text) > 0 {
					result.PushFront(extendee.Text)
				}
			}
		}

		for _, child := range n.Args {
			findNodeExtends(result, child)
		}

	default:
		// ignore all other node types.
	}
}
