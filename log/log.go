package log

import (
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"text/template"
	"time"
)

var (
	// FuncMap is exposed to add custom functions if desirable.
	FuncMap = template.FuncMap{
		"ftime": (*time.Time).Format,
	}

	dash = "-"
)

const (
	// common log format used by Apache: httpd.apache.org/docs/2.2/logs.html
	Common = `{{.RemoteAddr}} - {{.Username}} [{{ftime .Time "02/Jan/2006:15:04:05 -0700"}}] "{{.Method}} {{.RequestURI}} {{.Proto}}" {{.Status}} {{.ContentSize}}`

	// combined log format used by Apache: httpd.apache.org/docs/2.2/logs.html
	Combined = Common + ` "{{.Referer}}" "{{.UserAgent}}"`
)

// Logger is used internally to track the request/response information, but is
// exported for documentation purposes. Public methods/fields on Logger are
// available in the template.
type Logger struct {
	// The original request
	*http.Request

	// The start time of the request.
	Time time.Time

	// The status of the response
	Status int

	// The content length of the response. See ContentSize
	ContentLength int
}

// Username returns the Username or a "-"
func (l *Logger) Username() string {
	if u := l.URL.User; u != nil {
		return u.Username()
	}
	return dash
}

// RemoteAddr wraps Request.RemoteAddr to remove the port. If not available, this value will be a "-"
func (l *Logger) RemoteAddr() string {
	addr := l.Request.RemoteAddr
	for i := 0; i < len(addr); i++ {
		if addr[i] == ':' {
			addr = addr[:i]
		}
	}
	if len(addr) == 0 {
		return dash
	}
	return addr
}

// ContentSize tries to return the content byte size returned to the client not
// including the headers. If no content was returned (0), this value will be a
// "-". To log "0", use ContentLength
func (l *Logger) ContentSize() string {
	if l.ContentLength == 0 {
		return dash
	}
	return fmt.Sprintf("%d", l.ContentLength)
}

type logWriter struct {
	http.ResponseWriter
	logger *Logger
}

func (lw *logWriter) WriteHeader(code int) {
	lw.logger.Status = code
	lw.ResponseWriter.WriteHeader(code)
}

func (lw *logWriter) Write(p []byte) (int, error) {
	lw.logger.ContentLength += len(p)
	return lw.ResponseWriter.Write(p)
}

// LogMiddleware implements web.Middleware. It takes an io.Writer and template
// string and returns a Middleware which will log the request. See Common and
// Combined for some predefined templates.  See Logger for available fields and
// methods. LogMiddleware panics if template does not compile. If an error is
// returned from the template, that error will be logged to the default logger.
func LogMiddleware(out io.Writer, t string) func(http.Handler) http.HandlerFunc {
	if out == nil {
		out = os.Stdout
	}
	tmp := template.Must(template.New("log").Funcs(FuncMap).Parse(t))
	return func(next http.Handler) http.HandlerFunc {
		return func(w http.ResponseWriter, r *http.Request) {
			lgr := &Logger{
				r,
				time.Now(),
				http.StatusOK,
				0,
			}
			w = &logWriter{w, lgr}

			next.ServeHTTP(w, r)

			err := tmp.Execute(out, lgr)
			out.Write([]byte("\n"))
			if err != nil {
				log.Println(err)
			}
		}
	}
}
