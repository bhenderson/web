package head

import (
	"fmt"
	"net/http"
)

type headWriter struct {
	size int
	http.ResponseWriter
}

func (hw *headWriter) Write(p []byte) (int, error) {
	hw.size += len(p)
	return len(p), nil
}

// HeadMiddleware implements web.Middleware. If the request Method is "HEAD",
// it changes the Method to "GET", for the next http.Handler, but doesn't write
// any response. It does tries to preserve what would be the Content-Length header.
func HeadMiddleware(next http.Handler) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		isHead := r.Method == "HEAD"
		if isHead {
			r.Method = "GET"
			w = &headWriter{0, w}
		}

		next.ServeHTTP(w, r)

		if isHead {
			size := fmt.Sprintf("%d", w.(*headWriter).size)
			w.Header().Set("Content-Length", size)
		}
	}
}
