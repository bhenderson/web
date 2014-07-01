/*
Package web is a just for fun micro framework built on top of net/http

Example:

	package main

	import (
		"fmt"
		"net/http"

		"github.com/bhenderson/web"
	)

	func hello(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintln(w, "hello world")
	}

	func main() {
		web.Use(
			web.Log(os.Stdout, web.Combined),
		)
		http.HandleFunc("/", hello)
		http.ListenAndServe(":8080", web.Run(nil))
	}
*/
package web
