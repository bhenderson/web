package main

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/bhenderson/web"
)

func HelloWorld(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "hello ")
	time.Sleep(0 * time.Second)
	fmt.Fprintf(w, "world")
}

func main() {
	web.Use(
		web.Log(os.Stdout, web.CombinedLog),
		web.Flush,
		web.Head,
	)

	app := http.NewServeMux()
	app.Handle("/foo", &web.Method{
		Get: http.HandlerFunc(HelloWorld),
	})
	app.Handle("/users/", &web.Resource{
		Index:  http.HandlerFunc(Index),
		Show:   Lookup(Show),
		Update: Lookup(Update),
	})

	// 404 any other endpoints
	http.Handle("/", web.Run(app))
	log.Fatal(http.ListenAndServe(":8080", nil))
}

func Index(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "list of users:\n")
}

type User struct {
	ID   string
	Name string
}

type UserHandler func(*User) http.HandlerFunc

func Lookup(h UserHandler) web.ResourceHandleFunc {
	return func(id string) http.Handler {
		// simulate DB lookup or whatever
		if id == "abc" {
			u := &User{id, "John Smith"}
			return h(u)
		}
		// not found
		return http.NotFoundHandler()
	}
}

func Show(u *User) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintf(w, "Welcome %s!\n", u.Name)
	}
}

func Update(u *User) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintf(w, "User %s was updated!\n", u.Name)
	}
}
