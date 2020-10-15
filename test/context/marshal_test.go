package context

import (
	"reflect"
	"testing"

	"github.com/micro-plat/hydra/context"
)

func TestUnmarshalXML(t *testing.T) {
	type args struct {
		s string
	}
	tests := []struct {
		name    string
		args    args
		want    map[string]string
		wantErr bool
	}{
		{name: "转换正确的xml", args: args{s: `<xml><key>value</key></xml>`}, want: map[string]string{"key": "value"}, wantErr: false},
		{name: "转换错误的xml", args: args{s: `<xml><key>value</ky></xml>`}, want: map[string]string{}, wantErr: true},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := context.UnmarshalXML(tt.args.s)
			if (err != nil) != tt.wantErr {
				t.Errorf("UnmarshalXML() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("UnmarshalXML() = %v, want %v", got, tt.want)
			}
		})
	}
}
