package flush

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

// FlushMiddleware implements web.Middleware. It tries to flush the underlying
// writer on every write.
func FlushMiddleware(next http.Handler) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w = &flushWriter{w}
		next.ServeHTTP(w, r)
	}
}
