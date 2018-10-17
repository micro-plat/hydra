package metrics

import "sync"

// IQPS count events to produce exponentially-weighted moving average rates
// at one-, five-, and fifteen-minutes and a mean rate.
type IQPS interface {
	Mark(int32)
	M1() int32
	M5() int32
	M15() int32
}

// GetOrRegisterQPS returns an existing Meter or constructs and registers a
// new StandardMeter.
func GetOrRegisterQPS(name string, r Registry) IQPS {
	if nil == r {
		r = DefaultRegistry
	}
	return r.GetOrRegister(name, newStandardQPS).(IQPS)
}

// StandardRPS is the standard implementation of a Meter.
type StandardRPS struct {
	lock sync.RWMutex
	m1   *QPSC
	m5   *QPSC
	m15  *QPSC
}

func newStandardQPS() IQPS {
	return &StandardRPS{
		m1:  NewQPSC(60, 70),
		m5:  NewQPSC(300, 350),
		m15: NewQPSC(900, 910),
	}
}

func (s *StandardRPS) Mark(i int32) {
	s.lock.Lock()
	s.m1.Mark(i)
	s.m5.Mark(i)
	s.m15.Mark(i)
	s.lock.Unlock()
}
func (s *StandardRPS) M1() int32 {
	return s.m1.counter
}
func (s *StandardRPS) M5() int32 {
	return s.m5.counter
}
func (s *StandardRPS) M15() int32 {
	return s.m15.counter
}
