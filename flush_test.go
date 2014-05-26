package web

import (
	"github.com/stretchr/testify/assert"
	"net/http"
	"net/http/httptest"
	"testing"
)

func Test_WrapFlushWriter(t *testing.T) {
	var w http.ResponseWriter = httptest.NewRecorder()
	tr := NewTestResponse()

	lw := &logWriter{w, 200}
	fw := &flushWriter{w}
	w = WrapResponseWriter(lw, tr)
	w = WrapResponseWriter(fw, w)

	w.Write([]byte{0, 1, 2})

	assert.True(t, w.(*responseWriterPlus).
		ResponseWriter.(*flushWriter).
		ResponseWriter.(*httptest.ResponseRecorder).
		Flushed)
}
