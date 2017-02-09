package api_test

import (
	"fmt"
	"io"
	"net/http"
	"net/http/httptest"

	"github.com/bhenderson/web/api"
)

func ExampleH_Path() {
	h := api.Run(func(h api.H) {
		h.Path("foo", func(h api.H) {
			h.Return("hi from foo")
		})
	})

	serve(h, "GET", "/foo", nil)
	serve(h, "GET", "/bar", nil)
	// Output: 200 hi from foo
	// 404 Not Found
}

func ExampleH_Verb() {
	h := api.Run(func(h api.H) {
		defer h.Catch(func(h api.H) {
			fmt.Println("Status set to", h.Status)
		})
		h.Post(func(h api.H) {
			h.Return("post at foo")
		})
	})

	serve(h, "POST", "/foo", nil)
	serve(h, "GET", "/foo", nil)
	// Output: Status set to 200
	// 200 post at foo
	// Status set to 405
	// 405 Method Not Allowed
}

func serve(h http.Handler, verb, path string, body io.Reader) {
	w := httptest.NewRecorder()
	r := httptest.NewRequest(verb, path, body)
	h.ServeHTTP(w, r)

	fmt.Println(w.Code, w.Body)
}
