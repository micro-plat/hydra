package nfs

import (
	"fmt"
	"testing"

	"github.com/micro-plat/lib4go/assert"
)

func TestRead(t *testing.T) {
	m, err := readFiles("/users/ganyanfei/work/bin")
	assert.Equal(t, nil, err)
	fmt.Println(m)
	assert.Equal(t, 280, m.Len())

}
