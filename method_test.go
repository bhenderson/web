package web

import (
	"github.com/stretchr/testify/assert"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestMethod(t *testing.T) {
	handler := func(w http.ResponseWriter, r *http.Request) {
		http.Error(w, "hello!", http.StatusInternalServerError)
	}

	m := &Method{Get: handler}

	_, body := testMethod(t, "GET", m)
	assert.Equal(t, "hello!\n", body)

	c, _ := testMethod(t, "POST", m)
	assert.Equal(t, 404, c)
}

func TestMethod_Get_Any(t *testing.T) {
	getH := func(w http.ResponseWriter, r *http.Request) {
		http.Error(w, "hello!", http.StatusOK)
	}

	anyH := func(w http.ResponseWriter, r *http.Request) {
		http.Error(w, "any", http.StatusOK)
	}

	notFound := func(w http.ResponseWriter, r *http.Request) {
		http.Error(w, "found, it was not", http.StatusNotFound)
	}

	m := &Method{
		Get:      getH,
		Any:      anyH,
		NotFound: notFound,
	}

	_, body := testMethod(t, "GET", m)
	assert.Equal(t, "hello!\n", body)

	_, body = testMethod(t, "POST", m)
	assert.Equal(t, "any\n", body)
}

func TestMethod_Get_NotFound(t *testing.T) {
	getH := func(w http.ResponseWriter, r *http.Request) {
		http.Error(w, "hello!", http.StatusOK)
	}

	notFound := func(w http.ResponseWriter, r *http.Request) {
		http.Error(w, "found, it was not", http.StatusNotFound)
	}

	m := &Method{
		Get:      getH,
		NotFound: notFound,
	}

	_, body := testMethod(t, "GET", m)
	assert.Equal(t, "hello!\n", body)

	_, body = testMethod(t, "POST", m)
	assert.Equal(t, "found, it was not\n", body)
}

func testMethod(t testing.TB, method string, f http.Handler) (int, string) {
	req, err := http.NewRequest(method, "http://example.com/foo", nil)
	if err != nil {
		t.Fatal(err)
	}

	w := httptest.NewRecorder()
	f.ServeHTTP(w, req)

	return w.Code, w.Body.String()
}
