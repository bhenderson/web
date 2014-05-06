package main

import (
	"fmt"
	"github.com/bhenderson/health"
	"github.com/bhenderson/web"
	"net/http"
	"time"
)

func HelloWorld(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(404)
	fmt.Fprintf(w, "hello ")
	time.Sleep(2 * time.Second)
	fmt.Fprintf(w, "world")
}

func HelloBar(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusCreated)
	fmt.Fprintf(w, "hello bar")
}

func CreateBar(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "created bar")
}

func HelloUsers(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "hello all")
}

func HelloUser(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "hello %s", r.Form.Get("_resourceId"))
}

func HelloPost(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w,
		"%s: post %s",
		r.Form.Get("_user_id"),
		r.Form.Get("_post_id"),
	)
}

func CreateUser(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "create user")
}

// Handler
func GoodbyeUser() http.HandlerFunc {
	mux := http.NewServeMux()
	mux.HandleFunc("/edit", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintf(w, "goodbye!")
	})
	return mux
}

func main() {
	web.Use(
		web.Flush,
		web.Log,
	)

	app := http.NewServeMux()
	app.Handle(health.Path, health.DefaultHandler)
	app.HandleFunc("/foo", HelloWorld)
	// web.Handle(app, "GET", "/foo", HelloBar)
	// app.HandleFunc("/bar", HelloBar)
	app.Handle("/bar", &web.Method{
		Get:  HelloBar,
		Post: CreateBar,
	})
	// GET /users/         # index
	// GET /users/:id      # Show
	// GET /users/:id/edit # ???
	app.Handle("/users/", &web.Resource{
		Index:  HelloUsers,
		Show:   HelloUser,
		Create: CreateUser,
		Handler: &web.Resource{
			Prefix: "posts",
			Show:   HelloPost,
		},
	})

	http.Handle("/", web.Run(app))
	http.ListenAndServe(":8080", nil)
}
