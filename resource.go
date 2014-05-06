package web

import (
	"net/http"
	"strings"
)

type Resource struct {
	Prefix   string
	Index    http.HandlerFunc
	Show     http.HandlerFunc
	Create   http.HandlerFunc
	Update   http.HandlerFunc
	Delete   http.HandlerFunc
	NotFound http.HandlerFunc
	Handler  http.Handler
	method   *Method
}

func (rs *Resource) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	paths := parseComponents(r)

	if rs.Prefix != "" && rs.Prefix != paths[0] {
		// not found
		if rs.NotFound == nil {
			http.NotFound(w, r)
		} else {
			rs.NotFound(w, r)
		}
		return
	}

	// no resource id
	if len(paths) < 2 {
		m := &Method{
			Get:      rs.Index,
			NotFound: rs.NotFound,
		}
		m.ServeHTTP(w, r)
		return
	}

	if len(paths) > 2 && rs.Handler == nil {
		if rs.NotFound == nil {
			http.NotFound(w, r)
		} else {
			rs.NotFound(w, r)
		}
		return
	}

	// set resource id
	base := paths[0]
	id := paths[1]
	formId := "id"
	if rs.Prefix != "" {
		formId = rs.Prefix + "_id"
	}
	r.ParseForm()
	r.Form.Set(formId, id)

	if len(paths) > 2 {
		prefix := strings.Join([]string{"", base, id, ""}, "/")
		http.StripPrefix(prefix, rs.Handler).ServeHTTP(w, r)
		return
	}

	buildMethod(rs)
	rs.method.ServeHTTP(w, r)
}

func buildMethod(rs *Resource) {
	if rs.method == nil {
		rs.method = &Method{
			Get:      rs.Show,
			Put:      rs.Update,
			Post:     rs.Create,
			Delete:   rs.Delete,
			NotFound: rs.NotFound,
		}
	}
}

func parseComponents(r *http.Request) []string {
	path := r.URL.Path
	path = strings.TrimSpace(path)
	//Cut off the leading and trailing forward slashes, if they exist.
	//This cuts off the leading forward slash.
	if strings.HasPrefix(path, "/") {
		path = path[1:]
	}
	//This cuts off the trailing forward slash.
	if strings.HasSuffix(path, "/") {
		path = path[:len(path)-1]
	}
	//We need to isolate the individual components of the path.
	components := strings.Split(path, "/")
	return components
}
