// Package router implements nginx[1] style location routing.
//
// [1] http://nginx.org/en/docs/http/ngx_http_core_module.html#location
package router

import (
	"fmt"
	"net/http"
	"regexp"
	"sort"
)

func NewRouter() *Router {
	return &Router{
		requests: make(map[*http.Request]*locationHandler),
		exact:    make(map[string]locationHandler),
	}
}

type Router struct {
	// NotFound handles the case when no location matches. Defaults to
	// http.NotFound
	NotFound http.HandlerFunc

	exact     map[string]locationHandler
	locations []locationHandler
	regexps   []locationHandler

	requests map[*http.Request]*locationHandler
}

func (r *Router) ServeHTTP(w http.ResponseWriter, req *http.Request) {
	defer func() {
		// cleanup r.requests after processing.
		delete(r.requests, req)
	}()

	h := r.match(req)
	if h != nil {
		h.ServeHTTP(w, req)
	} else if r.NotFound != nil {
		r.NotFound(w, req)
	} else {
		http.NotFound(w, req)
	}
}

func (r *Router) match(req *http.Request) (lp *locationHandler) {
	// Check priority, but maintain order. This lets applications configure
	// "popular" routes.

	// check exacts
	// check prefixes
	// check regexps
	// check locations

	// Seems messy. should I pass path to all match?
	if l, ok := r.exact[req.URL.Path]; ok {
		return &l
	}

	// locations are now sorted by longest path, so we use the first matched.
	for _, nl := range r.locations {
		if lp = nl.match(req); lp != nil {
			if lp.noRegexs {
				return lp
			}
			break
		}
	}

	for _, nl := range r.regexps {
		if l := nl.match(req); l != nil {
			r.requests[req] = l
			return l
		}
	}

	// return matched location if no regexps
	return lp
}

// Location takes similar parameters to nginx Location directive (see package docs.)
/*
	=  LocationExact
	~  LocationRegexp
	~* LocationRegexp (?i)
	^~ LocationPrefix (stop regexp matching)
	"" Location (prefix matching)
*/
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
		r.addLocation(locationHandler{
			location: path,
			handler:  h,
		})
	default:
		panic(fmt.Sprintf("%q is not supported", kind))
	}
}

func (r *Router) LocationFunc(kind, path string, h http.HandlerFunc) {
	r.Location(kind, path, h)
}

// LocationExact matches the URL Path exactly. Route processing immediately
// stops.
//
// Same as Location("=", path, h)
func (r *Router) LocationExact(path string, h http.Handler) {
	r.exact[path] = locationHandler{
		exact:   true,
		handler: h,
	}
}

// LocationPrefix matches the beginning of the url path and skips regexps after longest match.
//
// Same as Location("^~", path, h)
func (r *Router) LocationPrefix(path string, h http.Handler) {
	r.addLocation(locationHandler{
		location: path,
		noRegexs: true,
		handler:  h,
	})
}

// LocationRegexp matches the URL Path against the regexp.
//
// Same as Location("~", path, h)
// To get case-insensitive matching, compile your regexp with (?i).
func (r *Router) LocationRegexp(path *regexp.Regexp, h http.Handler) {
	r.regexps = append(r.regexps, locationHandler{
		regexp:   path,
		capNames: path.SubexpNames()[1:],
		handler:  h,
	})
}

// Params returns the capture groups of a regexp location. If the cached
// version is not found, request processing is re-run to find the params. This
// is only the case if called outside r.ServeHTTP (like in a go routine).
func (r *Router) Params(req *http.Request) Params {
	if l, ok := r.requests[req]; ok {
		return l.params()
	}

	if l := r.match(req); l != nil {
		return l.params()
	}
	return nil
}

func (r *Router) addLocation(h locationHandler) {
	r.locations = append(r.locations, h)
	sort.Stable(locationHandlers(r.locations))
}

// sort by longest location
type locationHandlers []locationHandler

func (s locationHandlers) Len() int           { return len(s) }
func (s locationHandlers) Less(i, j int) bool { return -len(s[i].location) < -len(s[j].location) }
func (s locationHandlers) Swap(i, j int)      { s[i], s[j] = s[j], s[i] }

type locationHandler struct {
	location string
	exact    bool
	noRegexs bool

	regexp               *regexp.Regexp
	capNames, capResults []string

	handler http.Handler
}

func (h locationHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	h.handler.ServeHTTP(w, r)
}

// match needs to return a pointer in the case of modifying capResults
// takes care to return a copy
func (h locationHandler) match(r *http.Request) *locationHandler {
	if h.exact {
		return &h
	}

	path := r.URL.Path

	if h.regexp == nil {
		if len(path) >= len(h.location) && path[:len(h.location)] == h.location {
			return &h
		}
	} else {
		if res := h.regexp.FindStringSubmatch(path); len(res) > 0 {
			h.capResults = res[1:]
			return &h
		}
	}
	return nil
}

func (h locationHandler) params() (p Params) {
	// capNames and capResults will always be the same size
	for i, k := range h.capNames {
		p = append(p, Param{k, h.capResults[i]})
	}
	return
}

type Params []Param

type Param struct {
	Key   string
	Value string
}

func (ps Params) ByName(name string) (value string) {
	for _, p := range ps {
		if p.Key == name {
			return p.Value
		}
	}
	return ""
}
