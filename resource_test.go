package web

import (
	"fmt"
	"github.com/stretchr/testify/assert"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestResource(t *testing.T) {
	handler := func(n string) http.HandlerFunc {
		return func(w http.ResponseWriter, r *http.Request) {
			fmt.Fprintf(w, "%s %s", r.Method, n)
		}
	}
	rs := &Resource{
		Index:   handler("index"),
		Show:    handler("show"),
		Create:  handler("create"),
		Replace: handler("replace"),
		Update:  handler("update"),
		Delete:  handler("delete"),
	}

	var body string
	var c int

	_, body = testResource(t, "GET", "/users/", rs)
	assert.Equal(t, "GET index", body)

	_, body = testResource(t, "POST", "/users/", rs)
	assert.Equal(t, "POST create", body)

	_, body = testResource(t, "GET", "/users/a", rs)
	assert.Equal(t, "GET show", body)

	_, body = testResource(t, "PATCH", "/users/a", rs)
	assert.Equal(t, "PATCH update", body)

	_, body = testResource(t, "PUT", "/users/a", rs)
	assert.Equal(t, "PUT replace", body)

	_, body = testResource(t, "DELETE", "/users/a", rs)
	assert.Equal(t, "DELETE delete", body)

	// route not found
	c, _ = testResource(t, "GET", "/users/a/b", rs)
	assert.Equal(t, 404, c)

	// wrong method
	c, _ = testResource(t, "POST", "/users/a", rs)
	assert.Equal(t, 405, c)
}

func TestParseComponents(t *testing.T) {
	tests := []struct {
		path string
		exp  []string
	}{
		{
			"/abc",
			[]string{"abc"},
		},
		{
			"/a/bc/cde",
			[]string{"a", "bc", "cde"},
		},
		{
			"/a/bc/cde/",
			[]string{"a", "bc", "cde"},
		},
	}

	for _, test := range tests {
		paths := PathParts(test.path)
		assert.Equal(t, test.exp, paths)
	}
}

func testResource(t testing.TB, method, path string, f http.Handler) (int, string) {
	req, err := http.NewRequest(method, "http://example.com"+path, nil)
	if err != nil {
		t.Fatal(err)
	}
	w := httptest.NewRecorder()

	mux := http.NewServeMux()
	mux.Handle("/users/", f)
	mux.ServeHTTP(w, req)

	return w.Code, w.Body.String()
}
