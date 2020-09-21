package context

import (
	"testing"
)

var tid uint64

func BenchmarkGetGID(b *testing.B) {
	b.ResetTimer()
	tid = getGID()
	for i := 0; i < b.N; i++ {
		tid1 := getGID()
		if tid != tid1 {
			b.Error("获取的数据有误")
			return
		}
	}
}
