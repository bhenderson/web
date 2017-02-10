// +build ignore

package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net"
	"net/http"
	"os"
	"strconv"
	"text/template"
	"time"

	"github.com/bhenderson/web/api"
)

func main() {
	api.Use(HandleLogs)

	h := api.Run(func(h api.H) {
		h.Path("api", func(h api.H) {
			h.Use(HandleJSON)

			h.Path("user", func(h api.H) {
				h.Get(getUser)
			})

			h.Path("jira", func(h api.H) {
				h.Get(getJira)
				h.Post(postJira)
			})

			h.Path("applications", func(h api.H) {
				h.Get(getApps)
				h.Post(newApp)

				h.Path(":appId", func(h api.H) {
					app := lookupApp(h.PathSegment)
					if app == nil {
						h.Return(404)
					}

					h.Get(app.showApp)
					h.Put(app.updateApp)
					h.Path("tags", func(h api.H) {
						h.Get(app.appTags)
					})
				})
			})
		})

		defer h.Catch(func(h api.H) {
			log.Println("inside catch")
			if h.Status == 404 {
				h.SubPath = "/"
				h.ServeFiles(http.Dir("./public"))
			}
		})

		h.ServeFiles(http.Dir("./public"))
	})

	http.ListenAndServe(":8123", h)
}

func HandleLogs(f api.Handler) api.Handler {
	// common log format used by Apache: httpd.apache.org/docs/2.2/logs.html
	Common := `{{.RemoteAddr}} - {{.Username}} [{{.LocalTime}}] "{{.RequestLine}}" {{.Status}} {{.ContentSize}}`

	// combined log format used by Apache: httpd.apache.org/docs/2.2/logs.html
	Combined := Common + ` "{{.Referer}}" "{{.UserAgent}}"`

	myLogs := Combined + " ({{.Since}})\n"

	tmp := template.Must(template.New("logger").
		Parse(myLogs))

	return func(h api.H) {
		defer h.Catch(func(h api.H) {
			tmp.Execute(os.Stdout, Logger{h})
		})

		f(h)
	}
}

func HandleJSON(f api.Handler) api.Handler {
	return func(h api.H) {
		defer h.Catch(func(h api.H) {
			h.Header().Set("Content-Type", "application/json")

			var err error
			body := h.Response.Body
			body, err = json.Marshal(body)
			if err != nil {
				body = err
			}
			h.Return(body)
		})

		f(h)
	}
}

func getUser(h api.H) {
	u := map[string]string{
		"name":  "auser",
		"email": "auser@example.com",
	}

	h.Return(u)
}

func getJira(h api.H) {
	h.Return(map[string]string{
		"jira": "get stuff",
	})
}

func postJira(h api.H) {
	h.Return(map[string]string{
		"jira": "post stuff",
	})
}

func getApps(h api.H) {
	h.Return(map[string][]App{
		"apps": {{"app1"}, {"app2"}},
	})
}

func newApp(h api.H) {
	h.Return(App{
		"app1",
	})
}

func lookupApp(id string) *App {
	if id == "joe" {
		return nil
	}
	return &App{id}
}

type App struct {
	Name string
}

func (a *App) showApp(h api.H) {
	h.Return(a)
}

func (a *App) updateApp(h api.H) {
	h.ParseForm()
	a.Name = h.Form.Get("name")
	h.Return(a)
}

func (a *App) appTags(h api.H) {
	h.Return(map[string][]string{
		"tags": {"tag1", "tag2"},
	})
}

// Logging methods

type Logger struct {
	api.H
}

const dash = "-"

func (l Logger) LocalTime() string {
	return l.Time.Format("02/Jan/2006:15:04:05 -0700")
}

// Username returns the Username or a "-"
func (l Logger) Username() string {
	if l.URL != nil {
		if u := l.URL.User; u != nil {
			return u.Username()
		}
	}
	return dash
}

func (l Logger) RemoteAddr() string {
	host, _, err := net.SplitHostPort(l.Request.RemoteAddr)
	if err != nil {
		return dash
	}
	return host
}

func (l Logger) RequestLine() string {
	return fmt.Sprintf("%s %s %s", l.Method, l.RequestURI, l.Proto)
}

func (l Logger) ContentSize() string {
	ln := l.GetContentLength()
	if ln < 0 {
		return dash
	}
	return strconv.FormatInt(ln, 10)
}

func (l Logger) Since() time.Duration {
	return time.Since(l.Time)
}

// end Logging methods
