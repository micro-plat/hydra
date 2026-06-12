package aimiddleware

import (
	"context"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/micro-plat/hydra/pkgs"
)

type blockingReadCloser struct {
	closed chan struct{}
}

func newBlockingReadCloser() *blockingReadCloser {
	return &blockingReadCloser{closed: make(chan struct{})}
}

func (r *blockingReadCloser) Read(_ []byte) (int, error) {
	<-r.closed
	return 0, errors.New("closed")
}

func (r *blockingReadCloser) Close() error {
	select {
	case <-r.closed:
	default:
		close(r.closed)
	}
	return nil
}

func TestRPCErrorStatusHandlesNilResponse(t *testing.T) {
	if got := rpcErrorStatus(nil); got != http.StatusBadGateway {
		t.Fatalf("rpcErrorStatus(nil) = %d, want %d", got, http.StatusBadGateway)
	}
	if got := rpcErrorStatus(pkgs.NewRspnsByHD(http.StatusGatewayTimeout, "{}", nil)); got != http.StatusGatewayTimeout {
		t.Fatalf("rpcErrorStatus(timeout response) = %d, want %d", got, http.StatusGatewayTimeout)
	}
}

func TestCopySSEProxyClosesBodyOnRequestCancel(t *testing.T) {
	body := newBlockingReadCloser()
	ctx, cancel := context.WithCancel(context.Background())
	recorder := httptest.NewRecorder()
	done := make(chan error, 1)

	go func() {
		done <- copySSEProxy(ctx, recorder, SSEProxy{StatusCode: http.StatusOK, Body: body}, time.Minute)
	}()
	cancel()

	select {
	case <-body.closed:
	case <-time.After(time.Second):
		t.Fatalf("body was not closed after request cancellation")
	}
	select {
	case <-done:
	case <-time.After(time.Second):
		t.Fatalf("copySSEProxy did not return after request cancellation")
	}
}

func TestCopySSEProxyClosesBodyOnIdleTimeout(t *testing.T) {
	body := newBlockingReadCloser()
	recorder := httptest.NewRecorder()
	done := make(chan error, 1)

	go func() {
		done <- copySSEProxy(context.Background(), recorder, SSEProxy{StatusCode: http.StatusOK, Body: body}, 20*time.Millisecond)
	}()

	select {
	case <-body.closed:
	case <-time.After(time.Second):
		t.Fatalf("body was not closed after stream idle timeout")
	}
	select {
	case <-done:
	case <-time.After(time.Second):
		t.Fatalf("copySSEProxy did not return after stream idle timeout")
	}
}
