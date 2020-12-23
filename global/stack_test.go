package global

import (
	"testing"
)

var tid string

func BenchmarkGetGID(b *testing.B) {
	tid = GetGoroutineID()
	for i := 0; i < b.N; i++ {
		tid1 := GetGoroutineID()
		if tid != tid1 {
			b.Error("获取的数据有误")
			return
		}
	}
}
