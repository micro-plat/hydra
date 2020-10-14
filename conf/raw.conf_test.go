package conf

import (
	"encoding/json"
	"reflect"
	"testing"

	"github.com/micro-plat/lib4go/security/md5"
)

func TestNewRawConfByMap(t *testing.T) {
	dataN := map[string]interface{}{}
	dataV := map[string]interface{}{"t": "x"}
	dataB, _ := json.Marshal(dataV)
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
		{name: "nil数据初始化", args: args{data: nil, version: 0}, wantC: &RawConf{data: nil, version: 0, raw: []byte("null"), signature: md5.EncryptBytes([]byte("null"))}, wantErr: false},
		{name: "空数据初始化", args: args{data: dataN, version: 0}, wantC: &RawConf{data: dataN, version: 0, raw: []byte("{}"), signature: md5.EncryptBytes([]byte("{}"))}, wantErr: false},
		{name: "对象数据初始化", args: args{data: dataV, version: 0}, wantC: &RawConf{data: dataV, version: 0, raw: dataB, signature: md5.EncryptBytes(dataB)}, wantErr: false},
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

func TestNewRawConfByJson(t *testing.T) {
	type args struct {
		message []byte
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
			gotC, err := NewRawConfByJson(tt.args.message, tt.args.version)
			if (err != nil) != tt.wantErr {
				t.Errorf("NewRawConfByJson() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(gotC, tt.wantC) {
				t.Errorf("NewRawConfByJson() = %v, want %v", gotC, tt.wantC)
			}
		})
	}
}
