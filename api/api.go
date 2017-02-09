package api

import (
	"fmt"
	"io"
	"net/http"
	"time"
)

func Run(f Handler) Handler {
	return f
}

var DefaultMiddleware = []Middleware{}

func Use(ms ...Middleware) {
	DefaultMiddleware = append(DefaultMiddleware, ms...)
}

func newH(w http.ResponseWriter, r *http.Request) H {
	path, rest := nextPathSegment(r.URL.Path)
	h := H{
		Request:  r,
		Response: NewResponse(w),
		Time:     time.Now(),
		path:     path,
		rest:     rest,
	}
	return h
}

type Handler func(H)

func (f Handler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	h := newH(w, r)

	h.Use(handlePanic)
	h.Use(DefaultMiddleware...)
	h.Use(
		handleFinish,
		HandleStatus(http.StatusNotFound),
	)

	h.Handle(f)
}

type Middleware func(Handler) Handler

type H struct {
	*http.Request
	*Response

	Body interface{}

	path, rest string
	Middleware []Middleware

	Time time.Time
}

func (h H) Stream(v interface{}) {
	switch x := v.(type) {
	case apiResponse:
		h.Status = x.status
		h.Stream(x.body)
	case http.Handler:
		x.ServeHTTP(h, h.Request)
	case int:
		h.Status = x
		h.WriteString(http.StatusText(x))
	case []byte:
		h.Write(x)
	case string:
		h.WriteString(x)
	case io.Reader:
		io.Copy(h, x)
	case nil:
		// noop
	default:
		fmt.Fprintf(h, "%s", x)
	}
}

func (h H) Handle(f Handler) {
	// compile in reverse order
	for i := len(h.Middleware); i > 0; i-- {
		f = h.Middleware[i-1](f)
	}
	h.Middleware = h.Middleware[:0]

	h.Header().Del(allowHeader)

	f(h)
}

func (h H) HandleHTTP(f http.Handler) {
	// TODO
	h.Handle(func(h H) {
		f.ServeHTTP(h, h.Request)
		h.Return(nil)
	})
}

// ServeFiles is exactly like http.FileServer, except that it runs
// h.Middleware, and also wraps root such that errors returned by Open panic
func (h H) ServeFiles(root http.FileSystem) {
	h.HandleHTTP(http.FileServer(apiDir{root}))
}

func nextPathSegment(path string) (string, string) {
	for i := 0; i < len(path); i++ {
		if path[i] == '/' {
			return path[:i], path[i+1:]
		}
	}
	return path, ""
}

func (h *H) Use(ms ...Middleware) {
	h.Middleware = append(h.Middleware, ms...)
}

func (h H) UseBefore(ms ...Middleware) {
	h.Middleware = append(ms, h.Middleware...)
}

func (h H) PathSegment() string {
	return h.path
}

// Path runs f if the next PathSegment matches path.
// The first character of path can have special meaning.
//
//	: will match any segment
//	* will match the rest of the request path
func (h H) Path(path string, f Handler) {
	h.path, h.rest = nextPathSegment(h.rest)
	h.checkPath(h.path, path, f)
}

// PathEnd runs only when path matches the end of the request path
// It is a convenience function for
//
//	h.Path(path, func(h H) {
//		h.Path("", func(h H) {
//			f(h)
//		}
//	})
func (h H) PathEnd(path string, f Handler) {
	h.checkPath(h.rest, path, f)
}

func (h H) checkPath(reqPath, segPath string, f Handler) {
	if reqPath == segPath {
		h.Handle(f)
	}

	if len(segPath) > 0 {
		switch segPath[0] {
		case '*':
			h.Handle(f)
		case ':':
			if len(reqPath) > 0 {
				h.Handle(f)
			}
		}
	}
}

var allowHeader = http.CanonicalHeaderKey("Allow")

func (h H) Allow(verbs ...string) {
	// allowHeader MAY be separated into multiple headers
	// rfc2616-sec4.html#sec4.2
	h.Header()[allowHeader] = append(h.Header()[allowHeader], verbs...)
}

func (h H) hasAllowed() bool {
	_, ok := h.Header()[allowHeader]
	return ok
}

// Verb is a convenience function for
//
//	h.Allow(verb)
//	if h.Method == verb {
//		f(h)
//	}
func (h H) Verb(verb string, f Handler) {
	h.Allow(verb)
	if h.Method == verb {
		h.Handle(f)
	}
}

func (h H) Delete(f Handler) { h.Verb("DELETE", f) }
func (h H) Get(f Handler)    { h.Verb("GET", f) }
func (h H) Patch(f Handler)  { h.Verb("PATCH", f) }
func (h H) Post(f Handler)   { h.Verb("POST", f) }
func (h H) Put(f Handler)    { h.Verb("PUT", f) }

func (h H) Catch(f Handler) {
	err := recover()

	if err == nil && h.Status == 0 && h.hasAllowed() {
		defer h.Catch(f)
		h.Return(http.StatusMethodNotAllowed)
	}

	if res, ok := err.(apiResponse); ok {
		h.Status = res.status
		h.Body = res.body
	} else {
		h.Body = err
	}

	if f != nil {
		f(h)
		h.Return(h.Body)
	}
}

func handlePanic(f Handler) Handler {
	return func(h H) {
		defer func() {
			switch x := recover().(type) {
			case apiResponse, nil:
			default:
				panic(x)
			}
		}()

		f(h)
	}
}

func handleFinish(f Handler) Handler {
	return func(h H) {
		defer func() {
			h.Stream(recover())
		}()

		f(h)
	}
}

func HandleStatus(status int) Middleware {
	return func(f Handler) Handler {
		return func(h H) {
			defer func() {
				err := recover()

				if err == nil {
					err = status
				}

				h.Return(err)
			}()

			f(h)
		}
	}
}

func (h H) Return(body interface{}) {
	if h.Status == 0 && body == nil {
		panic(nil)
	}
	panic(newResponse(h.Status, body))
}

func (h H) Header() http.Header {
	return h.Response.Header()
}

func (h H) Write(p []byte) (int, error) {
	return h.Response.Write(p)
}

var _ http.FileSystem = apiDir{}

type apiDir struct {
	root http.FileSystem
}

func (d apiDir) Open(name string) (http.File, error) {
	f, err := d.root.Open(name)
	if err != nil {
		panic(newResponse(http.StatusNotFound, err))
	}
	return f, nil
}
