package web

import (
	"log"
	"net/http"
	"time"
)

type logWriter struct {
	http.ResponseWriter
	status int
}

func (lw *logWriter) WriteHeader(code int) {
	lw.status = code
	lw.ResponseWriter.WriteHeader(code)
}

func Log(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		now := time.Now()
		code := http.StatusOK
		w = &logWriter{w, code}
		next.ServeHTTP(w, r)
		if fw, ok := w.(*logWriter); ok {
			code = fw.status
		}
		log.Println(code, r.URL.Path, time.Since(now))
	})
}
