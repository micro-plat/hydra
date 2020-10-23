package registry

import (
	"testing"

	"github.com/micro-plat/hydra/registry/pub"
)

type master struct {
	master        bool
	shardingCount int
	path          string
	all           []string
	expSharding   int
	expMaster     bool
}

func TestFirstStart(t *testing.T) {
	masters := []*master{
		&master{master: false, shardingCount: 0, path: "abc", all: []string{}, expSharding: 0, expMaster: false},
		&master{master: false, shardingCount: 1, path: "abc", all: []string{}, expSharding: 0, expMaster: false},
		&master{master: false, shardingCount: 15, path: "abc", all: []string{}, expSharding: 0, expMaster: false},
		&master{master: false, shardingCount: 0, path: "abc", all: []string{"abc"}, expSharding: 0, expMaster: true},
		&master{master: false, shardingCount: 0, path: "abc", all: []string{"efg", "abc"}, expSharding: 0, expMaster: true},
		&master{master: false, shardingCount: 0, path: "abc", all: []string{"abc", "efg"}, expSharding: 0, expMaster: true},
		&master{master: false, shardingCount: 1, path: "/plat/00000_0002", all: []string{"192.168.0.100_0002", "192.168.0.100_0001"}, expSharding: 0, expMaster: false},
		&master{master: false, shardingCount: 1, path: "/plat/00000_0001", all: []string{"192.168.0.100_0002", "192.168.0.100_0001"}, expSharding: 0, expMaster: true},
		&master{master: false, shardingCount: 2, path: "/plat/00000_0002", all: []string{"192.168.0.100_0002", "192.168.0.100_0001"}, expSharding: 1, expMaster: true},
		&master{master: false, shardingCount: 2, path: "/plat/00000_0001", all: []string{"192.168.0.100_0002", "192.168.0.100_0001"}, expSharding: 0, expMaster: true},
		&master{master: false, shardingCount: 10, path: "/plat/00000_0002", all: []string{"192.168.0.100_0002", "192.168.0.100_0001"}, expSharding: 1, expMaster: true},
		&master{master: false, shardingCount: 10, path: "/plat/00000_0001", all: []string{"192.168.0.100_0002", "192.168.0.100_0001"}, expSharding: 0, expMaster: true},
		&master{master: false, shardingCount: 10, path: "/plat/00000_0002", all: []string{"192.168.0.100_0002", "192.168.0.100_0001"}, expSharding: 1, expMaster: true},
		&master{master: false, shardingCount: 10, path: "/plat/00000_0004", all: []string{"192.168.0.100_0002", "192.168.0.100_0001", "192.168.0.100_0003", "192.168.0.100_0004"}, expSharding: 3, expMaster: true},
		&master{master: false, shardingCount: 2, path: "/plat/00000_0004", all: []string{"192.168.0.100_0002", "192.168.0.100_0001", "192.168.0.100_0003", "192.168.0.100_0004"}, expSharding: 1, expMaster: true},
	}
	for _, m := range masters {
		actSharding, actMaster := pub.GetSharding(m.master, m.shardingCount, m.path, m.all)
		if actSharding != m.expSharding || actMaster != m.expMaster {
			t.Errorf("请求结果不一致:%v,[%d,%v]", m, actSharding, actMaster)
		}
	}
}

func TestFirstSeoncd(t *testing.T) {
	masters := []*master{
		&master{master: true, shardingCount: 0, path: "abc", all: []string{}, expSharding: 0, expMaster: true},
		&master{master: true, shardingCount: 1, path: "abc", all: []string{}, expSharding: 0, expMaster: true},
		&master{master: true, shardingCount: 15, path: "abc", all: []string{}, expSharding: 0, expMaster: true},
		&master{master: true, shardingCount: 0, path: "abc", all: []string{"abc"}, expSharding: 0, expMaster: true},
		&master{master: true, shardingCount: 0, path: "abc", all: []string{"efg", "abc"}, expSharding: 0, expMaster: true},
		&master{master: true, shardingCount: 0, path: "abc", all: []string{"abc", "efg"}, expSharding: 0, expMaster: true},
		&master{master: true, shardingCount: 1, path: "/plat/00000_0002", all: []string{"192.168.0.100_0002", "192.168.0.100_0001"}, expSharding: 0, expMaster: false},
		&master{master: true, shardingCount: 1, path: "/plat/00000_0001", all: []string{"192.168.0.100_0002", "192.168.0.100_0001"}, expSharding: 0, expMaster: true},
		&master{master: true, shardingCount: 2, path: "/plat/00000_0002", all: []string{"192.168.0.100_0002", "192.168.0.100_0001"}, expSharding: 1, expMaster: true},
		&master{master: true, shardingCount: 2, path: "/plat/00000_0001", all: []string{"192.168.0.100_0002", "192.168.0.100_0001"}, expSharding: 0, expMaster: true},
		&master{master: true, shardingCount: 10, path: "/plat/00000_0002", all: []string{"192.168.0.100_0002", "192.168.0.100_0001"}, expSharding: 1, expMaster: true},
		&master{master: true, shardingCount: 10, path: "/plat/00000_0001", all: []string{"192.168.0.100_0002", "192.168.0.100_0001"}, expSharding: 0, expMaster: true},
		&master{master: true, shardingCount: 10, path: "/plat/00000_0002", all: []string{"192.168.0.100_0002", "192.168.0.100_0001"}, expSharding: 1, expMaster: true},
		&master{master: true, shardingCount: 10, path: "/plat/00000_0004", all: []string{"192.168.0.100_0002", "192.168.0.100_0001", "192.168.0.100_0003", "192.168.0.100_0004"}, expSharding: 3, expMaster: true},
		&master{master: true, shardingCount: 2, path: "/plat/00000_0004", all: []string{"192.168.0.100_0002", "192.168.0.100_0001", "192.168.0.100_0003", "192.168.0.100_0004"}, expSharding: 1, expMaster: true},
	}
	for _, m := range masters {
		actSharding, actMaster := pub.GetSharding(m.master, m.shardingCount, m.path, m.all)
		if actSharding != m.expSharding || actMaster != m.expMaster {
			t.Errorf("请求结果不一致:%v,[%d,%v]", m, actSharding, actMaster)
		}
	}
}
