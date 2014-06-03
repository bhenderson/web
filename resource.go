package web

import (
	"net/http"
	"strings"
)

type Resource struct {
	Prefix string
	Index,
	Show,
	Create,
	Update,
	Delete,
	NotFound,
	Handler,
	method,
	index http.Handler
}

func (rs *Resource) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	buildResource(rs)

	paths := parseComponents(r)

	if len(paths) < 1 {
		rs.NotFound.ServeHTTP(w, r)
		return
	}

	if rs.Prefix != "" && rs.Prefix != paths[0] {
		// not found
		rs.NotFound.ServeHTTP(w, r)
		return
	}

	// no resource id
	if len(paths) < 2 {
		rs.index.ServeHTTP(w, r)
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
		if rs.Handler == nil {
			rs.NotFound.ServeHTTP(w, r)
		} else {
			prefix := strings.Join([]string{"", base, id, ""}, "/")
			http.StripPrefix(prefix, rs.Handler).ServeHTTP(w, r)
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
			Get:      rs.Index,
			NotFound: rs.NotFound,
		}
	}
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
