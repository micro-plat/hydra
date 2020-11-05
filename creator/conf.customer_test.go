package creator

import (
	"encoding/json"
	"strings"
	"testing"

	"github.com/micro-plat/hydra/test/assert"
	"github.com/micro-plat/lib4go/types"
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
		name    string
		b       CustomerBuilder
		args    args
		want    CustomerBuilder
		wantErr string
	}{
		{name: "没有入参", b: CustomerBuilder{}, args: args{name: "x1", s: []interface{}{}}, want: nil, wantErr: "配置：x1值不能为空"},
		{name: "入参是string", b: CustomerBuilder{}, args: args{name: "x1", s: []interface{}{"1111"}}, want: CustomerBuilder{"x1": json.RawMessage([]byte("1111"))}, wantErr: ""},
		{name: "入参是map", b: CustomerBuilder{}, args: args{name: "x1", s: []interface{}{map[string]string{"xx": "yy"}}}, want: CustomerBuilder{"x1": map[string]string{"xx": "yy"}}, wantErr: ""},
		{name: "入参是Ptr", b: CustomerBuilder{}, args: args{name: "x1", s: []interface{}{&args{}}}, want: CustomerBuilder{"x1": args{}}, wantErr: ""},
		{name: "入参是struct", b: CustomerBuilder{}, args: args{name: "x1", s: []interface{}{args{name: "ss"}}}, want: CustomerBuilder{"x1": args{name: "ss"}}, wantErr: ""},
		{name: "入参int", b: CustomerBuilder{}, args: args{name: "x1", s: []interface{}{32}}, want: nil, wantErr: "配置：x1值类型不支持"},
		{name: "入参float", b: CustomerBuilder{}, args: args{name: "x1", s: []interface{}{10.1}}, want: nil, wantErr: "配置：x1值类型不支持"},
		{name: "入参byte", b: CustomerBuilder{}, args: args{name: "x1", s: []interface{}{[]byte("10.1")}}, want: nil, wantErr: "配置：x1值类型不支持"},
		{name: "入参rune", b: CustomerBuilder{}, args: args{name: "x1", s: []interface{}{rune('1')}}, want: nil, wantErr: "配置：x1值类型不支持"},
	}
	for _, tt := range tests {
		func(name, wantErr string, b CustomerBuilder) {
			defer func() {
				e := types.GetString(recover())
				if e != "" {
					assert.Equal(t, true, strings.Contains(e, wantErr), name)
				} else {
					assert.Equal(t, wantErr, e, name)
				}
			}()
			got := b.Sub(tt.args.name, tt.args.s...)
			assert.Equal(t, tt.want, got, tt.name)
		}(tt.name, tt.wantErr, tt.b)
	}
}
