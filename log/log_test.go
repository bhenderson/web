package log

import (
	"testing"
)

func TestLogMiddleware(t *testing.T) {
	// don't panic
	LogMiddleware(nil, Common)
	LogMiddleware(nil, "")
}
