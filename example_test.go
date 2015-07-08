package web_test

import (
	"fmt"
	"net/http"
	"os"

	"github.com/bhenderson/web"
)

func hello() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintln(w, "hello world")
	}
}

func Example() {
	// Setup some middleware
	web.Use(
		// Use the combined log format to log all requests to stdout.
		web.Log(os.Stdout, web.CombinedLog),
	)

	// Use net/http to handle routing. Method sets up specific handlers
	// for specific HTTP methods.
	http.Handle("/", &web.Method{
		Get: hello(),
	})

	// Run using net/http DefaultServeMux
	http.ListenAndServe(":8080", web.Run(nil))
}
