package circuit

import (
	"testing"
	"time"
)

var num = NewSecondBucket(10)

func BenchmarkT1(b *testing.B) {

	for i := 0; i < b.N; i++ {
		num.Increment(uint64(i))
		num.Sum(time.Now())
	}
}
