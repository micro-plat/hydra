package circuit

import "time"

type StandardMetricCollector struct {
	timeRange   int64
	numRequests *SecondBucket
	successes   *SecondBucket
	failures    *SecondBucket
	rejects     *SecondBucket
	timeout     *SecondBucket

	shortCircuits *SecondBucket

	fallbackSuccesses *SecondBucket
	fallbackFailures  *SecondBucket
}

func NewStandardMetricCollector(timeRange int64) *StandardMetricCollector {
	d := &StandardMetricCollector{timeRange: timeRange}
	d.Reset()
	return d
}

func (d *StandardMetricCollector) NumRequests() *SecondBucket {
	return d.numRequests
}

// Successes returns the rolling number of successes
func (d *StandardMetricCollector) Successes() *SecondBucket {
	return d.successes
}

// Failures returns the rolling number of failures
func (d *StandardMetricCollector) Failures() *SecondBucket {
	return d.failures
}

// Rejects returns the rolling number of rejects
func (d *StandardMetricCollector) Rejects() *SecondBucket {
	return d.rejects
}

// Timeouts returns the rolling number of rejects
func (d *StandardMetricCollector) Timeouts() *SecondBucket {
	return d.timeout
}

// ShortCircuits returns the rolling number of short circuits
func (d *StandardMetricCollector) ShortCircuits() *SecondBucket {
	return d.shortCircuits
}

// FallbackSuccesses returns the rolling number of fallback successes
func (d *StandardMetricCollector) FallbackSuccesses() *SecondBucket {
	return d.fallbackSuccesses
}

// FallbackFailures returns the rolling number of fallback failures
func (d *StandardMetricCollector) FallbackFailures() *SecondBucket {
	return d.fallbackFailures
}

//Success 成功记数
func (d *StandardMetricCollector) Success(i uint64) (t int64) {
	d.numRequests.Increment(i)
	return d.successes.Increment(i)
}

//Failure 失败记数
func (d *StandardMetricCollector) Failure(i uint64) (t int64) {
	d.numRequests.Increment(i)
	return d.failures.Increment(i)
}

//Reject 拒绝访问
func (d *StandardMetricCollector) Reject(i uint64) (t int64) {
	d.numRequests.Increment(i)
	d.failures.Increment(i)
	return d.rejects.Increment(i)
}

//Timeout 超时请求
func (d *StandardMetricCollector) Timeout(i uint64) (t int64) {
	d.numRequests.Increment(i)
	d.failures.Increment(i)
	return d.timeout.Increment(i)
}

//ShortCircuit 熔断记数
func (d *StandardMetricCollector) ShortCircuit(i uint64) int64 {
	return d.shortCircuits.Increment(i)
}

//FallbackSuccess 熔断执行成功记数
func (d *StandardMetricCollector) FallbackSuccess(i uint64) int64 {
	d.shortCircuits.Increment(i)
	return d.fallbackSuccesses.Increment(i)
}

//FallbackFailure 熔断执行失败记数
func (d *StandardMetricCollector) FallbackFailure(i uint64) int64 {
	d.shortCircuits.Increment(i)
	return d.fallbackFailures.Increment(i)
}

func (m *StandardMetricCollector) FailurePercent(now time.Time) int {
	var errPct float64
	request := m.numRequests.Sum(now)
	failure := m.failures.Sum(now)
	if request > 0 {
		errPct = (float64(failure) / float64(request)) * 100
	}
	return int(errPct + 0.5)
}

func (m *StandardMetricCollector) RejectPercent(now time.Time) int {
	var errPct float64
	request := m.numRequests.Sum(now)
	failure := m.rejects.Sum(now)
	if request > 0 {
		errPct = (float64(failure) / float64(request)) * 100
	}
	return int(errPct + 0.5)
}

//Reset resets all metrics in this collector to 0.
func (d *StandardMetricCollector) Reset() {
	d.numRequests = NewSecondBucket(d.timeRange)
	d.successes = NewSecondBucket(d.timeRange)
	d.rejects = NewSecondBucket(d.timeRange)
	d.timeout = NewSecondBucket(d.timeRange)
	d.shortCircuits = NewSecondBucket(d.timeRange)
	d.failures = NewSecondBucket(d.timeRange)
	d.fallbackSuccesses = NewSecondBucket(d.timeRange)
	d.fallbackFailures = NewSecondBucket(d.timeRange)
}
