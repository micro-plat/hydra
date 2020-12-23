package global

import (
	"fmt"
	"strings"
)

const (
	ProtoZK      = "zk"
	ProtoRPC     = "rpc"
	ProtoHTTP    = "http"
	ProtoLM      = "lm"
	ProtoFS      = "fs"
	ProtoLMQ     = "lmq"
	ProtoREDIS   = "redis"
	ProtoInvoker = "ivk"
)

//ParseProto 解析协议信息
func ParseProto(address string) (string, string, error) {
	address = strings.Trim(address, " ")

	addr := strings.Split(address, "://")
	if len(addr) != 2 {
		return "", "", fmt.Errorf("%s协议格式错误,正确格式(proto://addr)", addr)
	}
	proto := addr[0]
	if proto == "" {
		return "", "", fmt.Errorf("%s缺少协议proto,正确格式(proto://addr)", address)
	}
	raddr := addr[1]
	if raddr == "" {
		return "", "", fmt.Errorf("%s缺少地址addr,正确格式(proto://addr)", address)
	}
	return proto, raddr, nil
}

//IsProto 是否是指定的协议
func IsProto(addr string, proto string) (string, bool) {
	p, addrs, _ := ParseProto(addr)
	return addrs, p == proto
}

//IsLocal 是否是本地服务
func IsLocal(proto string) bool {
	return strings.EqualFold(proto, ProtoLM) || strings.EqualFold(proto, ProtoLMQ)
}
