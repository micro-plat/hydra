package global

import (
	"fmt"
	"net"
	"strings"

	"github.com/asaskevich/govalidator"
)

//GetHostPort 获取服务器名及端口
func GetHostPort(addr string) (host string, port string, err error) {
	if !strings.Contains(addr, ":") {
		port = addr
	} else {
		host, port, err = net.SplitHostPort(addr)
	}
	if err != nil {
		return "", "", err
	}
	if !govalidator.IsPort(port) {
		return "", "", fmt.Errorf("端口不合法 %s", port)
	}
	if port == "80" {
		if err := CheckPrivileges(); err != nil {
			return "", "", err
		}
	}
	if host == "" {
		return "0.0.0.0", port, nil
	}
	return host, port, nil
}
