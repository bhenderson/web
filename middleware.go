package web

import (
	"io"
	"net/http"

	"github.com/bhenderson/web/flush"
	"github.com/bhenderson/web/head"
	"github.com/bhenderson/web/log"
)

// return a HandlerFunc because that's the common use case.
// if we changed it to return a Handler, we'd have to typecast every time
// we create an annonymous function.
// we accept a Handler however, because HandlerFunc is also a Handler. booyah!
type Middleware func(http.Handler) http.HandlerFunc

type Stack []Middleware

func (s *Stack) Use(ms ...Middleware) {
	*s = append(*s, ms...)
}

// Run takes a http.Handler (http.DefaultServeMux if nil) and builds the
// middleware stack to return a new http.Handler.
func (s *Stack) Run(app http.Handler) (f http.Handler) {
	if app == nil {
		f = http.DefaultServeMux
	} else {
		f = app
	}

	ms := *s
	// reverse
	for i := len(ms) - 1; i >= 0; i-- {
		// The simple case
		// f = ms[i](f)

		f = wrapMiddleware(ms[i], f)
	}
	return
}

func wrapMiddleware(mid Middleware, next http.Handler) http.HandlerFunc {
	// We wrap the next handler in our own handler so we can wrap the
	// response writer, making it so middleware writers don't have to
	// worry about losing Plusser methods.
	return func(w http.ResponseWriter, r *http.Request) {
		m := http.HandlerFunc(func(wr http.ResponseWriter, re *http.Request) {
			wr = WrapResponseWriter(wr, w)
			next.ServeHTTP(wr, re)
		})
		mid(m).ServeHTTP(w, r)
	}
}

var defaultStack = &Stack{}

// Use adds Middleware to the default stack.
func Use(ms ...Middleware) {
	defaultStack.Use(ms...)
}

// Run compiles the default stack of middleware and returns an http.Handler.
func Run(app http.Handler) http.Handler {
	return defaultStack.Run(app)
}

// Plusser is an interface for all the extra methods that http.response has :/
type Plusser interface {
	http.CloseNotifier
	http.Flusher
	http.Hijacker
	io.ReaderFrom
}

// responseWriterPlus implements http.ResponseWriter *and* all the extra
// methods that http.response exposes.
type responseWriterPlus struct {
	http.ResponseWriter
	Plusser
}

// WrapResponseWriter takes an http.ResponseWriter (wr) and wraps it with the
// functionality provided by wn. The return value tries hard to implement any
// extra methods that wr might also implement.
func WrapResponseWriter(wn, wr http.ResponseWriter) http.ResponseWriter {
	if wp, ok := wr.(Plusser); ok {
		return &responseWriterPlus{wn, wp}
	}
	return wn
}

const (
	// Log formats
	CombinedLog = log.Combined
	CommonLog   = log.Common
)

// Flush implements Middleware. See flush.FlushMiddleware for usage.
func Flush(next http.Handler) http.HandlerFunc {
	return flush.FlushMiddleware(next)
}

// Head implements Middleware. See head.HeadMiddleware for usage.
func Head(next http.Handler) http.HandlerFunc {
	return head.HeadMiddleware(next)
}

// Log returns a Middleware. See log.LogMiddleware for usage.
func Log(w io.Writer, t string) Middleware {
	return log.LogMiddleware(w, t)
}
