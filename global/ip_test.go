package global

import (
	"fmt"
	"testing"
)

func TestGetHostPort_Success(t *testing.T) {
	host, port, err := GetHostPort("192.168.0.1:9090")
	if err != nil {
		t.Error(err)
	}
	if host != "192.168.0.1" {
		t.Error("1")
	}
	if port != "9090" {
		t.Error("2")
	}

}

func TestGetHostPort_NoneIP(t *testing.T) {
	host, port, err := GetHostPort(":9090")
	fmt.Println(host, port, err)
	if err != nil {
		t.Error(err)
	}
	if host != "0.0.0.0" {
		t.Error("1")
	}
	if port != "9090" {
		t.Error("2")
	}

}

func TestGetHostPort_ErrorPort(t *testing.T) {
	host, port, err := GetHostPort(":aaa")
	fmt.Println(host, port, err)
	if err == nil {
		t.Error("端口测试不通过")
	}
}
