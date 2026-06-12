package wss

import (
	"fmt"
	"net/http"
	"sync"

	componentwss "github.com/micro-plat/hydra/components/wss"
	"github.com/micro-plat/hydra/hydra/servers/pkg/dispatcher"
)

type frameStreamWriter struct {
	id      string
	send    func(*Frame) error
	header  http.Header
	status  int
	size    int
	written bool
	mu      sync.Mutex
}

func newFrameStreamWriter(id string, send func(*Frame) error) *frameStreamWriter {
	return &frameStreamWriter{
		id:     id,
		send:   send,
		header: make(http.Header),
		status: http.StatusOK,
		size:   -1,
	}
}

func (w *frameStreamWriter) Status() int {
	w.mu.Lock()
	defer w.mu.Unlock()
	return w.status
}

func (w *frameStreamWriter) Size() int {
	w.mu.Lock()
	defer w.mu.Unlock()
	return w.size
}

func (w *frameStreamWriter) Data() []byte {
	return nil
}

func (w *frameStreamWriter) Header() http.Header {
	return w.header
}

func (w *frameStreamWriter) Written() bool {
	w.mu.Lock()
	defer w.mu.Unlock()
	return w.written
}

func (w *frameStreamWriter) WriteHeader(code int) {
	w.mu.Lock()
	if code > 0 {
		w.status = code
	}
	w.mu.Unlock()
}

func (w *frameStreamWriter) WriteHeaderNow() {
	_ = w.writeStart()
}

func (w *frameStreamWriter) WriteString(s string) (int, error) {
	return w.Write([]byte(s))
}

func (w *frameStreamWriter) Write(data []byte) (int, error) {
	if err := w.writeStart(); err != nil {
		return 0, err
	}
	body := make([]byte, len(data))
	copy(body, data)
	if err := w.send(&Frame{Type: frameStreamChunk, ID: w.id, Body: body}); err != nil {
		return 0, err
	}
	w.mu.Lock()
	if w.size < 0 {
		w.size = 0
	}
	w.size += len(data)
	w.mu.Unlock()
	return len(data), nil
}

func (w *frameStreamWriter) Flush() {}

func (w *frameStreamWriter) writeStart() error {
	w.mu.Lock()
	if w.written {
		w.mu.Unlock()
		return nil
	}
	w.written = true
	w.size = 0
	status := w.status
	header := headers(w.header)
	w.mu.Unlock()
	return w.send(&Frame{Type: frameStreamStart, ID: w.id, Status: status, Header: header})
}

type httpStreamWriter struct {
	dst     http.ResponseWriter
	header  http.Header
	status  int
	size    int
	written bool
}

func newHTTPStreamWriter(dst http.ResponseWriter) *httpStreamWriter {
	return &httpStreamWriter{
		dst:    dst,
		header: dst.Header(),
		status: http.StatusOK,
		size:   -1,
	}
}

func (w *httpStreamWriter) Status() int {
	return w.status
}

func (w *httpStreamWriter) Size() int {
	return w.size
}

func (w *httpStreamWriter) Data() []byte {
	return nil
}

func (w *httpStreamWriter) Header() http.Header {
	return w.header
}

func (w *httpStreamWriter) Written() bool {
	return w.written
}

func (w *httpStreamWriter) WriteHeader(code int) {
	if code > 0 {
		w.status = code
	}
}

func (w *httpStreamWriter) WriteHeaderNow() {
	if w.written {
		return
	}
	w.written = true
	w.size = 0
	w.dst.WriteHeader(statusOr(w.status, http.StatusOK))
}

func (w *httpStreamWriter) WriteString(s string) (int, error) {
	return w.Write([]byte(s))
}

func (w *httpStreamWriter) Write(data []byte) (int, error) {
	w.WriteHeaderNow()
	n, err := w.dst.Write(data)
	w.size += n
	w.Flush()
	return n, err
}

func (w *httpStreamWriter) Flush() {
	if f, ok := w.dst.(http.Flusher); ok {
		f.Flush()
	}
}

func writeStreamEnd(writer *frameStreamWriter, status int, err error) error {
	frame := &Frame{Type: frameStreamEnd, ID: writer.id, Status: statusOr(status, writer.Status()), Header: headers(writer.Header())}
	if err != nil {
		frame.Error = err.Error()
	}
	if writer.send == nil {
		return fmt.Errorf("wss stream writer is closed")
	}
	return writer.send(frame)
}

var _ dispatcher.ResponseWriter = (*frameStreamWriter)(nil)
var _ dispatcher.ResponseWriter = (*httpStreamWriter)(nil)
var _ http.ResponseWriter = (*frameStreamWriter)(nil)
var _ http.Flusher = (*frameStreamWriter)(nil)
var _ http.ResponseWriter = (*httpStreamWriter)(nil)
var _ http.Flusher = (*httpStreamWriter)(nil)

type componentStreamWriter struct {
	onChunk componentwss.ChunkHandler
	onStart componentwss.StartHandler
	header  http.Header
	status  int
	size    int
	written bool
}

func newComponentStreamWriter(onStart componentwss.StartHandler, onChunk componentwss.ChunkHandler) *componentStreamWriter {
	return &componentStreamWriter{
		onChunk: onChunk,
		onStart: onStart,
		header:  make(http.Header),
		status:  http.StatusOK,
		size:    -1,
	}
}

func (w *componentStreamWriter) Status() int {
	return w.status
}

func (w *componentStreamWriter) Size() int {
	return w.size
}

func (w *componentStreamWriter) Data() []byte {
	return nil
}

func (w *componentStreamWriter) Header() http.Header {
	return w.header
}

func (w *componentStreamWriter) Written() bool {
	return w.written
}

func (w *componentStreamWriter) WriteHeader(code int) {
	if code > 0 {
		w.status = code
	}
}

func (w *componentStreamWriter) WriteHeaderNow() {
	if w.written {
		return
	}
	w.written = true
	w.size = 0
	if w.onStart != nil {
		_ = w.onStart(&componentwss.Response{Status: statusOr(w.status, http.StatusOK), Header: componentwss.Header(headers(w.header))})
	}
}

func (w *componentStreamWriter) WriteString(s string) (int, error) {
	return w.Write([]byte(s))
}

func (w *componentStreamWriter) Write(data []byte) (int, error) {
	w.WriteHeaderNow()
	if w.onChunk != nil && len(data) > 0 {
		if err := w.onChunk(data); err != nil {
			return 0, err
		}
	}
	w.size += len(data)
	return len(data), nil
}

func (w *componentStreamWriter) Flush() {}

var _ dispatcher.ResponseWriter = (*componentStreamWriter)(nil)
var _ http.ResponseWriter = (*componentStreamWriter)(nil)
var _ http.Flusher = (*componentStreamWriter)(nil)
