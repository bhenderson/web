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

func main() {
	web.Use(
		web.Flush,
		web.Log,
	)

	var app = http.NewServeMux()
	app.Handle(health.Path, health.DefaultHandler)
	app.HandleFunc("/foo", HelloWorld)
	app.HandleFunc("/bar", HelloBar)

	http.Handle("/", web.Run(app))
	http.ListenAndServe(":8080", nil)
}
