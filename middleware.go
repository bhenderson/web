package web

import "net/http"

type Middleware func(http.HandlerFunc) http.HandlerFunc

type stack []Middleware

func (s *stack) Use(ms ...Middleware) {
	*s = append(*s, ms...)
}

func (s *stack) Run(app http.Handler) (f http.HandlerFunc) {
	f = app.ServeHTTP
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
func Run(app *http.Handler) http.HandlerFunc {
	return defaultStack.Run(app)
}
