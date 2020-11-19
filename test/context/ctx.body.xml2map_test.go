package context

import (
	"testing"

	"github.com/clbanning/mxj"
)

func BenchmarkNewMapXml(b *testing.B) {
	for i := 0; i < b.N; i++ {
		mxj.NewMapXml([]byte("<xml><a><c>1</c><a><b></b></xml>"))
	}
}
