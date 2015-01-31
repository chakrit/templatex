package main

import "os"
import "text/template"
import "text/template/parse"

const TemplateName = "main.template"

func main() {
	must := template.Must
	funcs := template.FuncMap{
		"extend": template_Extend,
	}

	t := must(template.New("").Funcs(funcs).ParseFiles(TemplateName))
	t = must(preprocess(t))
	if e := t.Execute(os.Stdout, nil); e != nil {
		panic(e)
	}
}

func template_Extend(name string) string {
	// only as a placeholder, so no-op
	return ""
}

func preprocess(t *template.Template) (*template.Template, error) {
	if t.Tree != nil {
		return preprocessTree(t, t.Tree)
	}

	result, children := t, t.Templates()
	for _, child := range children {
		if updated, e := preprocess(child); e != nil {
			return nil, e
		} else if updated != nil {
			result = updated
		}
	}

	return result, nil
}

func preprocessTree(t *template.Template, tree *parse.Tree) (*template.Template, error) {
	return preprocessNode(t, tree.Root)
}

func preprocessNode(t *template.Template, node parse.Node) (result *template.Template, e error) {
	defer func() {
		if recov := recover(); recov != nil {
			if err, ok := recov.(error); ok {
				result, e = nil, err
			}
		}
	}()

	result = t
	descend := func(child parse.Node) {
		if updated, e := preprocessNode(result, child); e != nil {
			panic(e)
		} else if updated != nil {
			result = updated
		}
	}

	switch n := node.(type) {
	case *parse.ListNode:
		for _, child := range n.Nodes {
			descend(child)
		}
	case *parse.ActionNode:
		descend(n.Pipe)
	case *parse.PipeNode:
		for _, child := range n.Cmds {
			descend(child)
		}

	case *parse.CommandNode:
		if len(n.Args) == 2 {
			if ident, ok := n.Args[0].(*parse.IdentifierNode); ok && ident.Ident == "extend" {
				if extendee, ok := n.Args[1].(*parse.StringNode); ok && len(extendee.Text) > 0 {
					result, e = preprocessExtend(result, extendee.Text)
					break
				}
			}
		}

		for _, child := range n.Args {
			descend(child)
		}

	default:
		return nil, nil // ignore all other node types.
	}

	return result, nil
}

func preprocessExtend(tmpl *template.Template, filename string) (*template.Template, error) {
	if extendee, e := template.ParseFiles(filename); e != nil {
		return nil, e
	} else if extendee != nil {
		return tmpl.AddParseTree(filename, extendee.Tree)
	}

	return tmpl, nil
}
