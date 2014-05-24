package main

import (
	"fmt"
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

func main() {
	web.Use(
		web.Log,
		web.Flush,
	)

	app := http.NewServeMux()
	app.HandleFunc("/foo", HelloWorld)

	http.Handle("/", web.Run(app))
	http.ListenAndServe(":8080", nil)
}
