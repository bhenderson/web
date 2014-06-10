package web

import (
	"net/http"
)

type ResourceHandleFunc func(string) http.Handler

func (rhf ResourceHandleFunc) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	paths := PathParts(r)
	var id string
	if len(paths) > 1 {
		id = paths[1]
	}
	rhf(id).ServeHTTP(w, r)
}

type Resource struct {
	// FormID (default "id") is the resource id key to be set in the params.
	FormID string

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

	paths := PathParts(r)

	if len(paths) < 1 {
		rs.NotFound.ServeHTTP(w, r)
		return
	}

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

	r.ParseForm()
	r.Form.Set(rs.FormID, paths[1])

	rs.method.ServeHTTP(w, r)
}

func buildResource(rs *Resource) {
	if rs.FormID == "" {
		rs.FormID = "id"
	}
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

func PathParts(r *http.Request) map[int]string {
	path := r.URL.Path
	paths := make(map[int]string)

	if path[0] == '/' {
		path = path[1:]
	}
	for i := 0; i < len(path); i++ {
		if path[i] == '/' {
			paths[len(paths)] = path[:i]
			path = path[i+1:]
			i = -1
		}
	}
	if len(path) > 0 {
		paths[len(paths)] = path
	}
	return paths
}
