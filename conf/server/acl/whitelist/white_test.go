package whitelist

import (
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
