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
		Index:  handler("index"),
		Show:   handler("show"),
		Create: handler("create"),
		Update: handler("update"),
		Delete: handler("delete"),
	}

	var body string
	var c int

	_, body = testResource(t, "GET", "/users/", rs)
	assert.Equal(t, "GET index", body)

	_, body = testResource(t, "GET", "/users/a", rs)
	assert.Equal(t, "GET show", body)

	_, body = testResource(t, "POST", "/users/a", rs)
	assert.Equal(t, "POST create", body)

	_, body = testResource(t, "PUT", "/users/a", rs)
	assert.Equal(t, "PUT update", body)

	_, body = testResource(t, "DELETE", "/users/a", rs)
	assert.Equal(t, "DELETE delete", body)

	c, _ = testResource(t, "GET", "/users/a/b", rs)
	assert.Equal(t, 404, c)

	c, _ = testResource(t, "POST", "/users/", rs)
	assert.Equal(t, 404, c)
}

func TestResource_Form(t *testing.T) {
	handler := func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintf(w, "%s %s", r.Form.Get("id"), r.Form.Get("foo"))
	}
	rs := &Resource{Show: handler}

	_, body := testResource(t, "GET", "/users/joe.smith?id=1&foo=bar", rs)
	assert.Equal(t, "joe.smith bar", body)

	c, _ := testResource(t, "POST", "/users/1", rs)
	assert.Equal(t, 404, c)
}

func TestResource_nesting(t *testing.T) {
	handler := func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintf(w, "%s %s %s",
			r.Form.Get("id"),
			r.Form.Get("posts_id"),
			r.URL.Path)
	}
	rs := &Resource{
		Handler: &Resource{
			Prefix: "posts",
			Show:   handler,
		},
	}

	_, body := testResource(t, "GET", "/users/1/posts/2", rs)
	assert.Equal(t, "1 2 posts/2", body)

	c, _ := testResource(t, "GET", "/users/1/foo/2", rs)
	assert.Equal(t, 404, c)
}

func testResource(t testing.TB, method, path string, f http.Handler) (int, string) {
	req, err := http.NewRequest(method, "http://example.com"+path, nil)
	if err != nil {
		t.Fatal(err)
	}

	w := httptest.NewRecorder()
	f.ServeHTTP(w, req)

	return w.Code, w.Body.String()
}
