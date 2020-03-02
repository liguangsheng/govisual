package govisual

import (
	"bytes"
	"fmt"
	"text/template"
)

var defaultTemplate = `digraph "{{.v.Config.Module}}" {
	label="{{.v.Config.Module}}";
	node [shape=Mrecord, style=solid];
{{range $key, $node := .v.Digraph.Nodes}}	"{{$node.Value.FullName}}" [label="{{nodeLabel $node}}", {{if eq $node.Out 0}}shape="box",{{end}} color="{{color $node.Value.Type}}", fontcolor="{{color $node.Value.Type}}"];
{{end}}{{range $key, $edge := .v.Digraph.Edges}}	"{{$edge.From.Value.FullName}}" -> "{{$edge.To.Value.FullName}}";
{{end}}
}`

func Render(v *Visualization, filenames ...string) string {
	var t = template.New("dot").Funcs(map[string]interface{}{
		"color":     color,
		"nodeLabel": nodeLabel,
	})

	if len(filenames) > 0 {
		t = template.Must(t.ParseFiles(filenames...))
	} else {
		t = template.Must(t.Parse(defaultTemplate))
	}

	var buf bytes.Buffer
	if err := t.Execute(&buf, map[string]interface{}{
		"v": v,
	}); err != nil {
		return err.Error()
	}
	return buf.String()
}

func color(importType ImportType) string {
	switch importType {
	case Self:
		return "black"
	case Org:
		return "orange"
	case Third:
		return "blue"
	case Sys:
		return "green"
	}
	return "black"
}

func nodeLabel(n *Node) string {
	pkg := n.Value.(*Package)

	var typ string
	switch pkg.Type {
	case Self:
		typ = "self"
	case Org:
		typ = "org"
	case Third:
		typ = "third"
	case Sys:
		typ = "sys"
	}

	return fmt.Sprintf("%s:%s:I%d:O%d", pkg.SimpleName, typ, n.In, n.Out)
}
