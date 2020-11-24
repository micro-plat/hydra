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
		name   string
		args   []interface{}
		repeat []interface{}
		want   CustomerBuilder
	}{
		{name: "1. 没有入参,初始化空对象", args: []interface{}{}, want: CustomerBuilder{"main": make(map[string]interface{})}},
		{name: "2. 单入参,初始化对象", args: []interface{}{"参数1"}, want: CustomerBuilder{"main": "参数1"}},
		{name: "3. 多入参,初始化对象", args: []interface{}{"参数2", "参数1"}, want: CustomerBuilder{"main": "参数2"}},
		{name: "4. 重复入参,初始化对象", args: []interface{}{"参数2"}, repeat: []interface{}{"参数1"}, want: CustomerBuilder{"main": "参数1"}},
	}
	for _, tt := range tests {
		got := newCustomerBuilder(tt.args...)
		if tt.repeat != nil {
			got = newCustomerBuilder(tt.repeat...)
		}
		assert.Equal(t, tt.want, got, tt.name)
	}
}

func TestCustomerBuilder_Sub(t *testing.T) {
	type args struct {
		name   string
		s      []interface{}
		repeat []interface{}
	}
	tests := []struct {
		name    string
		b       CustomerBuilder
		args    args
		want    CustomerBuilder
		wantErr string
	}{
		{name: "1. 无入参,构建自定义子节点", b: CustomerBuilder{}, args: args{name: "x1", s: []interface{}{}}, want: nil, wantErr: "配置：x1值不能为空"},
		{name: "2. 入参-string,构建自定义子节点", b: CustomerBuilder{}, args: args{name: "x1", s: []interface{}{"1111"}}, want: CustomerBuilder{"x1": json.RawMessage([]byte("1111"))}, wantErr: ""},
		{name: "3. 入参-map,构建自定义子节点", b: CustomerBuilder{}, args: args{name: "x1", s: []interface{}{map[string]string{"xx": "yy"}}}, want: CustomerBuilder{"x1": map[string]string{"xx": "yy"}}, wantErr: ""},
		{name: "4. 入参-Ptr,构建自定义子节点", b: CustomerBuilder{}, args: args{name: "x1", s: []interface{}{&args{}}}, want: CustomerBuilder{"x1": args{}}, wantErr: ""},
		{name: "5. 入参-struct,构建自定义子节点", b: CustomerBuilder{}, args: args{name: "x1", s: []interface{}{args{name: "ss"}}}, want: CustomerBuilder{"x1": args{name: "ss"}}, wantErr: ""},
		{name: "6. 入参-int,构建自定义子节点", b: CustomerBuilder{}, args: args{name: "x1", s: []interface{}{32}}, want: nil, wantErr: "配置：x1值类型不支持"},
		{name: "7. 入参-float,构建自定义子节点", b: CustomerBuilder{}, args: args{name: "x1", s: []interface{}{10.1}}, want: nil, wantErr: "配置：x1值类型不支持"},
		{name: "8. 入参-byte,构建自定义子节点", b: CustomerBuilder{}, args: args{name: "x1", s: []interface{}{[]byte("10.1")}}, want: nil, wantErr: "配置：x1值类型不支持"},
		{name: "9. 入参-rune,构建自定义子节点", b: CustomerBuilder{}, args: args{name: "x1", s: []interface{}{rune('1')}}, want: nil, wantErr: "配置：x1值类型不支持"},
		{name: "10. 重复入参-string,构建自定义子节点", b: CustomerBuilder{}, args: args{name: "x1", s: []interface{}{"1111"}, repeat: []interface{}{"2222"}}, want: CustomerBuilder{"x1": json.RawMessage([]byte("2222"))}, wantErr: ""},
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
			if tt.args.repeat != nil {
				got = b.Sub(tt.args.name, tt.args.repeat...)
			}
			assert.Equal(t, tt.want, got, tt.name)
		}(tt.name, tt.wantErr, tt.b)
	}
}
