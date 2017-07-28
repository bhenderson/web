// +build ignore

package main

// Rationale
// https://stackoverflow.com/questions/34408808/golang-embed-an-interface-with-additional-hidden-methods

import (
	"bytes"
	"fmt"
	"go/format"
	"io"
	"io/ioutil"
	"os"
	"text/template"
)

var all = []string{
	"http.CloseNotifier",
	"http.Flusher",
	"http.Hijacker",
	"io.ReaderFrom",
}

var tmp = template.Must(template.New("plusser").Parse(`package web

import (
	"io"
	"net/http"
)

// AUTO GENERATED. DO NOT EDIT
// Usually, implementations of http.ResponseWriter contain extra methods, so we
// auto generate all combinations of those to keep the functionality
{{ range . }}
type plusser{{.Bits}} interface { {{range .Interfaces }}
	{{if .Off}}// {{end}}{{.Interface}} {{end}}
}

type responseWriter{{.Bits}} struct {
	http.ResponseWriter
	plusser{{.Bits}}
}
{{end}}
func newPlusser(wr, pl http.ResponseWriter) http.ResponseWriter {
	switch x := pl.(type) { {{range .}}
	case plusser{{.Bits}}:
		return &responseWriter{{.Bits}}{wr, x} {{end}}
	}
	return wr
}
`))

type Interfaces struct {
	Bit        int
	Interfaces []Interface
}

func (ifs Interfaces) Bits() string {
	return fmt.Sprintf("%04b", ifs.Bit)
}

type Interface struct {
	Bit       int
	Interface string
}

func (ifs Interface) Off() bool {
	return ifs.Bit == 0
}

func main() {
	combo := compileAll(all)

	var buf bytes.Buffer

	err := tmp.Execute(&buf, combo)
	if err != nil {
		panic(err)
	}

	out, err := format.Source(buf.Bytes())
	if err != nil {
		io.Copy(os.Stdout, &buf)
		panic(err)
	}

	ioutil.WriteFile("plusser.go", out, 0666)
}

func compileAll(ifs []string) []Interfaces {
	all := make([]Interfaces, 1<<uint(len(ifs)))

	for i := range all {
		b := len(all) - i - 1
		all[i].Bit = b
		all[i].Interfaces = make([]Interface, len(ifs))
		for j := range ifs {
			all[i].Interfaces[j].Bit = b & (1 << uint(len(ifs)-j-1))
			all[i].Interfaces[j].Interface = ifs[j]
		}
	}

	return all
}
