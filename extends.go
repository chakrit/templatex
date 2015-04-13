package templatex

import (
	"container/list"
	"fmt"
	"html/template"
	"io"
	"io/ioutil"
	"os"
	"path/filepath"
	// "text/template" // also works
	"text/template/parse"
)

var _ = fmt.Println

type Template struct {
	*template.Template
	wd string
}

func Must(t *Template, e error) *Template {
	if e != nil {
		panic(e)
	}

	return t
}

func (t *Template) ExecuteTemplate(wr io.Writer, name string, data interface{}) error {
	if filename, e := normalizeFilename(t.wd, name); e != nil {
		return e

	} else {
		return t.Template.ExecuteTemplate(wr, filename, data)
	}
}

func normalizeFilename(wd, filename string) (result string, e error) {
	if wd, e = filepath.Abs(wd); e != nil {
		return filename, e
	}

	filename = filepath.Join(wd, filename)
	if filename, e = filepath.Abs(filename); e != nil {
		return filename, e
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
	if e = findExtensions(extensions, htmlTemplate); e != nil {
		return nil, e
	}

	for elem := extensions.Front(); elem != nil; elem = extensions.Front() {
		extensions.Remove(elem)

		extFilename := elem.Value.(string)
		if htmlTemplate.Lookup(extFilename) != nil {
			continue // ignore parsed template
		}

		var extTemplate *template.Template
		if extTemplate, e = parseHTMLTemplate(extFilename); e != nil {
			return nil, e

		} else if extTemplate.Tree != nil {
			if e = findExtensions(extensions, extTemplate); e != nil {
				return nil, e
			}

			htmlTemplate.AddParseTree(extTemplate.Name(), extTemplate.Tree)
			for _, inner := range extTemplate.Templates() {
				htmlTemplate.AddParseTree(inner.Name(), inner.Tree)
			}
		}
	}

	return &Template{Template: htmlTemplate, wd: wd}, nil
}

func findExtensions(result *list.List, t *template.Template) error {
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

	wd := filepath.Dir(t.Name())
	for _, node := range nodes {
		if e := findNodeExtends(result, node, wd); e != nil {
			return e
		}
	}

	return nil
}

func findNodeExtends(result *list.List, node parse.Node, wd string) error {
	switch n := node.(type) {
	case *parse.ListNode:
		for _, child := range n.Nodes {
			if e := findNodeExtends(result, child, wd); e != nil {
				return e
			}
		}
	case *parse.ActionNode:
		if e := findNodeExtends(result, n.Pipe, wd); e != nil {
			return e
		}
	case *parse.PipeNode:
		for _, child := range n.Cmds {
			if e := findNodeExtends(result, child, wd); e != nil {
				return e
			}
		}

	case *parse.CommandNode:
		if len(n.Args) == 2 {
			if ident, ok := n.Args[0].(*parse.IdentifierNode); ok && ident.Ident == "extends" {
				if extendee, ok := n.Args[1].(*parse.StringNode); ok && len(extendee.Text) > 0 {
					if extFilename, e := normalizeFilename(wd, extendee.Text); e != nil {
						return e

					} else {
						result.PushFront(extFilename)
						return nil
					}
				}
			}
		}

		for _, child := range n.Args {
			if e := findNodeExtends(result, child, wd); e != nil {
				return e
			}
		}

	default:
		// ignore all other node types.
	}

	return nil
}
