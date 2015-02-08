package main

import (
	"container/list"
	"fmt"
	"html/template"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"
	// "text/template" // also works
	"text/template/parse"
)

var _ = fmt.Println

type Template struct {
	*template.Template
	wd string
}

func normalizeFilename(wd, filename string) (string, error) {
	if absName, e := filepath.Abs(filename); e != nil {
		return "", e
	} else {
		filename = absName
	}

	if strings.HasPrefix(filename, wd) {
		filename = filename[len(wd)+1:]
	}

	return filename, nil
}

func parseHTMLTemplate(filename string) (*template.Template, error) {
	funcs := template.FuncMap{
		"extends": func(name string) string { return "" }, // dummy
	}

	if b, e := ioutil.ReadFile(filename); e != nil {
		return nil, e
	} else if result, e := template.New(filename).Funcs(funcs).Parse(string(b)); e != nil {
		return nil, e
	} else {
		return result, nil
	}
}

func ParseFile(wd, filename string) (result *Template, e error) {
	if wd == "" {
		if wd, e = os.Getwd(); e != nil {
			return nil, e
		}
	}

	if wd, e = filepath.Abs(wd); e != nil {
		return nil, e
	}

	if filename, e = normalizeFilename(wd, filename); e != nil {
		return nil, e
	}

	var htmlTemplate *template.Template
	if htmlTemplate, e = parseHTMLTemplate(filename); e != nil {
		return nil, e
	}

	extensions := list.New()
	findExtensions(extensions, htmlTemplate)

	for elem := extensions.Front(); elem != nil; elem = extensions.Front() {
		extensions.Remove(elem)

		extFilename := elem.Value.(string)
		if extFilename, e = normalizeFilename(wd, elem.Value.(string)); e != nil {
			return nil, e // TODO: Should tell clue that we're parsing extFilename.
		} else if htmlTemplate.Lookup(extFilename) != nil {
			continue // ignore parsed template
		}

		var extTemplate *template.Template
		if extTemplate, e = parseHTMLTemplate(extFilename); e != nil {
			return nil, e

		} else if extTemplate.Tree != nil {
			findExtensions(extensions, extTemplate)
			htmlTemplate.AddParseTree(extTemplate.Name(), extTemplate.Tree)
			for _, inner := range extTemplate.Templates() {
				htmlTemplate.AddParseTree(inner.Name(), inner.Tree)
			}
		}
	}

	return &Template{Template: htmlTemplate, wd: wd}, nil
}

func findExtensions(result *list.List, t *template.Template) {
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
