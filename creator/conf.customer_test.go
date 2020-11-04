package creator

import (
	"testing"

	"github.com/micro-plat/hydra/test/assert"
)

func Test_newCustomerBuilder(t *testing.T) {
	tests := []struct {
		name string
		args []interface{}
		want CustomerBuilder
	}{
		{name: "参数为空", args: []interface{}{}, want: CustomerBuilder{"main": make(map[string]interface{})}},
		{name: "单入参", args: []interface{}{"参数1"}, want: CustomerBuilder{"main": "参数1"}},
		{name: "多入参", args: []interface{}{"参数1", "参数2"}, want: CustomerBuilder{"main": "参数1"}},
	}
	for _, tt := range tests {
		got := newCustomerBuilder(tt.args...)
		assert.Equal(t, tt.want, got, tt.name)
	}
}

func TestCustomerBuilder_Sub(t *testing.T) {
	type args struct {
		name string
		s    []interface{}
	}
	tests := []struct {
		name string
		b    CustomerBuilder
		args args
		want ISUB
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		got := tt.b.Sub(tt.args.name, tt.args.s...)
		assert.Equal(t, tt.want, got, tt.name)
	}
}
