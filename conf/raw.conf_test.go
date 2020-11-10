package conf

import (
	"encoding/json"
	"reflect"
	"testing"

	"github.com/micro-plat/hydra/test/assert"
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
		{name: "nil数据初始化", args: args{data: nil, version: 0}, wantC: &RawConf{XMap: nil, version: 0, raw: []byte("null"), signature: md5.EncryptBytes([]byte("null"))}, wantErr: false},
		{name: "空数据初始化", args: args{data: dataN, version: 0}, wantC: &RawConf{XMap: dataN, version: 0, raw: []byte("{}"), signature: md5.EncryptBytes([]byte("{}"))}, wantErr: false},
		{name: "对象数据初始化", args: args{data: dataV, version: 0}, wantC: &RawConf{XMap: dataV, version: 0, raw: dataB, signature: md5.EncryptBytes(dataB)}, wantErr: false},
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
		{name: "反序列化失败初始化", args: args{message: []byte("{sdfsdfsdf}"), version: 0}, wantC: &RawConf{XMap: nil, version: 0, raw: []byte("{sdfsdfsdf}"), signature: md5.EncryptBytes([]byte("{sdfsdfsdf}"))}, wantErr: true},
		{name: "正常初始化", args: args{message: []byte(`{"sss":"11"}`), version: 0}, wantC: &RawConf{XMap: map[string]interface{}{"sss": "11"}, version: 0, raw: []byte(`{"sss":"11"}`), signature: md5.EncryptBytes([]byte(`{"sss":"11"}`))}, wantErr: false},
		{name: "无data,正常初始化", args: args{message: []byte(`test1`), version: 0}, wantC: &RawConf{XMap: nil, version: 0, raw: []byte("test1"), signature: md5.EncryptBytes([]byte(`test1`))}, wantErr: true},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotC, err := NewByJSON(tt.args.message, tt.args.version)
			if (err != nil) != tt.wantErr {
				t.Errorf("NewByJSON() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			if err == nil && !reflect.DeepEqual(gotC, tt.wantC) {
				t.Errorf("NewByJSON() = %v, want %v", gotC, tt.wantC)
			}
		})
	}
}

func TestRawConf_GetSignature(t *testing.T) {
	obj, _ := NewByJSON([]byte(`{"sss":"11"}`), 0)
	tests := []struct {
		name   string
		fields *RawConf
		want   string
	}{
		{name: "正常的数据md5", fields: obj, want: md5.EncryptBytes([]byte(`{"sss":"11"}`))},
	}
	for _, tt := range tests {
		got := tt.fields.GetSignature()
		assert.Equal(t, tt.want, got, tt.name)
	}
}

func TestRawConf_GetVersion(t *testing.T) {
	obj, _ := NewByJSON([]byte(`{"sss":"11"}`), 1000)
	obj1, _ := NewByJSON([]byte(`{"sss":"11"}`), 122)
	tests := []struct {
		name   string
		fields *RawConf
		want   int32
	}{
		{name: "获取对象的版本号1", fields: obj, want: 1000},
		{name: "获取对象的版本号2", fields: obj1, want: 122},
	}
	for _, tt := range tests {
		got := tt.fields.GetVersion()
		assert.Equal(t, tt.want, got, tt.name)
	}
}

func TestRawConf_GetRaw(t *testing.T) {
	obj, _ := NewByJSON([]byte(`{"sss":"11"}`), 1)
	obj1, _ := NewByJSON([]byte(`{"sss":"12121122"}`), 1)
	tests := []struct {
		name   string
		fields *RawConf
		want   []byte
	}{
		{name: "获取对象的byte1", fields: obj, want: []byte(`{"sss":"11"}`)},
		{name: "获取对象的byte2", fields: obj1, want: []byte(`{"sss":"12121122"}`)},
	}
	for _, tt := range tests {
		got := tt.fields.GetRaw()
		assert.Equal(t, tt.want, got, tt.name)
	}
}

func TestRawConf_GetArray(t *testing.T) {
	type fields struct {
		raw       json.RawMessage
		version   int32
		signature string
		data      map[string]interface{}
	}
	type args struct {
		key string
		def []interface{}
	}
	tests := []struct {
		name   string
		fields fields
		args   args
		wantR  []interface{}
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			j := &RawConf{
				raw:       tt.fields.raw,
				version:   tt.fields.version,
				signature: tt.fields.signature,
				XMap:      tt.fields.data,
			}
			if gotR := j.GetArray(tt.args.key, tt.args.def...); !reflect.DeepEqual(gotR, tt.wantR) {
				t.Errorf("RawConf.GetArray() = %v, want %v", gotR, tt.wantR)
			}
		})
	}
}

func TestRawConf_GetJSON(t *testing.T) {
	type fields struct {
		raw       json.RawMessage
		version   int32
		signature string
		data      map[string]interface{}
	}
	type args struct {
		section string
	}
	tests := []struct {
		name        string
		fields      fields
		args        args
		wantR       []byte
		wantVersion int32
		wantErr     bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			j := &RawConf{
				raw:       tt.fields.raw,
				version:   tt.fields.version,
				signature: tt.fields.signature,
				XMap:      tt.fields.data,
			}
			gotR, err := j.GetJSON(tt.args.section)
			if (err != nil) != tt.wantErr {
				t.Errorf("RawConf.GetJSON() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(gotR, tt.wantR) {
				t.Errorf("RawConf.GetJSON() gotR = %v, want %v", gotR, tt.wantR)
			}
			// if gotVersion != tt.wantVersion {
			// 	t.Errorf("RawConf.GetJSON() gotVersion = %v, want %v", gotVersion, tt.wantVersion)
			// }
		})
	}
}

func TestRawConf_HasSection(t *testing.T) {
	type fields struct {
		raw       json.RawMessage
		version   int32
		signature string
		data      map[string]interface{}
	}
	type args struct {
		section string
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
			j := &RawConf{
				raw:       tt.fields.raw,
				version:   tt.fields.version,
				signature: tt.fields.signature,
				XMap:      tt.fields.data,
			}
			if got := j.Has(tt.args.section); got != tt.want {
				t.Errorf("RawConf.HasSection() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestRawConf_GetSection(t *testing.T) {
	type fields struct {
		raw       json.RawMessage
		version   int32
		signature string
		data      map[string]interface{}
	}
	type args struct {
		section string
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		wantC   *RawConf
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			j := &RawConf{
				raw:       tt.fields.raw,
				version:   tt.fields.version,
				signature: tt.fields.signature,
				XMap:      tt.fields.data,
			}
			gotC, err := j.GetXMap(tt.args.section)
			if (err != nil) != tt.wantErr {
				t.Errorf("RawConf.GetSection() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(gotC, tt.wantC) {
				t.Errorf("RawConf.GetSection() = %v, want %v", gotC, tt.wantC)
			}
		})
	}
}
