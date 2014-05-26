package web

import (
	// "github.com/stretchr/testify/assert"
	"bufio"
	"io"
	"net"
	"net/http"
	"net/http/httptest"
	"testing"
)

type testResponse struct {
	// shortcut
	*httptest.ResponseRecorder
	Hijacked      bool
	CloseNotified bool
	ReadedFrom    bool
}

func (tr *testResponse) Hijack() (net.Conn, *bufio.ReadWriter, error) {
	tr.Hijacked = true
	return nil, nil, nil
}

func (tr *testResponse) CloseNotify() <-chan bool {
	c := make(chan bool)
	tr.CloseNotified = true
	return c
}

func (tr *testResponse) ReadFrom(r io.Reader) (n int64, err error) {
	tr.ReadedFrom = true
	return 0, nil
}

func NewTestResponse() *testResponse {
	rr := httptest.NewRecorder()
	return &testResponse{
		ResponseRecorder: rr,
	}
}

func Test_testResponse(t *testing.T) {
	tr := interface{}(NewTestResponse())

	if _, ok := tr.(http.ResponseWriter); !ok {
		t.Fatal("expected testResponse to implement http.ResponseWriter")
	}

	if _, ok := tr.(http.Flusher); !ok {
		t.Fatal("expected testResponse to implement http.Flusher")
	}

	if _, ok := tr.(http.CloseNotifier); !ok {
		t.Fatal("expected testResponse to implement http.CloseNotifier")
	}

	if _, ok := tr.(http.Hijacker); !ok {
		t.Fatal("expected testResponse to implement http.Hijacker")
	}

	if _, ok := tr.(io.ReaderFrom); !ok {
		t.Fatal("expected testResponse to implement io.ReaderFrom")
	}
}
