package router

import (
	"fmt"
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

	assertMatch(t, r, "/", "A")
	assertMatch(t, r, "/index.html", "B")
	assertMatch(t, r, "/documents/document.html", "C")
	assertMatch(t, r, "/images/1.gif", "D")
	assertMatch(t, r, "/documents/1.jpg", "E")
}

type testCapture struct {
	r *Router
}

func (t testCapture) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	p := t.r.Params(r)
	config = fmt.Sprintf("%v", p)
}

func TestCaptures(t *testing.T) {
	r := NewRouter()

	tc := testCapture{r}

	r.Location("~*", `\.(?P<format>gif|jpg|jpeg)$`, tc)

	assertMatch(t, r, "/images/1.gif", `[{format gif}]`)
}

func assertMatch(t *testing.T, r *Router, path, exp string) {
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
