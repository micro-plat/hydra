package conf

import (
	"strings"
	"testing"

	"github.com/micro-plat/lib4go/assert"
)

func BenchmarkEncrypt(b *testing.B) {
	b.ResetTimer()
	var input = []byte("taosytaosytaosytaosytaosytaosytaosy")
	for i := 0; i < b.N; i++ {
		Encrypt(input)
	}
}

func BenchmarkDecrypt(b *testing.B) {
	var input = []byte("encrypt:cbc/pkcs5:47b17dd320c67986a839e86c4da057a95ff57a008f168817daf12500e475dfccf032844585f9723c")
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		Decrypt(input)
	}
}

func Test_encrypt(t *testing.T) {
	tests := []struct {
		name  string
		input []byte
	}{
		{name: "1. conf-encrypt-空数据加密", input: []byte{}},
		{name: "2. conf-encrypt-数据加密", input: []byte("taosytaosytaosytaosytaosytaosytaosy")},
	}
	for _, tt := range tests {
		got := Encrypt(tt.input)
		list := strings.Split(got, ":")
		assert.Equal(t, len(list), 3, tt.name+",len")
		if len(list) >= 2 {
			assert.Equal(t, list[0], hd, tt.name+".hd")
			assert.Equal(t, list[1], mode, tt.name+",mode")

		}
	}
}

func Test_decrypt(t *testing.T) {
	input := []byte{}
	nildata := Encrypt(input)
	input1 := []byte("encryptapsytsetetapsytsetetapsytsetetapsytsete")
	data1 := Encrypt(input1)

	tests := []struct {
		name    string
		data    []byte
		want    []byte
		wantErr bool
	}{
		{name: "1. conf-decrypt-空数据解密", data: []byte(nildata), want: []byte{}, wantErr: false},
		{name: "2. conf-decrypt-小于加密前后缀的长度的数据解密", data: []byte("nildata"), want: []byte("nildata"), wantErr: false},
		{name: "3. conf-decrypt-不是由hd开头数据解密", data: []byte("nildatanildatanildatanildata"), want: []byte("nildatanildatanildatanildata"), wantErr: false},
		{name: "4. conf-decrypt-错误数据解密", data: []byte("encryptnildatanildatanildatanildata"), want: nil, wantErr: true},
		{name: "5. conf-decrypt-正确数据数据解密", data: []byte(data1), want: []byte("encryptapsytsetetapsytsetetapsytsetetapsytsete"), wantErr: false},
	}
	for _, tt := range tests {
		got, err := Decrypt(tt.data)
		assert.Equal(t, tt.wantErr, (err != nil), tt.name+".err")
		assert.Equal(t, tt.want, got, tt.name+",res")
	}
}
