package balancer

import (
	"google.golang.org/grpc/resolver"
)

type mockSubConn struct {
	Addr string
}

func (sc *mockSubConn) UpdateAddresses([]resolver.Address) {
	return
}

func (sc *mockSubConn) Connect() {
	return
}
func (sc *mockSubConn) String() string {
	return sc.Addr
}

// func TestLocalFirst(t *testing.T) {
// 	localip := "172.0.0.1"
// 	pickerbuilder := &lfPickerBuilder{
// 		localip: localip,
// 	}

// 	scs := map[balancer.SubConn]base.SubConnInfo{}
// 	scs[&mockSubConn{
// 		Addr: "1",
// 	}] = base.SubConnInfo{
// 		Address: resolver.Address{Addr: "172.0.0.1"},
// 	}
// 	scs[&mockSubConn{
// 		Addr: "2",
// 	}] = base.SubConnInfo{
// 		Address: resolver.Address{Addr: "173.0.0.1"},
// 	}

// 	info := base.PickerBuildInfo{
// 		ReadySCs: scs,
// 	}
// 	picker := pickerbuilder.Build(info)

// 	pickInfo := balancer.PickInfo{}

// 	result, err := picker.Pick(pickInfo)
// 	if err != nil {
// 		t.Error(err)
// 	}
// 	if fmt.Sprintf("%v", result.SubConn) != "1" {
// 		t.Error("平衡策略失败0")
// 	}

// 	result1, err := picker.Pick(pickInfo)
// 	if err != nil {
// 		t.Error(err)
// 	}
// 	if fmt.Sprintf("%v", result1.SubConn) != "1" {
// 		t.Error("平衡策略失败1")
// 	}

// 	result2, err := picker.Pick(pickInfo)
// 	if err != nil {
// 		t.Error(err)
// 	}
// 	if fmt.Sprintf("%v", result2.SubConn) != "1" {
// 		t.Error("平衡策略失败2")
// 	}
// }
