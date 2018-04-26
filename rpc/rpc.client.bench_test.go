package rpc

import "testing"
import "github.com/qxnw/lib4go/ut"

var invoker = NewRPCInvoker("/hydra", "merchant_api", "zk://192.168.106.172")

func init() {

}
func BenchmarkRPCClient1(b *testing.B) {
	invoker.PreInit("/order/request/success")
	//invoker.Request("/order/request/success", make(map[string]string), true)
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		s, _, _, _ := invoker.Request("/order/request/success", make(map[string]string), true)
		ut.Expect(b, s, 200)
	}
}
