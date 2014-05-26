package web

import (
	"github.com/stretchr/testify/assert"
	"net/http"
	"net/http/httptest"
	"testing"
)

func Test_WrapLogWriter(t *testing.T) {
	var w http.ResponseWriter = httptest.NewRecorder()
	tr := NewTestResponse()

	fw := &flushWriter{w}
	lw := &logWriter{w, 200}
	w = WrapResponseWriter(fw, tr)
	w = WrapResponseWriter(lw, w)

	w.WriteHeader(123)

	assert.Equal(t, 123, w.(*responseWriterPlus).
		ResponseWriter.(*logWriter).
		status)
}
