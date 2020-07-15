package global

import (
	"fmt"
	"strings"
)

const (
	ProtoZK   = "zk"
	ProtoRPC  = "rpc"
	ProtoHTTP = "http"
	ProtoLM   = "lm"
	ProtoFS   = "fs"
	ProtoLMQ  = "lmq"
)

//ParseProto 解析协议信息
func ParseProto(address string) (string, string, error) {
	addr := strings.Split(address, "://")
	if len(addr) != 2 {
		return "", "", fmt.Errorf("%s协议格式错误:proto://addr", addr)
	}
	proto := addr[0]
	raddr := addr[1]
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
