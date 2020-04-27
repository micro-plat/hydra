package dispatcher

import (
	"net/http"
)

const (
	noWritten     = -1
	defaultStatus = 200
)

type ResponseWriter interface {

	// Returns the HTTP response status code of the current request.
	Status() int

	// Returns the number of bytes already written into the response http body.
	// See Written()
	Size() int
	Data() []byte
	// Writes the string into the response body.
	WriteString(string) (int, error)
	WriteHeader(code int)
	Write([]byte) (int, error)
	// Returns true if the response body was already written.
	Written() bool

	// Forces to write the http header (status code + headers).
	WriteHeaderNow()
	Header() http.Header
}

type responseWriter struct {
	size   int
	status int
	data   []byte
	header http.Header
}

func (w *responseWriter) Copy() *responseWriter {
	var cp = *w
	return &cp
}

func (w *responseWriter) reset() {
	w.size = noWritten
	w.header = make(map[string][]string)
	w.status = defaultStatus
	w.data = nil
}

func (w *responseWriter) WriteHeader(code int) {
	if code > 0 && w.status != code {
		w.status = code
	}
}

func (w *responseWriter) WriteHeaderNow() {
	if !w.Written() {
		w.size = 0
	}
}

func (w *responseWriter) Write(data []byte) (n int, err error) {
	w.WriteHeaderNow()
	w.data = data
	w.size += len(data)
	return w.size, nil
}
func (w *responseWriter) Header() http.Header {
	return w.header
}
func (w *responseWriter) WriteString(s string) (n int, err error) {
	w.WriteHeaderNow()
	w.data = []byte(s)
	w.size += len(w.data)
	return w.size, nil
}
func (w *responseWriter) Status() int {
	return w.status
}
func (w *responseWriter) Data() []byte {
	return w.data
}
func (w *responseWriter) Size() int {
	return w.size
}

func (w *responseWriter) Written() bool {
	return w.size != noWritten
}
