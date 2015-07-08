package web

import (
	"net/http"
	"strings"
)

// ResourceHandleFunc is a function that will be given the second path element.
// ResourceHandleFunc is itself an http.Handler.
type ResourceHandleFunc func(string) http.Handler

func (rhf ResourceHandleFunc) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	paths := PathParts(r.URL.Path)
	var id string
	if len(paths) > 1 {
		id = paths[1]
	}
	rhf(id).ServeHTTP(w, r)
}

type Resource struct {
	// Handlers for a resource.
	// Index   -> GET /prefix/
	// Create  -> POST /prefix/
	// Show    -> GET /prefix/:id
	// Update  -> PATCH /prefix/:id
	// Replace -> PUT /prefix/:id
	// Delete  -> DELETE /prefix/:id
	Index,
	Create,

	Show,
	Replace,
	Update,
	Delete http.Handler

	// NotFound is called when given an invalid path
	NotFound,

	// MethodNotAllowed is called when given an invalid HTTP Method for a path
	MethodNotAllowed,

	// Handler is called for /prefix/:id/
	Handler,

	method,
	index http.Handler
}

func (rs *Resource) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	buildResource(rs)

	paths := PathParts(r.URL.Path)

	// no resource id
	if len(paths) < 2 {
		rs.index.ServeHTTP(w, r)
		return
	}

	if len(paths) > 2 {
		if rs.Handler == nil {
			rs.NotFound.ServeHTTP(w, r)
		} else {
			rs.Handler.ServeHTTP(w, r)
		}
		return
	}

	rs.method.ServeHTTP(w, r)
}

func buildResource(rs *Resource) {
	if rs.NotFound == nil {
		rs.NotFound = http.NotFoundHandler()
	}
	if rs.index == nil {
		rs.index = &Method{
			Get:              rs.Index,
			Post:             rs.Create,
			MethodNotAllowed: rs.MethodNotAllowed,
		}
	}
	if rs.method == nil {
		rs.method = &Method{
			Get:              rs.Show,
			Patch:            rs.Update,
			Put:              rs.Replace,
			Delete:           rs.Delete,
			MethodNotAllowed: rs.MethodNotAllowed,
		}
	}
}

// PathParts removes leading and trailing slash, then splits on slash
func PathParts(path string) []string {
	// removing leading /
	if path[0] == '/' {
		path = path[1:]
	}
	// remove trailing /
	if path[len(path)-1] == '/' {
		path = path[:len(path)-1]
	}
	return strings.Split(path, "/")
}
