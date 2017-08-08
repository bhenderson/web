package head

import (
	"net/http"
	"strconv"
)

type headWriter struct {
	http.ResponseWriter

	size int
}

func (hw *headWriter) Write(p []byte) (int, error) {
	if hw.Header().Get("Content-Type") == "" {
		hw.Header().Set("Content-Type", http.DetectContentType(p))
	}
	hw.size += len(p)
	return len(p), nil
}

func (hw *headWriter) WriteHeader(i int) {
	hw.Header().Set("Content-Length", strconv.Itoa(hw.size))
	hw.ResponseWriter.WriteHeader(i)
}

// HeadMiddleware implements web.Middleware. If the request Method is "HEAD",
// it changes the Method to "GET", for the next http.Handler, but doesn't write
// any response. It does tries to preserve what would be the Content-Length header.
func HeadMiddleware(next http.Handler) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method == "HEAD" {
			r.Method = "GET"
			w = &headWriter{w, 0}
		}

		next.ServeHTTP(w, r)

		w.WriteHeader(http.StatusOK)
	}
}
