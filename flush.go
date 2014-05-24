package web

import "net/http"

// TODO should flushWriter implement http.Flusher always?
type flushWriter struct {
	http.ResponseWriter
}

func (fw *flushWriter) Write(p []byte) (int, error) {
	if f, ok := fw.ResponseWriter.(http.Flusher); ok {
		defer f.Flush()
	}
	return fw.ResponseWriter.Write(p)
}

type flushWriterPlus struct {
	Plusser
}

func (fw *flushWriterPlus) Write(p []byte) (int, error) {
	if f, ok := fw.Plusser.(http.Flusher); ok {
		defer f.Flush()
	}
	return fw.Plusser.Write(p)
}

func NewFlusher(w http.ResponseWriter) http.ResponseWriter {
	if wp, ok := w.(Plusser); ok {
		return &flushWriterPlus{wp}
	}
	return &flushWriter{w}
}

// Flush implements Middleware
func Flush(next http.Handler) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w = NewFlusher(w)
		next.ServeHTTP(w, r)
	}
}
