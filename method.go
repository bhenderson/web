package web

import "net/http"

type Method struct {
	Get      http.HandlerFunc
	Post     http.HandlerFunc
	Put      http.HandlerFunc
	Delete   http.HandlerFunc
	Option   http.HandlerFunc
	Any      http.HandlerFunc
	NotFound http.HandlerFunc
}

// ServeHTTP implements the http.Handler interface for Method.
// Explicit methods are tried first, then Any is used as a fallback.
// NotFound is used for any missing method.
// falls back on http.NotFound
func (m *Method) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	var f http.HandlerFunc

	switch r.Method {
	case "DELETE":
		f = m.Delete
	case "GET":
		f = m.Get
	case "OPTION":
		f = m.Option
	case "POST":
		f = m.Post
	case "PUT":
		f = m.Put
	}

	if f == nil {
		f = m.Any
	}

	if f == nil {
		f = m.NotFound
	}

	if f == nil {
		f = http.NotFound
	}

	f(w, r)
}
