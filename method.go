package web

import "net/http"

type Method struct {
	Any,
	Delete,
	Get,
	NotFound,
	Option,
	Patch,
	Post,
	Put http.Handler
}

// ServeHTTP implements the http.Handler interface for Method.
// Explicit methods are tried first, then Any is used as a fallback.
// NotFound is used for any missing method.
// falls back on http.NotFound
func (m *Method) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	var f http.Handler

	switch r.Method {
	case "DELETE":
		f = m.Delete
	case "GET":
		f = m.Get
	case "OPTION":
		f = m.Option
	case "PATCH":
		f = m.Patch
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
		f = http.NotFoundHandler()
	}

	f.ServeHTTP(w, r)
}
