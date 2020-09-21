package server

import (
	"bytes"
	"strings"
	"testing"
)

func BenchmarkEncrypt(b *testing.B) {
	b.ResetTimer()
	var input = []byte("taosytaosytaosytaosytaosytaosytaosy")
	for i := 0; i < b.N; i++ {
		encrypt(input)
	}
}

func BenchmarkDecrypt(b *testing.B) {
	// var inputOr = []byte("taosytaosytaosytaosytaosytaosytaosy")
	var input = []byte("encrypt:cbc/pkcs5:47b17dd320c67986a839e86c4da057a95ff57a008f168817daf12500e475dfccf032844585f9723c")
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		decrypt(input)
		// if err != nil {
		// 	b.Errorf("性能测试解密错误,err:%+v", err)
		// 	return
		// }
		// if bytes.Compare(inputOr, res) != 0 {
		// 	b.Error("性能测试解密错误")
		// 	return
		// }
	}
}

func TestEncrypt(t *testing.T) {
	var input = []byte("taosytaosytaosytaosytaosytaosytaosy")
	data := encrypt(input)
	list := strings.Split(data, ":")
	if len(list) <= 2 {
		t.Error("加密结果错误")
		return
	}
	if hd != list[0] || mode != list[1] {
		t.Error("加密结果错误1")
		return
	}

	//空数据加密
	input = []byte{}
	data = encrypt(input)
	list = strings.Split(data, ":")
	if len(list) <= 2 {
		t.Error("加密结果错误")
		return
	}
	if hd != list[0] || mode != list[1] {
		t.Error("加密结果错误1")
		return
	}
	return
}

func TestDecrypt(t *testing.T) {
	var input = []byte("taosy")
	data := encrypt(input)
	res, err := decrypt([]byte(data))
	if err != nil {
		t.Error("解密结果错误")
		return
	}

	if bytes.Compare(input, res) != 0 {
		t.Error("解密结果错误1")
		return
	}

	//空数据加密
	input = []byte{}
	data = encrypt(input)
	res, err = decrypt([]byte(data))
	if err != nil {
		t.Error("解密结果错误2")
		return
	}

	if bytes.Compare(input, res) != 0 {
		t.Error("解密结果错误3")
		return
	}
}
