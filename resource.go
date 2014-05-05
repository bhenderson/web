package web

import (
	"net/http"
	"strings"
)

// TODO rename
var resourceId = "_resourceId"

// TODO needed?
func ResourceId(r *http.Request) string {
	r.Form.Get(resourceId)
}

type Resource struct {
	Index    http.HandlerFunc
	Show     http.HandlerFunc
	Update   http.HandlerFunc
	Create   http.HandlerFunc
	Delete   http.HandlerFunc
	NotFound http.HandlerFunc
	method   *Method
}

func (rs *Resource) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	// skip first component
	paths := parseComponents(r)[1:]

	// no resource id
	if len(paths) < 1 {
		m := &Method{
			Get:      rs.Index,
			NotFound: rs.NotFound,
		}
		m.ServeHTTP(w, r)
		return
	}

	if len(paths) > 1 {
		if rs.NotFound == nil {
			http.NotFound(w, r)
		} else {
			rs.NotFound(w, r)
		}
		return
	}

	// set resource id
	id := paths[0]
	formId := resourceId
	r.ParseForm()
	r.Form.Set(formId, id)

	rs.setMethod()
	rs.method.ServeHTTP(w, r)
}

func (rs *Resource) setMethod() {
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
		cut_off_last_char_len := len(path) - 1
		path = path[:cut_off_last_char_len]
	}
	//We need to isolate the individual components of the path.
	components := strings.Split(path, "/")
	return components
}
