package web

import (
	"io"
	"net/http"

	"github.com/bhenderson/web/flush"
	"github.com/bhenderson/web/head"
	"github.com/bhenderson/web/log"
	"github.com/bhenderson/web/session"
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
		f = ms[i](f)
	}
	return
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

//go:generate go run genplusser.go

// WrapResponseWriter takes an http.ResponseWriter (wr) and wraps it with the
// functionality provided by wn. The return value tries hard to implement any
// extra methods that wr might also implement.
func WrapResponseWriter(wn, wr http.ResponseWriter) http.ResponseWriter {
	return newPlusser(wn, wr)
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

func Session(secret, name string) Middleware {
	return session.SessionMiddleware(secret, name)
}
