package panic

import (
	"net/http"
)

const perror = "<html><head><title>Error</title></head><body><h1>Internal Server Error</h1></body></html>"

func PanicMiddleware(next http.Handler) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		defer func() {
			if err := recover(); err != nil {
				http.Error(w, perror, http.StatusInternalServerError)
			}
		}()
		next.ServeHTTP(w, r)
	}
}
