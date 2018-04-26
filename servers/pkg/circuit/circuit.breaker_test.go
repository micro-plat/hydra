package circuit

import (
	"testing"
	"time"

	"github.com/qxnw/lib4go/ut"
)

func TestForceOpen1(t *testing.T) {
	breaker := NewCircuitBreaker()
	breaker.ToggleForceOpen(true)
	ut.Expect(t, breaker.IsOpen(), true)
	ut.Expect(t, breaker.AllowRequest(), false)
}
func TestForceOpen2(t *testing.T) {
	breaker := NewCircuitBreaker()
	breaker.ToggleForceOpen(false)
	ut.Expect(t, breaker.IsOpen(), false)
	ut.Expect(t, breaker.AllowRequest(), true)
}
func TestOpen3(t *testing.T) {
	breaker := NewCircuitBreaker()
	ut.Expect(t, breaker.IsOpen(), false)
	ut.Expect(t, breaker.AllowRequest(), true)
}
func TestRequestSuccess(t *testing.T) {
	breaker := NewCircuitBreaker()
	var value uint64 = 10
	now := breaker.ReportEvent(EventSuccess, value)
	ut.Expect(t, breaker.metrics.NumRequests().Sum(time.Unix(now, 0)), value)
	ut.Expect(t, breaker.metrics.Successes().Sum(time.Unix(now, 0)), value)

}
func TestRequestFailed(t *testing.T) {
	breaker := NewCircuitBreaker()
	var value uint64 = 10
	var def uint64
	now := breaker.ReportEvent(EventFailure, value)
	ut.Expect(t, breaker.metrics.NumRequests().Sum(time.Unix(now, 0)), value)
	ut.Expect(t, breaker.metrics.Failures().Sum(time.Unix(now, 0)), value)
	ut.Expect(t, breaker.metrics.Successes().Sum(time.Unix(now, 0)), def)
}
func TestRequestTimeout(t *testing.T) {
	breaker := NewCircuitBreaker()
	var value uint64 = 10
	var def uint64
	now := breaker.ReportEvent(EventTimeout, value)
	ut.Expect(t, breaker.metrics.NumRequests().Sum(time.Unix(now, 0)), value)
	ut.Expect(t, breaker.metrics.Timeouts().Sum(time.Unix(now, 0)), value)
	ut.Expect(t, breaker.metrics.Failures().Sum(time.Unix(now, 0)), value)
	ut.Expect(t, breaker.metrics.Successes().Sum(time.Unix(now, 0)), def)
}
func TestRequestReject(t *testing.T) {
	breaker := NewCircuitBreaker()
	var value uint64 = 10
	var def uint64
	now := breaker.ReportEvent(EventReject, value)
	ut.Expect(t, breaker.metrics.NumRequests().Sum(time.Unix(now, 0)), value)
	ut.Expect(t, breaker.metrics.Rejects().Sum(time.Unix(now, 0)), value)
	ut.Expect(t, breaker.metrics.Failures().Sum(time.Unix(now, 0)), value)
	ut.Expect(t, breaker.metrics.Successes().Sum(time.Unix(now, 0)), def)

}

func TestRequestCicuit(t *testing.T) {
	breaker := NewCircuitBreaker()
	var value uint64 = 10
	var def uint64
	now := breaker.ReportEvent(EventShortCircuit, value)
	ut.Expect(t, breaker.metrics.NumRequests().Sum(time.Unix(now, 0)), def)
	ut.Expect(t, breaker.metrics.ShortCircuits().Sum(time.Unix(now, 0)), value)
	ut.Expect(t, breaker.metrics.Failures().Sum(time.Unix(now, 0)), def)
}
func TestRequestFallbackSuccess(t *testing.T) {
	breaker := NewCircuitBreaker()
	var value uint64 = 10
	var def uint64
	now := breaker.ReportEvent(EventFallbackSuccess, value)
	ut.Expect(t, breaker.metrics.NumRequests().Sum(time.Unix(now, 0)), def)
	ut.Expect(t, breaker.metrics.ShortCircuits().Sum(time.Unix(now, 0)), value)
	ut.Expect(t, breaker.metrics.Failures().Sum(time.Unix(now, 0)), def)
	ut.Expect(t, breaker.metrics.FallbackSuccesses().Sum(time.Unix(now, 0)), value)
	ut.Expect(t, breaker.metrics.FallbackFailures().Sum(time.Unix(now, 0)), def)
}

func TestRequestFallbackFailed(t *testing.T) {
	breaker := NewCircuitBreaker()
	var value uint64 = 10
	var def uint64
	now := breaker.ReportEvent(EventFallbackFailure, value)
	ut.Expect(t, breaker.metrics.NumRequests().Sum(time.Unix(now, 0)), def)
	ut.Expect(t, breaker.metrics.ShortCircuits().Sum(time.Unix(now, 0)), value)
	ut.Expect(t, breaker.metrics.Failures().Sum(time.Unix(now, 0)), def)
	ut.Expect(t, breaker.metrics.FallbackFailures().Sum(time.Unix(now, 0)), value)
	ut.Expect(t, breaker.metrics.FallbackSuccesses().Sum(time.Unix(now, 0)), def)
}

func TestRequestRequestCircuit(t *testing.T) {
	breaker := NewCircuitBreaker(WithRPS(10))
	var value uint64 = 11
	now := breaker.ReportEvent(EventSuccess, value)
	ut.Expect(t, breaker.isOpen(time.Unix(now, 0)), true)
	ut.Expect(t, breaker.allowRequest(time.Unix(now, 0)), false)
}
func TestRequestRequestCircuitPass(t *testing.T) {
	breaker := NewCircuitBreaker(WithRPS(10))
	var value uint64 = 11
	now := breaker.ReportEvent(EventSuccess, value) + 11
	ut.Expect(t, breaker.isOpen(time.Unix(now, 0)), false)
	ut.Expect(t, breaker.allowSingleTest(time.Unix(now, 0)), false)
}
func TestRequestErrorCircuit(t *testing.T) {
	breaker := NewCircuitBreaker(WithRPS(10))
	var value uint64 = 11
	now := breaker.ReportEvent(EventSuccess, value) + 11
	breaker.ReportEvent(EventFailure, value)
	ut.Expect(t, breaker.isOpen(time.Unix(now, 0)), false)
	ut.Expect(t, breaker.allowSingleTest(time.Unix(now, 0)), false)
}
