package web

import (
	"io"
	"net/http"
)

// AUTO GENERATED. DO NOT EDIT
// Usually, implementations of http.ResponseWriter contain extra methods, so we
// auto generate all combinations of those to keep the functionality

type plusser1111 interface {
	http.CloseNotifier
	http.Flusher
	http.Hijacker
	io.ReaderFrom
}

type responseWriter1111 struct {
	http.ResponseWriter
	plusser1111
}

type plusser1110 interface {
	http.CloseNotifier
	http.Flusher
	http.Hijacker
	// io.ReaderFrom
}

type responseWriter1110 struct {
	http.ResponseWriter
	plusser1110
}

type plusser1101 interface {
	http.CloseNotifier
	http.Flusher
	// http.Hijacker
	io.ReaderFrom
}

type responseWriter1101 struct {
	http.ResponseWriter
	plusser1101
}

type plusser1100 interface {
	http.CloseNotifier
	http.Flusher
	// http.Hijacker
	// io.ReaderFrom
}

type responseWriter1100 struct {
	http.ResponseWriter
	plusser1100
}

type plusser1011 interface {
	http.CloseNotifier
	// http.Flusher
	http.Hijacker
	io.ReaderFrom
}

type responseWriter1011 struct {
	http.ResponseWriter
	plusser1011
}

type plusser1010 interface {
	http.CloseNotifier
	// http.Flusher
	http.Hijacker
	// io.ReaderFrom
}

type responseWriter1010 struct {
	http.ResponseWriter
	plusser1010
}

type plusser1001 interface {
	http.CloseNotifier
	// http.Flusher
	// http.Hijacker
	io.ReaderFrom
}

type responseWriter1001 struct {
	http.ResponseWriter
	plusser1001
}

type plusser1000 interface {
	http.CloseNotifier
	// http.Flusher
	// http.Hijacker
	// io.ReaderFrom
}

type responseWriter1000 struct {
	http.ResponseWriter
	plusser1000
}

type plusser0111 interface {
	// http.CloseNotifier
	http.Flusher
	http.Hijacker
	io.ReaderFrom
}

type responseWriter0111 struct {
	http.ResponseWriter
	plusser0111
}

type plusser0110 interface {
	// http.CloseNotifier
	http.Flusher
	http.Hijacker
	// io.ReaderFrom
}

type responseWriter0110 struct {
	http.ResponseWriter
	plusser0110
}

type plusser0101 interface {
	// http.CloseNotifier
	http.Flusher
	// http.Hijacker
	io.ReaderFrom
}

type responseWriter0101 struct {
	http.ResponseWriter
	plusser0101
}

type plusser0100 interface {
	// http.CloseNotifier
	http.Flusher
	// http.Hijacker
	// io.ReaderFrom
}

type responseWriter0100 struct {
	http.ResponseWriter
	plusser0100
}

type plusser0011 interface {
	// http.CloseNotifier
	// http.Flusher
	http.Hijacker
	io.ReaderFrom
}

type responseWriter0011 struct {
	http.ResponseWriter
	plusser0011
}

type plusser0010 interface {
	// http.CloseNotifier
	// http.Flusher
	http.Hijacker
	// io.ReaderFrom
}

type responseWriter0010 struct {
	http.ResponseWriter
	plusser0010
}

type plusser0001 interface {
	// http.CloseNotifier
	// http.Flusher
	// http.Hijacker
	io.ReaderFrom
}

type responseWriter0001 struct {
	http.ResponseWriter
	plusser0001
}

type plusser0000 interface {
	// http.CloseNotifier
	// http.Flusher
	// http.Hijacker
	// io.ReaderFrom
}

type responseWriter0000 struct {
	http.ResponseWriter
	plusser0000
}

func newPlusser(wr, pl http.ResponseWriter) http.ResponseWriter {
	switch x := pl.(type) {
	case plusser1111:
		return &responseWriter1111{wr, x}
	case plusser1110:
		return &responseWriter1110{wr, x}
	case plusser1101:
		return &responseWriter1101{wr, x}
	case plusser1100:
		return &responseWriter1100{wr, x}
	case plusser1011:
		return &responseWriter1011{wr, x}
	case plusser1010:
		return &responseWriter1010{wr, x}
	case plusser1001:
		return &responseWriter1001{wr, x}
	case plusser1000:
		return &responseWriter1000{wr, x}
	case plusser0111:
		return &responseWriter0111{wr, x}
	case plusser0110:
		return &responseWriter0110{wr, x}
	case plusser0101:
		return &responseWriter0101{wr, x}
	case plusser0100:
		return &responseWriter0100{wr, x}
	case plusser0011:
		return &responseWriter0011{wr, x}
	case plusser0010:
		return &responseWriter0010{wr, x}
	case plusser0001:
		return &responseWriter0001{wr, x}
	case plusser0000:
		return &responseWriter0000{wr, x}
	}
	return wr
}
