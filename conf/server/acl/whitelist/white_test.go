package whitelist

import (
	"reflect"
	"testing"
)

//百名单匹配暂时不要测试  匹配方案没有确定
func TestWhiteList_IsAllow(t *testing.T) {
	type fields struct {
		Disable bool
		IPS     []*IPList
	}
	type args struct {
		path string
		ip   string
	}
	tests := []struct {
		name   string
		fields fields
		args   args
		want   bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			w := &WhiteList{
				Disable: tt.fields.Disable,
				IPS:     tt.fields.IPS,
			}
			if got := w.IsAllow(tt.args.path, tt.args.ip); got != tt.want {
				t.Errorf("WhiteList.IsAllow() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestNew(t *testing.T) {
	tests := []struct {
		name string
		args []Option
		want *WhiteList
	}{
		{name: "初始化默认百名单配置", args: []Option{}, want: &WhiteList{IPS: make([]*IPList, 0, 1), Disable: false}},
		{name: "初始化Disable百名单配置", args: []Option{WithDisable()}, want: &WhiteList{IPS: make([]*IPList, 0, 1), Disable: true}},
		{name: "初始化Enable百名单配置", args: []Option{WithEnable()}, want: &WhiteList{IPS: make([]*IPList, 0, 1), Disable: false}},
		{name: "初始化默认百名单配置", args: []Option{}, want: &WhiteList{IPS: make([]*IPList, 0, 1)}},
		{name: "初始化自定义ip百名单配置", args: []Option{WithIPList(NewIPList("/t1/t2/*", []string{"192.168.0.101"}...))}, want: &WhiteList{IPS: []*IPList{&IPList{Requests: []string{"/t1/t2/*"}, IPS: []string{"192.168.0.101"}}}, Disable: false}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := New(tt.args...)
			if tt.name != "初始化自定义ip百名单配置" && !reflect.DeepEqual(got, tt.want) {
				t.Errorf("New() = %v, want %v", got, tt.want)
			}

			if tt.name == "初始化自定义ip百名单配置" {
				for i, item := range got.IPS {
					if !reflect.DeepEqual(item.IPS, tt.want.IPS[i].IPS) ||
						!reflect.DeepEqual(item.Requests, tt.want.IPS[i].Requests) {
						t.Errorf("New1() = %v, want %v", item, tt.want.IPS[i])
					}
				}
			}
		})
	}
}
