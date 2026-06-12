package context

import (
	"context"
	"net/http/httptest"
	"strings"
	"testing"
	"time"
)

func TestSSEDataCloseKeepsBufferedData(t *testing.T) {
	stream := NewSSEData(2)
	stream.Push("data: one")
	stream.Push("data: two")
	stream.Close()

	ok, first := stream.Pop()
	if !ok || first != "data: one" {
		t.Fatalf("first pop = (%v,%q), want buffered first item", ok, first)
	}
	ok, second := stream.Pop()
	if !ok || second != "data: two" {
		t.Fatalf("second pop = (%v,%q), want buffered second item", ok, second)
	}
	ok, _ = stream.Pop()
	if ok {
		t.Fatalf("pop after drained close = true, want false")
	}
}

func TestSSEDataLoopWriteWithContextStopsOnCancel(t *testing.T) {
	stream := NewSSEData(1)
	ctx, cancel := context.WithCancel(context.Background())
	done := make(chan struct{})
	recorder := httptest.NewRecorder()

	go func() {
		defer close(done)
		stream.LoopWriteWithContext(ctx, recorder)
	}()

	stream.Push("data: one")
	deadline := time.After(time.Second)
	for !strings.Contains(recorder.Body.String(), "data: one") {
		select {
		case <-deadline:
			t.Fatalf("response body %q does not contain pushed event", recorder.Body.String())
		default:
			time.Sleep(time.Millisecond)
		}
	}
	cancel()

	select {
	case <-done:
	case <-time.After(time.Second):
		t.Fatalf("LoopWriteWithContext did not stop after context cancellation")
	}
}

func TestSSEDataPushAfterCloseDoesNotBlock(t *testing.T) {
	stream := NewSSEData(0)
	stream.Close()
	done := make(chan struct{})
	go func() {
		defer close(done)
		stream.Push("data: ignored")
	}()

	select {
	case <-done:
	case <-time.After(time.Second):
		t.Fatalf("Push after Close blocked")
	}
}
