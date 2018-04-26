package conf

import (
	"fmt"
	"testing"

	"github.com/micro-plat/lib4go/ut"
)

func TestCryptNow(t *testing.T) {
	input := "hydra杨"
	v := encrypt([]byte(input))
	fmt.Println(v)
	r, err := decrypt([]byte(v))
	ut.Expect(t, err, nil)
	ut.Expect(t, string(r), input)
}
func TestCryptNow2(t *testing.T) {
	input := "hydra杨"
	r, err := decrypt([]byte(input))
	ut.Expect(t, err, nil)
	ut.Expect(t, string(r), input)
}
