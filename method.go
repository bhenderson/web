package web

import (
	"net/http"
	"strings"
)

// Method is an http.Handler that routes requests according to the HTTP verb.
type Method struct {
	// GetAsHead sets the Get handler to be used for HEAD requests if the Head
	// handler is not defined.
	GetAsHead bool

	// Optional Method Handlers
	Delete,
	Get,
	Head,
	Options,
	Patch,
	Post,
	Put http.Handler

	// Any handler will respond to any method.
	Any http.Handler

	// MethodNotAllowed handler is called if the appropriate method handler is
	// not defined. Defaults to list of defined methods.
	MethodNotAllowed http.Handler
}

// ServeHTTP implements the http.Handler interface for Method.
// Explicit methods are tried first, then Any is used as a fallback.
// MethodNotAllowed is used for any missing method with defaults.
func (m *Method) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	var f http.Handler

	switch r.Method {
	case "DELETE":
		f = m.Delete
	case "HEAD":
		f = m.Head
		if f == nil && m.GetAsHead {
			f = m.Get
		}
	case "GET":
		f = m.Get
	case "PATCH":
		f = m.Patch
	case "POST":
		f = m.Post
	case "PUT":
		f = m.Put
	case "OPTIONS":
		f = m.Options
		if f == nil {
			// The default OPTIONS response.
			// rfc2616 9.2
			allowed := allowedMethods(m)
			setAllowHeader(w, allowed...)
			// empty body
			return
		}
	}

	if f == nil {
		f = m.Any
	}

	if f == nil {
		f = m.MethodNotAllowed
	}

	if f == nil {
		allowed := allowedMethods(m)
		MethodNotAllowed(w, r, allowed...)
		return
	}

	f.ServeHTTP(w, r)
}

// MethodNotAllowed replies to the request with an HTTP 405 method not allowed
// error. An optional list of allowed methods will be set in the Allow header.
func MethodNotAllowed(w http.ResponseWriter, r *http.Request, allowed ...string) {
	// rfc2616 14.7
	setAllowHeader(w, allowed...)
	http.Error(w, "405 method not allowed", http.StatusMethodNotAllowed)
}

// MethodNotAllowedHandler returns a http.Handler with a list of allowed
// methods set.
func MethodNotAllowedHandler(allowed ...string) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		MethodNotAllowed(w, r, allowed...)
	})
}

func setAllowHeader(w http.ResponseWriter, allowed ...string) {
	w.Header().Set("Allow", strings.Join(allowed, ", "))
}

func allowedMethods(m *Method) []string {
	a := make([]string, 0, 6)

	// not including OPTIONS for now, but open to discussion.
	if m.Delete != nil {
		a = append(a, "DELETE")
	}
	if m.Get != nil {
		a = append(a, "GET")
		if m.Head == nil && m.GetAsHead {
			a = append(a, "HEAD")
		}
	}
	if m.Head != nil {
		a = append(a, "HEAD")
	}
	if m.Patch != nil {
		a = append(a, "PATCH")
	}
	if m.Post != nil {
		a = append(a, "POST")
	}
	if m.Put != nil {
		a = append(a, "PUT")
	}
	return a
}
