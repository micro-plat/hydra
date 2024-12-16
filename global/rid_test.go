package global

import "testing"

func BenchmarkCreate(b *testing.B) {
	for i := 0; i < b.N; i++ {
		RID.GetXRequestID()
	}
}
