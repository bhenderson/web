package api

import (
	"fmt"
	"io"
	"net"
	"net/http"
	"strings"
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
		Request:        r,
		ResponseWriter: w,
		Time:           time.Now(),
		path:           path,
		rest:           rest,
	}
	return h
}

type Handler func(H)

func (f Handler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	h := newH(w, r)

	h.Use(handleSafePanic)
	h.Use(DefaultMiddleware...)
	h.Use(
		handleFinish,
		HandleStatus(http.StatusNotFound),
	)

	h.Handle(f)
}

type Middleware func(Handler) Handler

type H struct {
	Status int
	Body   interface{}

	*http.Request

	http.ResponseWriter
	wroteHeader bool

	path, rest string
	Middleware []Middleware

	Time time.Time

	// cheat by using a map (a 0 capacity, fixed pointer slice)
	allowedVerbs map[string]struct{}
}

func (h H) Stream(v interface{}) {
	switch x := v.(type) {
	case apiResponse:
		h.Status = x.status
		h.Stream(x.body)
	case http.Handler:
		x.ServeHTTP(&h, h.Request)
	case int:
		h.Status = x
		h.WriteString(http.StatusText(x))
	case []byte:
		h.Write(x)
	case string:
		h.WriteString(x)
	case io.Reader:
		io.Copy(&h, x)
	case nil:
		// noop
	default:
		fmt.Fprintf(&h, "%s", x)
	}
}

func (h H) Handle(f Handler) {
	// compile in reverse order
	for i := len(h.Middleware); i > 0; i-- {
		f = h.Middleware[i-1](f)
	}
	h.Middleware = h.Middleware[:0]

	h.allowedVerbs = emptyAllowedVerbs()

	f(h)
}

func (h H) HandleHTTP(f http.Handler) {
	// TODO
	h.Handle(func(h H) {
		f.ServeHTTP(&h, h.Request)
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

func (h H) Allow(verbs ...string) {
	if len(verbs) == 0 {
		for v := range h.allowedVerbs {
			verbs = append(verbs, v)
		}
	}
	for _, v := range verbs {
		if h.Method == v {
			return
		}
	}

	h.Header().Set("Allow", strings.Join(verbs, ", "))
	h.Return(http.StatusMethodNotAllowed)
}

// Verb is a convenience function for
//
//	if h.Method == verb {
//		f(h)
//	}
func (h H) Verb(verb string, f Handler) {
	h.allowedVerbs[verb] = struct{}{}
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

	if err == nil && h.Status == 0 && len(h.allowedVerbs) != 0 {
		defer h.Catch(f)
		h.Allow()
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

func handleSafePanic(f Handler) Handler {
	return func(h H) {
		defer func() {
			err := recover()
			if err != nil {
				panic(err)
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
	return h.ResponseWriter.Header()
}

func (h *H) WriteHeader(status int) {
	if status == 0 {
		status = http.StatusOK
	}
	h.ResponseWriter.WriteHeader(status)

	if h.wroteHeader {
		return
	}

	h.Status = status
	h.wroteHeader = true
}

func (h *H) Write(p []byte) (int, error) {
	if !h.wroteHeader {
		h.WriteHeader(h.Status)
	}
	return h.ResponseWriter.Write(p)
}

func (h *H) WriteString(p string) (int, error) {
	return h.Write([]byte(p))
}

// Logging methods
const dash = "-"

func (h H) LocalTime() string {
	return h.Time.Format("02/Jan/2006:15:04:05 -0700")
}

// Username returns the Username or a "-"
func (h H) Username() string {
	if h.URL != nil {
		if u := h.URL.User; u != nil {
			return h.Username()
		}
	}
	return dash
}

func (h H) RemoteAddr() string {
	host, _, err := net.SplitHostPort(h.Request.RemoteAddr)
	if err != nil {
		return dash
	}
	return host
}

func (h H) RequestLine() string {
	return fmt.Sprintf("%s %s %s", h.Method, h.RequestURI, h.Proto)
}

// ContentSize tries to return the content byte size returned to the client not
// including the headers. If no content was returned (0), this value will be a
// "-". To log "0", use ContentLength
func (h H) ContentSize() string {
	return dash // I don't have a good way of calculating this yet.
}

// Since returns the elapsed time of the request in nanoseconds.
// Suggested usage:
//	"{{.Since.Seconds}}s"
func (h H) Since() time.Duration {
	return time.Since(h.Time)
}

// end logging

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

func emptyAllowedVerbs() map[string]struct{} {
	return make(map[string]struct{})
}
