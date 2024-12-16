package conf

import (
	"encoding/json"
	"reflect"
	"testing"

	"github.com/micro-plat/lib4go/assert"
	"github.com/micro-plat/lib4go/security/md5"
)

func TestNewByMap(t *testing.T) {
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
		{name: "1. rawconf-NewByMap-nil数据初始化", args: args{data: nil, version: 0}, wantC: &RawConf{XMap: nil, version: 0, raw: []byte("null"), signature: md5.EncryptBytes([]byte("null"))}, wantErr: false},
		{name: "2. rawconf-NewByMap-空数据初始化", args: args{data: dataN, version: 0}, wantC: &RawConf{XMap: dataN, version: 0, raw: []byte("{}"), signature: md5.EncryptBytes([]byte("{}"))}, wantErr: false},
		{name: "3. rawconf-NewByMap-对象数据初始化", args: args{data: dataV, version: 0}, wantC: &RawConf{XMap: dataV, version: 0, raw: dataB, signature: md5.EncryptBytes(dataB)}, wantErr: false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotC, err := NewByMap(tt.args.data, tt.args.version)
			if (err != nil) != tt.wantErr {
				t.Errorf("NewByMap() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(gotC, tt.wantC) {
				t.Errorf("NewByMap() = %v, want %v", gotC, tt.wantC)
			}
		})
	}
}

func TestNewByJSON(t *testing.T) {
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
		{name: "1. rawconf-NewByJSON-反序列化失败初始化", args: args{message: []byte("{sdfsdfsdf}"), version: 0}, wantC: &RawConf{XMap: nil, version: 0, raw: []byte("{sdfsdfsdf}"), signature: md5.EncryptBytes([]byte("{sdfsdfsdf}"))}, wantErr: true},
		{name: "2. rawconf-NewByJSON-正常初始化", args: args{message: []byte(`{"sss":"11"}`), version: 0}, wantC: &RawConf{XMap: map[string]interface{}{"sss": "11"}, version: 0, raw: []byte(`{"sss":"11"}`), signature: md5.EncryptBytes([]byte(`{"sss":"11"}`))}, wantErr: false},
		{name: "3. rawconf-NewByJSON-无data,正常初始化", args: args{message: []byte(`test1`), version: 0}, wantC: &RawConf{XMap: nil, version: 0, raw: []byte("test1"), signature: md5.EncryptBytes([]byte(`test1`))}, wantErr: false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotC, err := NewByText(tt.args.message, tt.args.version)
			if (err != nil) != tt.wantErr {
				t.Errorf("NewByText() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			if err == nil && !reflect.DeepEqual(gotC, tt.wantC) {
				t.Errorf("NewByText() = %v, want %v", gotC, tt.wantC)
			}
		})
	}
}

func TestRawConf_GetSignature(t *testing.T) {
	obj, _ := NewByText([]byte(`{"sss":"11"}`), 0)
	tests := []struct {
		name   string
		fields *RawConf
		isbool bool
		want   string
	}{
		{name: "1. rawconf-GetSignature-正常的数据md5", isbool: true, fields: obj, want: md5.EncryptBytes([]byte(`{"sss":"11"}`))},
		{name: "2. rawconf-GetSignature-非md5", isbool: false, fields: obj, want: "2222"},
	}
	for _, tt := range tests {
		got := tt.fields.GetSignature()
		if tt.isbool {
			assert.Equal(t, tt.want, got, tt.name)
		} else {
			assert.NotEqual(t, tt.want, got, tt.name)
		}
	}
}

func TestRawConf_GetVersion(t *testing.T) {
	obj, _ := NewByText([]byte(`{"sss":"11"}`), 1000)
	obj1, _ := NewByText([]byte(`{"sss":"11"}`), 122)
	tests := []struct {
		name   string
		fields *RawConf
		want   int32
	}{
		{name: "1. rawconf-GetVersion-获取对象的版本号1", fields: obj, want: 1000},
		{name: "2. rawconf-GetVersion-获取对象的版本号2", fields: obj1, want: 122},
	}
	for _, tt := range tests {
		got := tt.fields.GetVersion()
		assert.Equal(t, tt.want, got, tt.name)
	}
}

func TestRawConf_GetRaw(t *testing.T) {
	obj, _ := NewByText([]byte(`{"sss":"11"}`), 1)
	obj1, _ := NewByText([]byte(`{"sss":"12121122"}`), 1)
	tests := []struct {
		name   string
		fields *RawConf
		want   []byte
	}{
		{name: "1. rawconf-GetRaw-获取对象的byte1", fields: obj, want: []byte(`{"sss":"11"}`)},
		{name: "2. rawconf-GetRaw-获取对象的byte2", fields: obj1, want: []byte(`{"sss":"12121122"}`)},
	}
	for _, tt := range tests {
		got := tt.fields.GetRaw()
		assert.Equal(t, tt.want, got, tt.name)
	}
}
