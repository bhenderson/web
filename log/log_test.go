package log

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestLogMiddleware(t *testing.T) {
	// don't panic
	tmp := ""
	assert.NotPanics(t, func() {
		LogMiddleware(discard, tmp)
	}, "expected %s to not panic", tmp)

	tmp = "{{.DoesNotExist}}"
	assert.Panics(t, func() {
		LogMiddleware(discard, tmp)
	}, "expected %s to panic", tmp)
}
