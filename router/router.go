// Package router implements nginx[1] style location routing.
//
// [1] http://nginx.org/en/docs/http/ngx_http_core_module.html#location
package router

import (
	"fmt"
	"net/http"
	"regexp"
)

func NewRouter() *Router {
	return &Router{
		requests: make(map[*http.Request]locationHandler),
	}
}

type Router struct {
	locations []locationHandler
	regexps   []locationHandler

	requests map[*http.Request]locationHandler
	captures Params
}

func (rr *Router) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	h := rr.match(r)
	if h != nil {
		h.ServeHTTP(w, r)
	} else {
		http.NotFound(w, r)
	}
}

func (rr *Router) match(r *http.Request) (h http.Handler) {
	// Check priority, but maintain order. This lets applications configure
	// "popular" routes.

	// check exacts
	// check prefixes
	// check regexps
	// check locations

	var l locationHandler
	found := false
	for _, nl := range rr.locations {
		if nl.match(r) {
			found = true
			// return exact matches right away
			if nl.exact {
				return nl
			}
			// remember longest location
			if len(nl.location) > len(l.location) {
				l = nl
			}
		}
	}

	if found && l.notRegex {
		return l
	}

	for _, nl := range rr.regexps {
		if nl.match(r) {
			return nl
		}
	}

	// return matched location if no regexps
	return l
}

func (r *Router) Location(kind, path string, h http.Handler) {
	switch kind {
	case "=":
		r.LocationExact(path, h)
	case "~":
		re := regexp.MustCompile(path)
		r.LocationRegexp(re, h)
	case "~*": // case insensitive
		re := regexp.MustCompile("(?i)" + path)
		r.LocationRegexp(re, h)
	case "^~":
		r.LocationPrefix(path, h)
	case "":
		r.locations = append(r.locations, locationHandler{
			location: path,
			handler:  h,
		})
	default:
		panic(fmt.Sprintf("%q is not supported", kind))
	}
}

func (r *Router) LocationExact(path string, h http.Handler) {
	r.locations = append(r.locations, locationHandler{
		location: path,
		handler:  h,
		exact:    true,
	})
}

// LocationPrefix mataches the beginning of the url path and skips regexps after longest match.
func (r *Router) LocationPrefix(path string, h http.Handler) {
	r.locations = append(r.locations, locationHandler{
		location: path,
		notRegex: true,
		handler:  h,
	})
}

func (r *Router) LocationRegexp(path *regexp.Regexp, h http.Handler) {
	r.regexps = append(r.regexps, locationHandler{
		regexp:  path,
		handler: h,
	})
}

type locationHandler struct {
	location string
	exact    bool
	notRegex bool
	regexp   *regexp.Regexp
	handler  http.Handler
}

func (h locationHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	h.handler.ServeHTTP(w, r)
}

func (h locationHandler) match(r *http.Request) bool {
	path := r.URL.Path
	if h.exact && h.location == path {
		return true
	}

	if h.regexp == nil {
		if !h.exact && len(path) >= len(h.location) && path[:len(h.location)] == h.location {
			return true
		}
	} else {
		if res := h.regexp.FindStringSubmatch(path); len(res) > 0 {
			// TODO store capture groups
			return true
		}
	}
	return false
}

type Params []Param

type Param struct {
	Key   string
	Value string
}
