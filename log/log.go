package log

import (
	"fmt"
	"io"
	"log"
	"net"
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

// internal devNull io.Writer
type devNull int

func (devNull) Write(p []byte) (int, error) {
	return len(p), nil
}

const (
	// common log format used by Apache: httpd.apache.org/docs/2.2/logs.html
	Common = `{{.RemoteAddr}} - {{.Username}} [{{.LocalTime}}] "{{.RequestLine}}" {{.Status}} {{.ContentSize}}`

	// combined log format used by Apache: httpd.apache.org/docs/2.2/logs.html
	Combined = Common + ` "{{.Referer}}" "{{.UserAgent}}"`
)

func defaultTmp() *template.Template {
	return template.New("log").Funcs(FuncMap)
}

var (
	discard     = devNull(0)
	blankLogger = &Logger{Request: http.Request{}}
	// make sure the templates work
	_ = template.Must(defaultTmp().Parse(Common)).Execute(discard, blankLogger)
	_ = template.Must(defaultTmp().Parse(Combined)).Execute(discard, blankLogger)
)

// Logger is used internally to track the request/response information, but is
// exported for documentation purposes. Public methods/fields on Logger are
// available in the template.
type Logger struct {
	// A copy of the original request
	http.Request

	// The start time of the request.
	Time time.Time

	// The status of the response
	Status int

	// The content length of the response. See ContentSize
	ContentLength int
}

func (l *Logger) LocalTime() string {
	return l.Time.Format("02/Jan/2006:15:04:05 -0700")
}

// Username returns the Username or a "-"
func (l *Logger) Username() string {
	if l.URL != nil {
		if u := l.URL.User; u != nil {
			return u.Username()
		}
	}
	return dash
}

// RemoteAddr wraps Request.RemoteAddr to remove the port. If not available, this value will be a "-"
func (l *Logger) RemoteAddr() string {
	host, _, err := net.SplitHostPort(l.Request.RemoteAddr)
	if err != nil {
		return dash
	}
	return host
}

func (l *Logger) RequestLine() string {
	return fmt.Sprintf("%s %s %s", l.Method, l.RequestURI, l.Proto)
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

// Since returns the elapsed time of the request in nanoseconds.
// Suggested usage:
//	"{{.Since.Seconds}}s"
func (l *Logger) Since() time.Duration {
	return time.Since(l.Time)
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
	if lw.logger.Status <= 0 {
		lw.logger.Status = http.StatusOK
	}
	lw.logger.ContentLength += len(p)
	return lw.ResponseWriter.Write(p)
}

// LogMiddleware takes an io.Writer and template string and returns a
// web.Middleware which will log the request. See Common and Combined for some
// predefined templates. See Logger for available fields and methods.
// LogMiddleware panics if template does not compile. If an error is returned
// from the template, that error will be logged to the default logger.
func LogMiddleware(out io.Writer, t string) func(http.Handler) http.HandlerFunc {
	if out == nil {
		out = os.Stdout
	}
	// add newline to template string if not there.
	// We can't just write a newline after the template executes because it
	// needs to be automic.
	if len(t) > 0 && t[len(t)-1] != '\n' {
		t = t + "\n"
	}
	tmp := template.Must(defaultTmp().Parse(t))
	// check template methods
	if err := tmp.Execute(discard, blankLogger); err != nil {
		panic(err)
	}

	return func(next http.Handler) http.HandlerFunc {
		return func(w http.ResponseWriter, r *http.Request) {
			// read only (does it matter?)
			// I guess this is just a shallow copy...
			copyR := *r
			lgr := &Logger{
				copyR,
				time.Now(),
				http.StatusOK,
				0,
			}
			w = &logWriter{w, lgr}

			next.ServeHTTP(w, r)

			err := tmp.Execute(out, lgr)
			if err != nil {
				log.Println(err)
			}
		}
	}
}
