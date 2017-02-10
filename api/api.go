package api

import (
	"fmt"
	"io"
	"log"
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
	h := H{
		Request:  r,
		Response: NewResponse(w),
		Time:     time.Now(),
		SubPath:  r.URL.Path,
	}
	return h
}

type halt struct{}

func Halt() {
	panic(halt{})
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

	h.Path("", f)
}

type Middleware func(Handler) Handler

type H struct {
	*http.Request
	*Response

	PathSegment, SubPath string

	Middleware []Middleware

	Time time.Time
}

func (h H) Stream(v interface{}) {
	switch x := v.(type) {
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
		h.Stream(h.Status)
	default:
		fmt.Fprintf(h, "%s", x)
	}
}

func (h H) Handle(f Handler) {
	if f == nil {
		return
	}
	h.Use(
		handlePanics,
		handleAllowed,
	)

	// compile in reverse order
	for i := len(h.Middleware); i > 0; i-- {
		f = h.Middleware[i-1](f)
	}
	h.Middleware = h.Middleware[:0]

	h.delAllowed()
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

func (h *H) Next() {
	path := h.SubPath
	for i := 0; i < len(path); i++ {
		if path[i] == '/' {
			h.PathSegment, h.SubPath = path[:i], path[i+1:]
			return
		}
	}
	h.PathSegment, h.SubPath = path, ""
}

func (h *H) Use(ms ...Middleware) {
	h.Middleware = append(h.Middleware, ms...)
}

func (h H) UseBefore(ms ...Middleware) {
	h.Middleware = append(ms, h.Middleware...)
}

// Path runs f if the next PathSegment matches path.
// The first character of path can have special meaning.
//
//	: will match any segment
//	* will match the rest of the request path
func (h H) Path(path string, f Handler) {
	h.Next()
	h.checkPath(path, f)
}

// TODO evaluate removing
// PathEnd runs only when path matches the end of the request path
// It uses the same matching rules as Path
// func (h H) PathEnd(path string, f Handler) {
// h.path, h.rest = h.rest, ""
// h.checkPath(h.path, path, f)
// }

func (h H) checkPath(path string, f Handler) {
	if h.PathSegment == path {
		h.Handle(f)
	}

	if len(path) > 0 {
		switch path[0] {
		case '*':
			h.Handle(f)
		case ':':
			if len(h.PathSegment) > 0 {
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

func (h H) delAllowed() {
	h.Header().Del(allowHeader)
}

// Verb runs f if the Method matches and we're at the end of the path
func (h H) Verb(verb string, f Handler) {
	h.Allow(verb)
	if h.Method == verb && h.SubPath == "" {
		h.Handle(f)
	}
}

// Common Methods rfc=7231#section-4.1
func (h H) Get(f Handler)    { h.Verb("GET", f) }
func (h H) Head(f Handler)   { h.Verb("HEAD", f) }
func (h H) Post(f Handler)   { h.Verb("POST", f) }
func (h H) Put(f Handler)    { h.Verb("PUT", f) }
func (h H) Delete(f Handler) { h.Verb("DELETE", f) }

// Connect // is only used for Proxies
// Options // internally handled, not ment to be customized
func (h H) Trace(f Handler) { h.Verb("Trace", f) }

func (h H) Catch(f Handler) {
	r := recover()

	if f != nil {
		f(h)
	}

	panic(r)
}

func handlePanic(f Handler) Handler {
	return func(h H) {
		defer func() {
			switch x := recover().(type) {
			case nil:
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
			r := recover()
			if _, ok := r.(halt); ok {
				r = h.Response.Body
			}
			h.Stream(r)
		}()

		f(h)
	}
}

func handleAllowed(f Handler) Handler {
	return func(h H) {
		f(h)
		if h.hasAllowed() {
			if h.Method == "OPTIONS" {
				h.Return(http.StatusOK)
			} else {
				h.Return(http.StatusMethodNotAllowed)
			}
		}
		h.delAllowed()
	}
}

func HandleStatus(status int) Middleware {
	return func(f Handler) Handler {
		return func(h H) {
			f(h)
			h.Return(status)
		}
	}
}

func (h H) Return(body interface{}) {
	switch x := body.(type) {
	case apiError:
		// pass
	case error:
		body = apiError{
			x,
			callers(1),
		}
		if h.Status == 0 {
			h.Status = http.StatusInternalServerError
		}
	case int:
		h.Status = x
		// keep body as an int
	case *Response:
		h.Return(*x)
	case Response:
		h.Status = x.Status
		h.Return(x.Body)
	default:
		if h.Status == 0 {
			h.Status = http.StatusOK
		}
	}
	h.Response.Body = body
	Halt()
}

// Do we want to treat regular panics as different from Return?
func handlePanics(f Handler) Handler {
	return func(h H) {
		defer func() {
			r := recover()
			if r == nil {
				return
			}
			if _, ok := r.(halt); !ok {
				h.Status = 500
			}
			panic(r)
		}()

		f(h)
	}
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
		panic(Response{Status: http.StatusNotFound, Body: err})
	}
	return f, nil
}

func debug(v interface{}) {
	log.Printf("%#v\n", v)
}
