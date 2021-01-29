package pkgs

import (
	"testing"

	"github.com/asaskevich/govalidator"
	"github.com/micro-plat/lib4go/assert"
)

type st struct {
	Service string `valid:"spath,required"`
}

func TestVDT(t *testing.T) {
	cases := []struct {
		input  *st
		result bool
	}{
		{input: &st{Service: "/a"}, result: true},
		{input: &st{Service: "b"}, result: false},
		{input: &st{Service: "http://abc"}, result: false},
	}
	for _, c := range cases {
		b, _ := govalidator.ValidateStruct(c.input)
		assert.Equal(t, c.result, b)
	}
}
