package api

import (
	"fmt"
	"io"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestHandleStatus(t *testing.T) {
	t.Run("GET_foo", func(t *testing.T) {
		assertRequest(t,
			"GET", "/foo", nil,
			200, "",
			func(h H) {
				h.PathEnd("foo", func(h H) {
					h.Status = 200
					h.Return(nil)
				})
			},
		)
	})

	t.Run("GET_bar", func(t *testing.T) {
		assertRequest(t,
			"GET", "/bar", nil,
			404, "Not Found",
			func(h H) {
				h.PathEnd("foo", func(h H) {
					h.Status = 200
					h.Return(nil)
				})
			},
		)
	})

	t.Run("GET_not_allowed", func(t *testing.T) {
		var allowed string
		assertRequest(t,
			"GET", "/foo", nil,
			405, "Method Not Allowed",
			func(h H) {
				defer h.Catch(func(h H) {
					allowed = h.Header().Get("Allow")
				})
				h.PathEnd("foo", func(h H) {
					h.Allow("PUT, POST")
				})
			},
		)
		assert.Equal(t, "PUT, POST", allowed)
	})

	t.Run("Auto_405", func(t *testing.T) {
		var allowed string
		assertRequest(t,
			"GET", "/foo", nil,
			405, "Method Not Allowed",
			func(h H) {
				defer h.Catch(func(h H) {
					allowed = h.Header().Get("Allow")
				})
				h.Delete(nil) // allowed verbs reset within path
				h.Path("foo", func(h H) {
					h.Put(nil)
					h.Post(nil)
				})
			},
		)
		assert.Equal(t, "PUT, POST", allowed)
	})
}

func assertRequest(t *testing.T, verb, path string, body io.Reader, status int, result interface{}, f Handler) {
	defer func() {
		e := recover()
		assert.Nil(t, e, fmt.Sprintf("should not be nil: %#v", e))
	}()
	w := httptest.NewRecorder()
	r := httptest.NewRequest(verb, path, body)

	Run(f).ServeHTTP(w, r)
	assert.Equal(t, status, w.Code, "status")
	assert.Equal(t, result, w.Body.String(), "result")
}
