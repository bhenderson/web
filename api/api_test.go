package api

import (
	"fmt"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestHandle(t *testing.T) {
	allowed := func(vs ...string) http.Header {
		return http.Header{"Allow": vs}
	}
	_ = allowed

	tcs := []struct {
		verb, path string
		body       io.Reader
		status     int
		result     interface{}
		headers    http.Header
		f          Handler
	}{
		{"ANY", "/returnNil", nil, 200, nil, nil, func(h H) {
			h.Path("returnNil", func(h H) {
				h.Status = 200
				h.Return(nil)
			})
		}},
		{"ANY", "/setStatus", nil, 123, "hi 123", nil, func(h H) {
			h.Path("setStatus", func(h H) {
				h.Status = 123
				h.Return("hi 123")
			})
		}},
		{"GET", "/a/b", nil, 404, nil, nil, func(h H) {
			h.Path("a", nil)
		}},
		{"GET", "/only/updates", nil, 405, nil, allowed("PUT", "POST"), func(h H) {
			h.Path("only", func(h H) {
				h.Get(func(h H) { h.Return("not this one") })
				h.Path("updates", func(h H) {
					h.Put(nil)
					h.Post(nil)
				})
			})
		}},
		{"OPTIONS", "/foo", nil, 200, nil, allowed("GET"), func(h H) {
			h.Path("foo", func(h H) {
				h.Get(nil)
			})
		}},
		{"ANY", "/auser/b/c", nil, 200, "auser", nil, func(h H) {
			h.Path(":id", func(h H) {
				h.Return(h.PathSegment)
			})
		}},
		{"ANY", "/anypath", nil, 405, nil, allowed("GET", "POST"), func(h H) {
			h.Get(nil)
			h.Post(nil)
		}},
		{"ANY", "/panics", nil, 500, "some error", nil, func(h H) {
			panic("some error")
		}},
	}

	for _, tc := range tcs {
		t.Run(fmt.Sprintf("%s %s -> %d", tc.verb, tc.path, tc.status), func(t *testing.T) {
			if tc.headers == nil {
				tc.headers = http.Header{}
			}
			if tc.result == nil {
				tc.result = http.StatusText(tc.status)
			}
			assertRequest(t, tc.verb, tc.path, tc.body, tc.status, tc.result, tc.headers, tc.f)
		})
	}
}

func assertRequest(t *testing.T, verb, path string, body io.Reader, status int, result interface{}, headers http.Header, f Handler) {
	defer func() {
		e := recover()
		assert.Nil(t, e, fmt.Sprintf("run should not panic: %#v", e))
	}()
	w := httptest.NewRecorder()
	r := httptest.NewRequest(verb, path, body)

	f.ServeHTTP(w, r)

	assert.Equal(t, status, w.Code, "status")
	assert.Equal(t, result, w.Body.String(), "result")
	assert.Equal(t, headers, w.HeaderMap)
}
