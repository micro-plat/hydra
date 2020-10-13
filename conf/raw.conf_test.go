package conf

import (
	"reflect"
	"testing"
)

func TestNewRawConfByMap(t *testing.T) {
	type args struct {
		data    map[string]interface{}
		version int32
	}
	tests := []struct {
		name    string
		args    args
		wantC   *RawConf
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotC, err := NewRawConfByMap(tt.args.data, tt.args.version)
			if (err != nil) != tt.wantErr {
				t.Errorf("NewRawConfByMap() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(gotC, tt.wantC) {
				t.Errorf("NewRawConfByMap() = %v, want %v", gotC, tt.wantC)
			}
		})
	}
}
