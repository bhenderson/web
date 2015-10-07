// +build ignore

package main

import (
	"bytes"
	"go/format"
	"io/ioutil"
	"text/template"
)

var tmp = template.Must(template.New("plusser").Parse(`package web

import (
	"io"
	"net/http"
)

// AUTO GENERATED. DO NOT EDIT
// Usually, implementations of http.ResponseWriter contain extra methods, so we
// auto generate all combinations of those to keep the functionality
{{range $idx, $el := .}}
type plusser{{$idx}} interface { {{range .}}
	{{.Pkg}}.{{.Name}} {{end}}
}

type responseWriter{{$idx}} struct {
	http.ResponseWriter
	plusser{{$idx}}
}
{{end}}
func newPlusser(wr, pl http.ResponseWriter) http.ResponseWriter {
	switch x := pl.(type) { {{range $idx, $el := .}}
	case plusser{{$idx}}:
		return &responseWriter{{$idx}}{wr, x} {{end}}
	}
	return wr
}
`))

type IS struct {
	Name, Pkg string
}

type IST []IS

func main() {
	all := IST{
		{"CloseNotifier", "http"},
		{"Flusher", "http"},
		{"Hijacker", "http"},
		{"ReaderFrom", "io"},
	}

	combo := compileAll(all)

	var buf bytes.Buffer

	tmp.Execute(&buf, combo)

	out, err := format.Source(buf.Bytes())
	if err != nil {
		panic(err)
	}

	ioutil.WriteFile("plusser.go", out, 0666)
}

func compileAll(all IST) []IST {
	if len(all) == 0 {
		return []IST{}
	}

	next := compileAll(all[1:])
	cur := []IST{}

	for _, c := range next {
		cur = append(cur, append(IST{all[0]}, c...))
	}
	for _, c := range next {
		if len(c) > 1 {
			cur = append(cur, c)
		}
	}
	for _, c := range all {
		cur = append(cur, IST{c})
	}
	return cur
}
