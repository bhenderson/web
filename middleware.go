package web

import "net/http"

// return a HandlerFunc because that's the common use case.
// if we changed it to return a Handler, we'd have to typecast every time
// we create an annonymous function.
// we accept a Handler however, because that's more "generic"
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
		f = ms[i](f)
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
