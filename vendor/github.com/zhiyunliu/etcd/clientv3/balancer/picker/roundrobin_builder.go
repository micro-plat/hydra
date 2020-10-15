package picker

import (
	"sync"

	"go.uber.org/zap"
	"google.golang.org/grpc/balancer"
	"google.golang.org/grpc/balancer/base"
)

type rrPickerBuilder struct {
}

func (b *rrPickerBuilder) Build(info base.PickerBuildInfo) balancer.Picker {
	if len(info.ReadySCs) == 0 {
		return base.NewErrPicker(balancer.ErrNoSubConnAvailable)
	}
	var scs []balancer.SubConn
	for sc := range info.ReadySCs {
		scs = append(scs, sc)
	}
	return &rrPicker{
		subConns: scs,
		// Start at a random index, as the same RR balancer rebuilds a new
		// picker when SubConn states change, and we don't want to apply excess
		// load to the first server in the list.
		next: 0,
	}
}

type rrPicker struct {
	lg *zap.Logger

	// subConns is the snapshot of the roundrobin balancer when this picker was
	// created. The slice is immutable. Each Get() will do a round robin
	// selection from it and return the selected SubConn.
	subConns []balancer.SubConn

	mu   sync.Mutex
	next int
}

func (p *rrPicker) Pick(balancer.PickInfo) (balancer.PickResult, error) {
	p.mu.Lock()
	sc := p.subConns[p.next]
	p.next = (p.next + 1) % len(p.subConns)
	p.mu.Unlock()

	doneFunc := func(info balancer.DoneInfo) {
	}

	return balancer.PickResult{
		SubConn: sc,
		Done:    doneFunc,
	}, nil
}
