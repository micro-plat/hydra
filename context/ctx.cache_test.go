package context

import (
	"testing"

	"github.com/micro-plat/hydra/global"
)

var tid string

func BenchmarkGetGID(b *testing.B) {
	b.ResetTimer()
	tid = global.GetGoroutineID()
	for i := 0; i < b.N; i++ {
		tid1 := global.GetGoroutineID()
		if tid != tid1 {
			b.Error("获取的数据有误")
			return
		}
	}
}
