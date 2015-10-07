package web

import (
	"io"
	"net/http"
)

// AUTO GENERATED. DO NOT EDIT
// Usually, implementations of http.ResponseWriter contain extra methods, so we
// auto generate all combinations of those to keep the functionality

type plusser0 interface {
	http.CloseNotifier
	http.Flusher
	http.Hijacker
	io.ReaderFrom
}

type responseWriter0 struct {
	http.ResponseWriter
	plusser0
}

type plusser1 interface {
	http.CloseNotifier
	http.Flusher
	http.Hijacker
}

type responseWriter1 struct {
	http.ResponseWriter
	plusser1
}

type plusser2 interface {
	http.CloseNotifier
	http.Flusher
	io.ReaderFrom
}

type responseWriter2 struct {
	http.ResponseWriter
	plusser2
}

type plusser3 interface {
	http.CloseNotifier
	http.Hijacker
	io.ReaderFrom
}

type responseWriter3 struct {
	http.ResponseWriter
	plusser3
}

type plusser4 interface {
	http.CloseNotifier
	http.Flusher
}

type responseWriter4 struct {
	http.ResponseWriter
	plusser4
}

type plusser5 interface {
	http.CloseNotifier
	http.Hijacker
}

type responseWriter5 struct {
	http.ResponseWriter
	plusser5
}

type plusser6 interface {
	http.CloseNotifier
	io.ReaderFrom
}

type responseWriter6 struct {
	http.ResponseWriter
	plusser6
}

type plusser7 interface {
	http.Flusher
	http.Hijacker
	io.ReaderFrom
}

type responseWriter7 struct {
	http.ResponseWriter
	plusser7
}

type plusser8 interface {
	http.Flusher
	http.Hijacker
}

type responseWriter8 struct {
	http.ResponseWriter
	plusser8
}

type plusser9 interface {
	http.Flusher
	io.ReaderFrom
}

type responseWriter9 struct {
	http.ResponseWriter
	plusser9
}

type plusser10 interface {
	http.Hijacker
	io.ReaderFrom
}

type responseWriter10 struct {
	http.ResponseWriter
	plusser10
}

type plusser11 interface {
	http.CloseNotifier
}

type responseWriter11 struct {
	http.ResponseWriter
	plusser11
}

type plusser12 interface {
	http.Flusher
}

type responseWriter12 struct {
	http.ResponseWriter
	plusser12
}

type plusser13 interface {
	http.Hijacker
}

type responseWriter13 struct {
	http.ResponseWriter
	plusser13
}

type plusser14 interface {
	io.ReaderFrom
}

type responseWriter14 struct {
	http.ResponseWriter
	plusser14
}

func newPlusser(wr, pl http.ResponseWriter) http.ResponseWriter {
	switch x := pl.(type) {
	case plusser0:
		return &responseWriter0{wr, x}
	case plusser1:
		return &responseWriter1{wr, x}
	case plusser2:
		return &responseWriter2{wr, x}
	case plusser3:
		return &responseWriter3{wr, x}
	case plusser4:
		return &responseWriter4{wr, x}
	case plusser5:
		return &responseWriter5{wr, x}
	case plusser6:
		return &responseWriter6{wr, x}
	case plusser7:
		return &responseWriter7{wr, x}
	case plusser8:
		return &responseWriter8{wr, x}
	case plusser9:
		return &responseWriter9{wr, x}
	case plusser10:
		return &responseWriter10{wr, x}
	case plusser11:
		return &responseWriter11{wr, x}
	case plusser12:
		return &responseWriter12{wr, x}
	case plusser13:
		return &responseWriter13{wr, x}
	case plusser14:
		return &responseWriter14{wr, x}
	}
	return wr
}
