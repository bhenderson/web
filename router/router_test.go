package router

import (
	"net/http"
	"testing"
)

var config string

type testHandler string

func (t testHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	config = string(t)
}

func TestLocation(t *testing.T) {
	r := NewRouter()

	r.Location("=", "/", testHandler("A"))
	r.Location("", "/", testHandler("B"))
	r.Location("", "/documents/", testHandler("C"))
	r.Location("^~", "/images/", testHandler("D"))
	r.Location("~*", `\.(gif|jpg|jpeg)$`, testHandler("E"))

	assertMatch := func(path, exp string) {
		config = ""
		req, _ := http.NewRequest("GET", path, nil)
		r.ServeHTTP(nil, req)
		if exp != config {
			t.Errorf(
				"expected %q got %q\n",
				exp, config,
			)
		}
	}

	assertMatch("/", "A")
	assertMatch("/index.html", "B")
	assertMatch("/documents/document.html", "C")
	assertMatch("/images/1.gif", "D")
	assertMatch("/documents/1.jpg", "E")
}
