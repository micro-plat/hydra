package circuit

import (
	"sync/atomic"
	"time"
)

type option struct {
	RPS        int
	FPPS       int
	RJTPS      int
	Timeout    int
	TimeWindow int64
}

//Option 配置选项
type Option func(*option)

//WithRPS 每秒请求数
func WithRPS(i int) Option {
	return func(o *option) {
		o.RPS = i
	}
}

//WithEPPS 每秒失败比例
func WithFPPS(i int) Option {
	return func(o *option) {
		o.FPPS = i
	}
}

//WithReject 每秒拒绝访问数
func WithReject(i int) Option {
	return func(o *option) {
		o.RJTPS = i
	}
}

//WithTimeout 每秒超时请求数
func WithTimeout(i int) Option {
	return func(o *option) {
		o.Timeout = i
	}
}

//WithSleepWindow 每秒失败数
func WithSleepWindow(i int64) Option {
	return func(o *option) {
		o.TimeWindow = i
	}
}

//CircuitBreaker 熔断管理
type CircuitBreaker struct {
	Name                   string
	open                   int32
	forceOpen              int32
	openedOrLastTestedTime int64
	metrics                *StandardMetricCollector
	*option
}

// NewCircuitBreaker creates a CircuitBreaker with associated Health
func NewCircuitBreaker(opts ...Option) *CircuitBreaker {
	c := &CircuitBreaker{
		forceOpen: -1,
		open:      -1,
		option: &option{
			RPS:        0,
			FPPS:       -1,
			RJTPS:      -1,
			TimeWindow: 10,
		},
	}
	for _, opt := range opts {
		opt(c.option)
	}
	c.metrics = NewStandardMetricCollector(c.TimeWindow)

	return c
}

// ToggleForceOpen allows manually causing the fallback logic for all instances
// of a given command.
func (circuit *CircuitBreaker) ToggleForceOpen(toggle bool) {
	if toggle {
		circuit.forceOpen = 0
		return
	}
	circuit.forceOpen = -1
}

// IsOpen is called before any Command execution to check whether or
// not it should be attempted. An "open" circuit means it is disabled.
func (circuit *CircuitBreaker) IsOpen() bool {
	return circuit.isOpen(time.Now())
}

//IsOpenByTime 指定时间范围内
func (circuit *CircuitBreaker) isOpen(now time.Time) bool {
	o := circuit.forceOpen == 0 || (circuit.open == 0 && circuit.TimeWindow > 0)
	if o {
		return true
	}
	if circuit.RPS == 0 || circuit.metrics.NumRequests().Sum(now) < uint64(circuit.RPS) {
		return false
	}

	if !circuit.IsHealthy(now) {
		circuit.setOpen()
	}
	return true
}

//GetCircuitStatus 获取熔断状态
func (circuit *CircuitBreaker) GetCircuitStatus() (isOpen bool, canRequest bool) {
	now := time.Now()
	isOpen = circuit.isOpen(now)
	canRequest = !isOpen || circuit.allowSingleTest(now)
	return
}

// AllowRequest is checked before a command executes, ensuring that circuit state and metric health allow it.
// When the circuit is open, this call will occasionally return true to measure whether the external service
// has recovered.
func (circuit *CircuitBreaker) AllowRequest() bool {
	return circuit.allowRequest(time.Now())
}
func (circuit *CircuitBreaker) allowRequest(now time.Time) bool {
	return !circuit.isOpen(now) || circuit.allowSingleTest(now)
}

func (circuit *CircuitBreaker) allowSingleTest(now time.Time) bool {
	if circuit.TimeWindow == 0 {
		return true
	}
	nowNano := now.UnixNano()
	openedOrLastTestedTime := atomic.LoadInt64(&circuit.openedOrLastTestedTime)
	if circuit.open == 0 && nowNano > openedOrLastTestedTime+circuit.TimeWindow*int64(time.Millisecond) {
		swapped := atomic.CompareAndSwapInt64(&circuit.openedOrLastTestedTime, openedOrLastTestedTime, nowNano)
		return swapped
	}
	return false
}

func (circuit *CircuitBreaker) setOpen() {
	if atomic.CompareAndSwapInt32(&circuit.open, -1, 0) {
		circuit.openedOrLastTestedTime = time.Now().UnixNano()
	}
}

func (circuit *CircuitBreaker) setClose() {
	if atomic.CompareAndSwapInt32(&circuit.open, 0, -1) {
		circuit.metrics.Reset()
	}
}

//IsHealthy 当前服务器健康状况
func (circuit *CircuitBreaker) IsHealthy(t time.Time) bool {
	return (circuit.FPPS < 0 || circuit.metrics.FailurePercent(t) > circuit.FPPS) && (circuit.RJTPS < 0 || circuit.metrics.RejectPercent(t) > circuit.RJTPS)
}

// ReportEvent records command metrics for tracking recent error rates and exposing data to the dashboard.
func (circuit *CircuitBreaker) ReportEvent(event string, i uint64) int64 {
	if event == EventSuccess && circuit.open == 0 {
		circuit.setClose()
	}
	switch event {
	case EventFailure:
		return circuit.metrics.Failure(i)
	case EventFallbackFailure:
		return circuit.metrics.FallbackFailure(i)
	case EventFallbackSuccess:
		return circuit.metrics.FallbackSuccess(i)
	case EventReject:
		return circuit.metrics.Reject(i)
	case EventShortCircuit:
		return circuit.metrics.ShortCircuit(i)
	case EventSuccess:
		return circuit.metrics.Success(i)
	case EventTimeout:
		return circuit.metrics.Timeout(i)

	}
	return 0
}

var (
	EventSuccess      = "SUCCESS"
	EventFailure      = "FAILURE"
	EventReject       = "REJECT"
	EventTimeout      = "TIMEOUT"
	EventShortCircuit = "SHORT_CIRCUIT"

	EventFallbackSuccess = "FALLBACK_SUCCESS"
	EventFallbackFailure = "FALLBACK_FAILURE"
)
