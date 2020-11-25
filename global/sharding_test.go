package global

import "testing"

type master struct {
	name          string
	master        bool
	shardingCount int
	path          string
	cldrs         []string
	expSharding   int
	expMaster     bool
}

func TestFirstStart(t *testing.T) {
	masters := []*master{
		&master{name: "1.1无子节点-默认false-shardingCount=0", master: false, shardingCount: 0, path: "abc", cldrs: []string{}, expSharding: 0, expMaster: false},
		&master{name: "1.2无子节点-默认false-shardingCount=1", master: false, shardingCount: 1, path: "abc", cldrs: []string{}, expSharding: 0, expMaster: false},
		&master{name: "1.2无子节点-默认false-shardingCount=大数", master: false, shardingCount: 100, path: "abc", cldrs: []string{}, expSharding: 0, expMaster: false},

		&master{name: "2.1无子节点-默认true-shardingCount=0", master: true, shardingCount: 0, path: "abc", cldrs: []string{}, expSharding: 0, expMaster: true},
		&master{name: "2.2无子节点-默认true-shardingCount=1", master: true, shardingCount: 1, path: "abc", cldrs: []string{}, expSharding: 0, expMaster: true},
		&master{name: "3.2无子节点-默认true-shardingCount=大数", master: true, shardingCount: 100, path: "abc", cldrs: []string{}, expSharding: 0, expMaster: true},

		&master{name: "3.1单子节点-默认false-shardingCount=0-存在", master: false, shardingCount: 0, path: "/abc/001", cldrs: []string{"abc_001"}, expSharding: 0, expMaster: true},
		&master{name: "3.2单子节点-默认false-shardingCount=1-存在", master: false, shardingCount: 1, path: "/abc/001", cldrs: []string{"abc_001"}, expSharding: 0, expMaster: true},
		&master{name: "3.3单子节点-默认false-shardingCount=大数-存在", master: false, shardingCount: 100, path: "/abc/001", cldrs: []string{"abc_001"}, expSharding: 0, expMaster: true},

		&master{name: "4.1单子节点-默认false-shardingCount=0-不存在", master: false, shardingCount: 0, path: "/abc/001", cldrs: []string{"abc_002"}, expSharding: 0, expMaster: false},
		&master{name: "4.2单子节点-默认false-shardingCount=1-不存在", master: false, shardingCount: 1, path: "/abc/001", cldrs: []string{"abc_002"}, expSharding: 0, expMaster: false},
		&master{name: "4.3单子节点-默认false-shardingCount=大数-不存在", master: false, shardingCount: 100, path: "/abc/001", cldrs: []string{"abc_002"}, expSharding: 0, expMaster: false},

		&master{name: "5.1单子节点-默认true-shardingCount=0-存在", master: true, shardingCount: 0, path: "/abc/001", cldrs: []string{"abc_001"}, expSharding: 0, expMaster: true},
		&master{name: "5.2单子节点-默认true-shardingCount=1-存在", master: true, shardingCount: 1, path: "/abc/001", cldrs: []string{"abc_001"}, expSharding: 0, expMaster: true},
		&master{name: "5.3单子节点-默认true-shardingCount=大数-存在", master: true, shardingCount: 100, path: "/abc/001", cldrs: []string{"abc_001"}, expSharding: 0, expMaster: true},

		&master{name: "6.1单子节点-默认true-shardingCount=0-不存在", master: true, shardingCount: 0, path: "/abc/001", cldrs: []string{"abc_002"}, expSharding: 0, expMaster: false},
		&master{name: "6.2单子节点-默认true-shardingCount=1-不存在", master: true, shardingCount: 1, path: "/abc/001", cldrs: []string{"abc_002"}, expSharding: 0, expMaster: false},
		&master{name: "6.3单子节点-默认true-shardingCount=大数-不存在", master: true, shardingCount: 100, path: "/abc/001", cldrs: []string{"abc_002"}, expSharding: 0, expMaster: false},

		&master{name: "7.1多子节点-默认false-shardingCount=0-存在-单个", master: false, shardingCount: 0, path: "/abc/001", cldrs: []string{"abc_001", "def_002", "ghi_003"}, expSharding: 0, expMaster: true},
		&master{name: "7.2多子节点-默认false-shardingCount=1-存在-单个", master: false, shardingCount: 1, path: "/abc/001", cldrs: []string{"abc_001", "def_002", "ghi_003"}, expSharding: 0, expMaster: true},
		&master{name: "7.3多子节点-默认false-shardingCount=大数-存在-单个", master: false, shardingCount: 100, path: "/abc/001", cldrs: []string{"abc_001", "def_002", "ghi_003"}, expSharding: 0, expMaster: true},

		&master{name: "8.1多子节点-默认false-shardingCount=0-存在-多个", master: false, shardingCount: 0, path: "/abc/001", cldrs: []string{"abc_001", "def_001", "ghi_001"}, expSharding: 0, expMaster: true},
		&master{name: "8.2多子节点-默认false-shardingCount=1-存在-多个", master: false, shardingCount: 1, path: "/abc/001", cldrs: []string{"abc_001", "def_001", "ghi_001"}, expSharding: 0, expMaster: true},
		&master{name: "8.3多子节点-默认false-shardingCount=大数-存在-多个", master: false, shardingCount: 100, path: "/abc/001", cldrs: []string{"abc_001", "def_001", "ghi_001"}, expSharding: 0, expMaster: true},

		&master{name: "9.1多子节点-默认false-shardingCount=0-不存在", master: false, shardingCount: 0, path: "/abc/002", cldrs: []string{"abc_001", "def_001", "ghi_001"}, expSharding: 0, expMaster: false},
		&master{name: "9.2多子节点-默认false-shardingCount=1-不存在", master: false, shardingCount: 1, path: "/abc/002", cldrs: []string{"abc_001", "def_001", "ghi_001"}, expSharding: 0, expMaster: false},
		&master{name: "9.3多子节点-默认false-shardingCount=大数-不存在", master: false, shardingCount: 100, path: "/abc/002", cldrs: []string{"abc_001", "def_001", "ghi_001"}, expSharding: 0, expMaster: false},

		&master{name: "10.1多子节点-默认true-shardingCount=0-存在-单个", master: true, shardingCount: 0, path: "/abc/001", cldrs: []string{"abc_001", "def_002", "ghi_003"}, expSharding: 0, expMaster: true},
		&master{name: "10.2多子节点-默认true-shardingCount=1-存在-单个", master: true, shardingCount: 1, path: "/abc/001", cldrs: []string{"abc_001", "def_002", "ghi_003"}, expSharding: 0, expMaster: true},
		&master{name: "10.3多子节点-默认true-shardingCount=大数-存在-单个", master: true, shardingCount: 100, path: "/abc/001", cldrs: []string{"abc_001", "def_002", "ghi_003"}, expSharding: 0, expMaster: true},

		&master{name: "11.1多子节点-默认true-shardingCount=0-存在-多个", master: true, shardingCount: 0, path: "/abc/001", cldrs: []string{"abc_001", "def_001", "ghi_001"}, expSharding: 0, expMaster: true},
		&master{name: "11.2多子节点-默认true-shardingCount=1-存在-多个", master: true, shardingCount: 1, path: "/abc/001", cldrs: []string{"abc_001", "def_001", "ghi_001"}, expSharding: 0, expMaster: true},
		&master{name: "11.3多子节点-默认true-shardingCount=大数-存在-多个", master: true, shardingCount: 100, path: "/abc/001", cldrs: []string{"abc_001", "def_001", "ghi_001"}, expSharding: 0, expMaster: true},

		&master{name: "12.1多子节点-默认true-shardingCount=0-不存在", master: true, shardingCount: 0, path: "/abc/002", cldrs: []string{"abc_001", "def_001", "ghi_001"}, expSharding: 0, expMaster: false},
		&master{name: "12.2多子节点-默认true-shardingCount=1-不存在", master: true, shardingCount: 1, path: "/abc/002", cldrs: []string{"abc_001", "def_001", "ghi_001"}, expSharding: 0, expMaster: false},
		&master{name: "12.3多子节点-默认true-shardingCount=大数-不存在", master: true, shardingCount: 100, path: "/abc/002", cldrs: []string{"abc_001", "def_001", "ghi_001"}, expSharding: 0, expMaster: false},
	}
	for _, m := range masters {
		actSharding, actMaster := IsMaster(m.master, m.shardingCount, m.path, m.cldrs)
		if actSharding != m.expSharding || actMaster != m.expMaster {
			t.Errorf("%s:请求结果不一致:[%d,%v]", m.name, actSharding, actMaster)
		}
	}
}
