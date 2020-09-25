package global

import (
	"testing"
)

func TestParseProto_Success(t *testing.T) {
	proto, raddr, err := ParseProto("zk://192.168.0.1")
	if err != nil {
		t.Error(err)
	}
	if proto != "zk" {
		t.Error("1")
	}
	if raddr != "192.168.0.1" {
		t.Error("2")
	}
}

func TestIsProto_Success(t *testing.T) {
	raddr, b := IsProto("zk://192.168.0.1", "zk")
	if b != true {
		t.Error("解析出错")
	}

	if raddr != "192.168.0.1" {
		t.Error("IP地址解析错误")
	}
}
