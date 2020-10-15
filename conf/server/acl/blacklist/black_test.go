package blacklist

import (
	"reflect"
	"testing"

	"github.com/micro-plat/hydra/conf"
	_ "github.com/micro-plat/hydra/registry/registry/localmemory"
)

//该方法依赖于conf-match函数的测试,match没有问题,该函数也就没有问题啦
func TestBlackList_IsDeny(t *testing.T) {
	type fields struct {
		Disable bool
		IPS     []string
		ipm     *conf.PathMatch
	}
	type args struct {
		ip string
	}
	tests := []struct {
		name   string
		fields fields
		args   args
		want   bool
	}{
		{name: "空对象匹配", args: args{ip: "192.168.5.101"}, fields: fields{IPS: []string{}, ipm: conf.NewPathMatch([]string{}...)}, want: false},
		{name: "路径匹配", args: args{ip: "/t1/t2/tt"}, fields: fields{IPS: []string{"/t1/**"}, ipm: conf.NewPathMatch([]string{"/t1/**"}...)}, want: true},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			w := &BlackList{
				Disable: tt.fields.Disable,
				IPS:     tt.fields.IPS,
				ipm:     tt.fields.ipm,
			}
			if got := w.IsDeny(tt.args.ip); got != tt.want {
				t.Errorf("BlackList.IsDeny() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestNew(t *testing.T) {
	type args struct {
		opts []Option
	}
	tests := []struct {
		name string
		args args
		want *BlackList
	}{
		{name: "初始化空对象", args: args{opts: nil}, want: &BlackList{IPS: []string{}, ipm: conf.NewPathMatch([]string{}...)}},
		{name: "初始化单ip对象", args: args{opts: []Option{WithIP("19.168.0.101")}}, want: &BlackList{IPS: []string{"19.168.0.101"}, ipm: conf.NewPathMatch([]string{"19.168.0.101"}...)}},
		{name: "初始化多ip对象", args: args{opts: []Option{WithIP("19.168.0.101"), WithIP("19.168.0.102")}}, want: &BlackList{IPS: []string{"19.168.0.101", "19.168.0.102"}, ipm: conf.NewPathMatch([]string{"19.168.0.101", "19.168.0.102"}...)}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := New(tt.args.opts...)
			if got.Disable != tt.want.Disable || !reflect.DeepEqual(got.IPS, tt.want.IPS) || !reflect.DeepEqual(*(got.ipm), *(tt.want.ipm)) {
				t.Errorf("New() = %v, want %v", got, tt.want)
			}
		})
	}
}
