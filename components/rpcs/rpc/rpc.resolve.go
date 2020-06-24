package rpc

import (
	"fmt"
	"net"
	"strings"

	"github.com/asaskevich/govalidator"
)

const serviceRoot = "/%s/services/rpc/%s%s/providers"

//ResolvePath   解析注册中心地址
//order.request@merchant
func ResolvePath(address string, defPlatName string) (isip bool, service string, platName string, err error) {
	raddress := strings.TrimSpace(strings.TrimRight(address, "@"))
	addrs := strings.SplitN(raddress, "@", 2)
	if len(addrs) == 0 || addrs[0] == "" {
		return false, "", "", fmt.Errorf("服务名不能为空 %s", address)
	}
	if addrs[0] == "" {
		return false, "", "", fmt.Errorf("服务地址%s不能为空", address)
	}
	service = "/" + strings.Trim(strings.Replace(addrs[0], ".", "/", -1), "/")
	platName = defPlatName

	if len(addrs) > 1 && addrs[1] != "" {
		platName = addrs[1]
	}
	if isIPPort(platName) {
		return true, service, platName, nil
	}
	return false, service, platName, nil
}
func isIPPort(s string) bool {
	host, port, err := net.SplitHostPort(s)
	if err != nil {
		return false
	}
	return govalidator.IsIP(host) && govalidator.IsPort(port)
}
