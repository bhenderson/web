package session

import (
	"net/http"

	"github.com/andreadipersio/securecookie"
)

type sessionWriter struct {
	http.ResponseWriter
	secret, name string

	wroteHeader bool
}

func (w *sessionWriter) Write(buf []byte) (int, error) {
	if !w.wroteHeader {
		w.WriteHeader(http.StatusOK)
	}
	return w.ResponseWriter.Write(buf)
}

func (w *sessionWriter) WriteHeader(i int) {
	w.wroteHeader = true
	SignCookies(w, w.secret, w.name)
	w.ResponseWriter.WriteHeader(i)
}

func SessionMiddleware(secret, name string) func(http.Handler) http.HandlerFunc {
	return func(next http.Handler) http.HandlerFunc {
		return func(w http.ResponseWriter, r *http.Request) {
			DecodeCookies(r, secret, name)
			w = &sessionWriter{
				ResponseWriter: w,
				secret:         secret,
				name:           name,
			}
			next.ServeHTTP(w, r)
		}
	}
}

const (
	SetCookie = "Set-Cookie"
	GetCookie = "Cookie"
)

func SignCookies(w http.ResponseWriter, secret, name string) {
	r := &http.Response{
		Header: w.Header(),
	}
	cookies := r.Cookies()
	w.Header().Del(SetCookie)
	for _, c := range cookies {
		if c.Name == name {
			securecookie.SignCookie(c, secret)
		}
		http.SetCookie(w, c)
	}
}

func DecodeCookies(r *http.Request, secret, name string) {
	cookies := r.Cookies()
	r.Header.Del(GetCookie)
	for _, c := range cookies {
		if c.Name == name {
			if value, err := securecookie.DecodeSignedValue(secret, c.Name, c.Value); err == nil {
				c.Value = value
			}
		}
		r.AddCookie(c)
	}
}
