package cache

import (
	"reflect"
	"testing"
)

func Test_getNames(t *testing.T) {
	type args struct {
		address string
	}
	tests := []struct {
		name      string
		args      args
		wantProto string
		wantRaddr []string
		wantErr   bool
	}{
		{name: "1", args: args{address: ""}, wantErr: true},
		{name: "2", args: args{address: "address"}, wantErr: true},
		{name: "3", args: args{address: "memcached://192.168.0.1:11211"}, wantProto: "memcached", wantRaddr: []string{"192.168.0.1:11211"}, wantErr: false},
		{name: "4", args: args{address: "redis://192.168.0.1:11211"}, wantProto: "redis", wantRaddr: []string{"192.168.0.1:11211"}, wantErr: false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotProto, gotRaddr, err := getNames(tt.args.address)
			if (err != nil) != tt.wantErr {
				t.Errorf("getNames() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if gotProto != tt.wantProto {
				t.Errorf("getNames() gotProto = %v, want %v", gotProto, tt.wantProto)
			}
			if !reflect.DeepEqual(gotRaddr, tt.wantRaddr) {
				t.Errorf("getNames() gotRaddr = %v, want %v", gotRaddr, tt.wantRaddr)
			}
		})
	}
}
