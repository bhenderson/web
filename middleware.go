package web

import (
	"io"
	"net/http"
)

// return a HandlerFunc because that's the common use case.
// if we changed it to return a Handler, we'd have to typecast every time
// we create an annonymous function.
// we accept a Handler however, because HandlerFunc is also a Handler. booyah!
type Middleware func(http.Handler) http.HandlerFunc

type stack []Middleware

func (s *stack) Use(ms ...Middleware) {
	*s = append(*s, ms...)
}

func (s *stack) Run(app http.Handler) (f http.Handler) {
	f = app
	ms := *s
	// reverse
	for i := len(ms) - 1; i >= 0; i-- {
		// The simple case
		// f = ms[i](f)

		// We wrap the next handler in our own handler so we can wrap the
		// response writer, making it so middleware writers don't have to
		// worry about losing Plusser methods.
		next := f
		i := i
		f = http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			mid := http.HandlerFunc(func(wr http.ResponseWriter, re *http.Request) {
				wr = WrapResponseWriter(wr, w)
				next.ServeHTTP(wr, re)
			})
			ms[i](mid).ServeHTTP(w, r)
		})
	}
	return
}

var defaultStack = &stack{}

func Use(ms ...Middleware) {
	defaultStack.Use(ms...)
}

// http.Handle("/", web.Run(app))
func Run(app http.Handler) http.Handler {
	return defaultStack.Run(app)
}

// all the extra methods that http.response has :/
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
