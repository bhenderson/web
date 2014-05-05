package web

import "net/http"

type flushWriter struct {
	http.ResponseWriter
}

func (fw *flushWriter) Write(p []byte) (int, error) {
	if f, ok := fw.ResponseWriter.(http.Flusher); ok {
		defer f.Flush()
	}
	return fw.ResponseWriter.Write(p)
}

// Flush implements Middleware
func Flush(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w = &flushWriter{w}
		next.ServeHTTP(w, r)
	})
}
