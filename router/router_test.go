package router

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"reflect"
	"testing"
)

var config string

type testHandler string

func (t testHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	fmt.Println(t)
}

func ExampleRouter_Location() {
	r := NewRouter()

	// testHandler outputs the string when called
	r.Location("=", "/", testHandler("A"))
	r.Location("", "/", testHandler("B"))
	r.Location("", "/documents/", testHandler("C"))
	r.Location("^~", "/images/", testHandler("D"))
	r.Location("~*", `\.(gif|jpg|jpeg)$`, testHandler("E"))

	// serve builds a *http.Request and calls r.ServeHTTP
	serve(r, "/")
	serve(r, "/index.html")
	serve(r, "/documents/document.html")
	serve(r, "/images/1.gif")
	serve(r, "/documents/1.jpg")
	serve(r, "/documents/1.JPG")
	// Output:
	// A
	// B
	// C
	// D
	// E
	// E
}

func serve(r *Router, path string) *httptest.ResponseRecorder {
	req, _ := http.NewRequest("GET", path, nil)
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)
	return w
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

	// test Params after cleared
	req, _ := http.NewRequest("GET", "/docs/2.jpg", nil)
	act := r.Params(req)
	exp := Params{{"format", "jpg"}}
	assertEqual(t, exp, act)
}

func TestNotFound(t *testing.T) {
	r := NewRouter()
	r.NotFound = func(w http.ResponseWriter, r *http.Request) {
		config = "route not found"
	}

	tc := testCapture{r}

	r.Location("~*", `\.(?P<format>gif|jpg|jpeg)$`, tc)

	assertMatch(t, r, "/foobar", "route not found")

	r.NotFound = nil
	w := serve(r, "/barfoo")

	assertEqual(t, "404 page not found\n", w.Body.String())
}

func TestPanic(t *testing.T) {
	r := NewRouter()

	assertPanic(t, `"foo" is not supported`, func() {
		r.Location("foo", "/path", nil)
	})

	assertPanic(t, "regexp: Compile(`[bad regex`): error parsing regexp: missing closing ]: `[bad regex`", func() {
		r.Location("~", "[bad regex", nil)
	})
}

func assertPanic(t *testing.T, exp interface{}, f func()) {
	defer func() {
		assertEqual(t, exp, recover())
	}()
	f()
}

func assertMatch(t *testing.T, r *Router, path, exp string) {
	config = ""
	serve(r, path)
	assertEqual(t, exp, config)
}

func assertEqual(t *testing.T, exp, act interface{}) {
	if !reflect.DeepEqual(exp, act) {
		t.Errorf(
			"expected %#v got %#v\n",
			exp, act,
		)
	}
}
