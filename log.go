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

type logWriterPlus struct {
	Plusser
	status int
}

func (lw *logWriterPlus) WriteHeader(code int) {
	lw.status = code
	lw.Plusser.WriteHeader(code)
}

func NewLogger(w http.ResponseWriter, code int) http.ResponseWriter {
	if wp, ok := w.(Plusser); ok {
		return &logWriterPlus{wp, code}
	}
	return &logWriter{w, code}
}

func Log(next http.Handler) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		now := time.Now()
		code := http.StatusOK
		w = NewLogger(w, code)
		next.ServeHTTP(w, r)
		if fw, ok := w.(*logWriter); ok {
			code = fw.status
		}
		log.Println(code, r.URL.Path, time.Since(now))
	}
}
