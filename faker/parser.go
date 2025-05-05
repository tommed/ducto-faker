package faker

import (
	"fmt"
	"io/ioutil"
	"strings"
	"text/template"
	"text/template/parse"
)

// Placeholder represents a single interpolation tag in a template.
type Placeholder struct {
	FieldName string
	Type      string
	Params    map[string]string
}

// ParseTemplate reads a template file and extracts all inline-modifier placeholders using {{ }} with quoted strings.
func ParseTemplate(path string) ([]Placeholder, error) {
	data, err := ioutil.ReadFile(path)
	if err != nil {
		return nil, fmt.Errorf("failed to read template %s: %w", path, err)
	}

	t, err := template.New("tmpl").Parse(string(data))
	if err != nil {
		return nil, fmt.Errorf("failed to parse template %s: %w", path, err)
	}

	var out []Placeholder
	if t.Tree == nil || t.Tree.Root == nil {
		return out, nil
	}
	walkNode(t.Tree.Root, &out)
	return out, nil
}

func walkNode(node parse.Node, out *[]Placeholder) {
	switch n := node.(type) {
	case *parse.ListNode:
		for _, child := range n.Nodes {
			walkNode(child, out)
		}

	case *parse.ActionNode:
		if n.Pipe != nil {
			for _, cmd := range n.Pipe.Cmds {
				for _, arg := range cmd.Args {
					if str, ok := arg.(*parse.StringNode); ok {
						raw := strings.TrimSpace(str.Text)
						parts := strings.Split(raw, ":")
						if len(parts) >= 2 {
							*out = append(*out, Placeholder{
								FieldName: parts[0],
								Type:      parts[1],
								Params:    parseParams(parts[2:]),
							})
						}
					}
				}
			}
		}

	case *parse.IfNode:
		walkNode(n.List, out)
		if n.ElseList != nil {
			walkNode(n.ElseList, out)
		}

	case *parse.RangeNode:
		walkNode(n.List, out)
		if n.ElseList != nil {
			walkNode(n.ElseList, out)
		}

	case *parse.WithNode:
		walkNode(n.List, out)
		if n.ElseList != nil {
			walkNode(n.ElseList, out)
		}
	}
}

func parseParams(items []string) map[string]string {
	m := make(map[string]string, len(items))
	for _, kv := range items {
		parts := strings.SplitN(kv, "=", 2)
		if len(parts) == 2 {
			m[parts[0]] = parts[1]
		}
	}
	return m
}
