package api

import "net/http"

func NewResponse(w http.ResponseWriter) *Response {
	return &Response{
		ResponseWriter: w,
		contentLength:  -1,
	}
}

type Response struct {
	http.ResponseWriter

	Status int
	Body   interface{}

	wroteHeader   bool
	contentLength int64
}

func (r *Response) WriteHeader(status int) {
	if status == 0 {
		status = http.StatusOK
	}
	r.ResponseWriter.WriteHeader(status)

	if r.wroteHeader {
		return
	}

	r.Status = status
	r.wroteHeader = true
}

func (r *Response) Write(p []byte) (int, error) {
	if !r.wroteHeader {
		r.WriteHeader(r.Status)
	}
	if r.contentLength == -1 {
		r.contentLength = 0
	}
	r.contentLength += int64(len(p))
	return r.ResponseWriter.Write(p)
}

func (r *Response) WriteString(p string) (int, error) {
	return r.Write([]byte(p))
}

func (r *Response) GetContentLength() int64 {
	return r.contentLength
}
