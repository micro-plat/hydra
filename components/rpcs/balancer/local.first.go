package balancer

import (
	"strings"
	"sync"

	rpcconf "github.com/micro-plat/hydra/conf/vars/rpc"
	"github.com/micro-plat/hydra/global"
	"google.golang.org/grpc/balancer"
	"google.golang.org/grpc/balancer/base"
)

// newBuilder creates a new roundrobin balancer builder.
func newBuilder(localip string) balancer.Builder {
	return base.NewBalancerBuilder(rpcconf.LocalFirst, &lfPickerBuilder{
		localip: localip,
	}, base.Config{HealthCheck: true})
}

func init() {
	balancer.Register(newBuilder(global.LocalIP()))
}

type lfPickerBuilder struct {
	localip string
}

func (builder *lfPickerBuilder) Build(info base.PickerBuildInfo) balancer.Picker {
	if len(info.ReadySCs) == 0 {
		return base.NewErrPicker(balancer.ErrNoSubConnAvailable)
	}
	var scs []balancer.SubConn
	var ifs []base.SubConnInfo
	for sc, ifv := range info.ReadySCs {
		scs = append(scs, sc)
		ifs = append(ifs, ifv)
	}

	return &lfPicker{
		subConns: scs,
		subInfos: ifs,
		next:     0,
		// Start at a random index, as the same RR balancer rebuilds a new
		// picker when SubConn states change, and we don't want to apply excess
		// load to the first server in the list.
		localip: builder.localip,
	}
}

type lfPicker struct {
	// subConns is the snapshot of the roundrobin balancer when this picker was
	// created. The slice is immutable. Each Get() will do a round robin
	// selection from it and return the selected SubConn.
	subConns []balancer.SubConn
	subInfos []base.SubConnInfo

	next    int
	mu      sync.Mutex
	localip string
}

func (p *lfPicker) Pick(info balancer.PickInfo) (balancer.PickResult, error) {
	p.mu.Lock()
	defer p.mu.Unlock()

	var hasFirst bool
	var idx int
	//检查是否有优先匹配项
	for i, v := range p.subInfos {
		idx = i
		hasFirst = strings.HasPrefix(v.Address.Addr, p.localip)
		if hasFirst {
			break
		}
	}
	if hasFirst {
		sc := p.subConns[idx]
		return balancer.PickResult{SubConn: sc}, nil
	}
	sc := p.subConns[p.next]
	p.next = (p.next + 1) % len(p.subConns)
	return balancer.PickResult{SubConn: sc}, nil

}
