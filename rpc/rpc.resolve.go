package rpc

import (
	"fmt"
	"net"
	"strings"

	"github.com/asaskevich/govalidator"
)

//ResolvePath   解析注册中心地址
//domain:hydra,server:merchant_cron
//order.request@merchant_api.hydra 解析为:service: /order/request,server:merchant_api,domain:hydra
//order.request 解析为 service: /order/request,server:merchant_cron,domain:hydra
//order.request@merchant_rpc 解析为 service: /order/request,server:merchant_rpc,domain:hydra
func ResolvePath(address string, d string, s string) (isip bool, service string, domain string, server string, err error) {
	raddress := strings.TrimRight(address, "@")
	addrs := strings.SplitN(raddress, "@", 2)
	if len(addrs) == 1 {
		if addrs[0] == "" {
			return false, "", "", "", fmt.Errorf("服务地址%s不能为空", address)
		}
		service = "/" + strings.Trim(strings.Replace(raddress, ".", "/", -1), "/")
		domain = d
		server = s
		return
	}
	if addrs[0] == "" {
		return false, "", "", "", fmt.Errorf("%s错误，服务名不能为空", address)
	}
	if addrs[1] == "" {
		return false, "", "", "", fmt.Errorf("%s错误，服务名，域不能为空", address)
	}
	if isIPPort(addrs[1]) {
		return true, addrs[0], d, addrs[1], nil
	}

	service = "/" + strings.Trim(strings.Replace(addrs[0], ".", "/", -1), "/") //处理服务名中的特殊字符
	raddr := strings.Split(strings.TrimRight(addrs[1], "."), ".")
	if len(raddr) >= 2 && raddr[0] != "" && raddr[1] != "" {
		domain = raddr[len(raddr)-1]
		server = strings.Join(raddr[0:len(raddr)-1], ".")
		return
	}
	if len(raddr) == 1 {
		if raddr[0] == "" {
			return false, "", "", "", fmt.Errorf("%s错误，服务器名称不能为空", address)
		}
		domain = d
		server = raddr[0]
		return
	}
	if raddr[0] == "" && raddr[1] == "" {
		return false, "", "", "", fmt.Errorf(`%s错误,未指定服务器名称和域名称`, addrs[1])
	}
	domain = raddr[1]
	server = s

	return
}
func isIPPort(s string) bool {
	host, port, err := net.SplitHostPort(s)
	if err != nil {
		return false
	}
	return govalidator.IsIP(host) && govalidator.IsPort(port)
}
