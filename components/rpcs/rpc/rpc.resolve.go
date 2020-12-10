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
	if addrs[0] == "" {
		return false, "", "", fmt.Errorf("服务地址不能为空:%s", address)
	}
	service = "/" + strings.Trim(strings.Replace(addrs[0], ".", "/", -1), "/")
	platName = defPlatName

	if len(addrs) > 1 && addrs[1] != "" {
		platName = addrs[1]
	}

	if newAddr, b := isIPPort(platName); b {
		return true, service, newAddr, nil
	}
	return false, service, platName, nil
}
func isIPPort(s string) (string, bool) {
	if strings.Contains(s, "://") {
		parties := strings.Split(s, "://")
		if len(parties) != 2 {
			return "", false
		}
		s = parties[1]
	}

	host, port, err := net.SplitHostPort(s)
	if err != nil {
		return "", false
	}
	return fmt.Sprintf("tcp://%s:%s", host, port), govalidator.IsIP(host) && govalidator.IsPort(port)
}
