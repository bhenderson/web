package web

import (
	"github.com/stretchr/testify/assert"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestMethod_Get_Any(t *testing.T) {
	var getH http.HandlerFunc = func(w http.ResponseWriter, r *http.Request) {
		http.Error(w, "hello!", http.StatusOK)
	}

	var anyH http.HandlerFunc = func(w http.ResponseWriter, r *http.Request) {
		http.Error(w, "any", http.StatusOK)
	}

	var notFound http.HandlerFunc = func(w http.ResponseWriter, r *http.Request) {
		http.Error(w, "found, it was not", http.StatusNotFound)
	}

	m := &Method{
		Get:              getH,
		Any:              anyH,
		MethodNotAllowed: notFound,
	}

	_, body := testMethod(t, "GET", m)
	assert.Equal(t, "hello!\n", body)

	_, body = testMethod(t, "POST", m)
	assert.Equal(t, "any\n", body)
}

func TestMethod_Get_MethodNotFound(t *testing.T) {
	var getH http.HandlerFunc = func(w http.ResponseWriter, r *http.Request) {
		http.Error(w, "hello!", http.StatusOK)
	}

	var putH http.HandlerFunc = func(w http.ResponseWriter, r *http.Request) {
		http.Error(w, "created!", http.StatusCreated)
	}

	m := &Method{
		Get: getH,
		Put: putH,
	}

	w := testRequest(t, "GET", m)
	assert.Equal(t, 200, w.Code)
	assert.Equal(t, "hello!\n", w.Body.String())

	w = testRequest(t, "POST", m)
	assert.Equal(t, 405, w.Code)
	assert.Equal(t, "405 method not allowed\n", w.Body.String())
	assert.Equal(t, "GET, PUT", w.HeaderMap.Get("Allow"))
}

func TestMethod_Options(t *testing.T) {
	var getH http.HandlerFunc = func(w http.ResponseWriter, r *http.Request) {
		http.Error(w, "hello!", http.StatusOK)
	}

	var putH http.HandlerFunc = func(w http.ResponseWriter, r *http.Request) {
		http.Error(w, "created!", http.StatusCreated)
	}

	m := &Method{
		Get: getH,
		Put: putH,
	}

	w := testRequest(t, "OPTIONS", m)
	assert.Equal(t, 200, w.Code)
	assert.Equal(t, "", w.Body.String())
	assert.Equal(t, "GET, PUT", w.HeaderMap.Get("Allow"))
	// Server implementation apparently takes care of this.
	// assert.Equal(t, "0", w.HeaderMap.Get("Content-Length"))
}

func TestAllowedMethods(t *testing.T) {
	handler := &Method{}
	tests := []struct {
		msg string
		m   *Method
		exp []string
	}{
		{
			"all",
			&Method{
				Delete:  handler,
				Get:     handler,
				Head:    handler,
				Options: handler,
				Patch:   handler,
				Post:    handler,
				Put:     handler,
			},
			[]string{
				"DELETE",
				"GET",
				"HEAD",
				"PATCH",
				"POST",
				"PUT",
			},
		},
		{
			"empty",
			&Method{},
			[]string{},
		},
		{
			"get as head",
			&Method{
				GetAsHead: true,
				Get:       handler,
			},
			[]string{
				"GET",
				"HEAD",
			},
		},
	}
	for _, test := range tests {
		assert.Equal(t, test.exp, allowedMethods(test.m), test.msg)
	}
}

func testRequest(t testing.TB, method string, f http.Handler) *httptest.ResponseRecorder {
	req, err := http.NewRequest(method, "http://example.com/foo", nil)
	if err != nil {
		t.Fatal(err)
	}

	w := httptest.NewRecorder()
	f.ServeHTTP(w, req)

	return w
}

func testMethod(t testing.TB, method string, f http.Handler) (int, string) {
	w := testRequest(t, method, f)

	return w.Code, w.Body.String()
}
