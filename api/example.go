// +build ignore

package main

import (
	"encoding/json"
	"net/http"
	"os"
	"text/template"

	"github.com/bhenderson/web/api"
)

func main() {
	api.Use(HandleLogs)

	h := api.Run(func(h api.H) {
		h.Path("api", func(h api.H) {
			h.Use(HandleJSON)

			h.PathEnd("user", func(h api.H) {
				h.Get(getUser)
			})

			h.PathEnd("jira", func(h api.H) {
				h.Get(getJira)
				h.Post(postJira)
			})

			h.Path("applications", func(h api.H) {
				h.Path("", func(h api.H) {
					h.Get(getApps)
					h.Post(newApp)
				})

				h.Path(":appId", func(h api.H) {
					app := lookupApp(h.PathSegment())
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

		h.Path("*", func(h api.H) {
			defer h.Catch(func(h api.H) {
				if h.Status >= 400 {
					h.Request.URL.Path = "/"
					h.ServeFiles(http.Dir("./public"))
				}
			})

			h.ServeFiles(http.Dir("./public"))
		})
	})

	http.ListenAndServe(":8123", h)
}

func HandleLogs(f api.Handler) api.Handler {
	// common log format used by Apache: httpd.apache.org/docs/2.2/logs.html
	Common := `{{.RemoteAddr}} - {{.Username}} [{{.LocalTime}}] "{{.RequestLine}}" {{.Status}} {{.ContentSize}}`

	// combined log format used by Apache: httpd.apache.org/docs/2.2/logs.html
	Combined := Common + ` "{{.Referer}}" "{{.UserAgent}}"`

	myLogs := Combined + " ({{.Since}})\n"

	tmp := template.Must(template.New("logger").Parse(myLogs))

	return func(h api.H) {
		defer h.Catch(func(h api.H) {
			tmp.Execute(os.Stdout, h)
		})

		f(h)
	}
}

func HandleJSON(f api.Handler) api.Handler {
	return func(h api.H) {
		defer h.Catch(func(h api.H) {
			h.Header().Set("Content-Type", "application/json")

			var err error
			h.Body, err = json.Marshal(h.Body)
			if err != nil {
				h.Body = err
			}
			h.Return(h.Body)
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
