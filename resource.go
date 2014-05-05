package web

import (
	"net/http"
	"strings"
)

type Resource struct {
	Index    http.HandlerFunc
	Show     http.HandlerFunc
	Create   http.HandlerFunc
	Update   http.HandlerFunc
	Delete   http.HandlerFunc
	NotFound http.HandlerFunc
	Resource http.Handler
}

func (rs *Resource) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	paths := parseComponents(r)

	// no resource id
	if len(paths) < 2 {
		m := &Method{
			Get:      rs.Index,
			NotFound: rs.NotFound,
		}
		m.ServeHTTP(w, r)
		return
	}

	if len(paths) > 2 {
		if rs.Resource == nil {
			if rs.NotFound == nil {
				http.NotFound(w, r)
			} else {
				rs.NotFound(w, r)
			}
		} else {
			ps := [4]string{}
			ps[1] = paths[0]
			ps[2] = paths[1]
			prefix := strings.Join(ps[:], "/")
			h := http.StripPrefix(prefix, rs.Resource)
			h.ServeHTTP(w, r)
		}
		return
	}

	// set resource id
	id := paths[1]
	r.ParseForm()
	r.Form.Set("id", id)

	m := &Method{
		Get:      rs.Show,
		Put:      rs.Update,
		Post:     rs.Create,
		Delete:   rs.Delete,
		NotFound: rs.NotFound,
	}
	m.ServeHTTP(w, r)
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
