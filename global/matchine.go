package global

import (
	"net"
	"sync"

	"fmt"
	"strings"

	"github.com/asaskevich/govalidator"
	"github.com/micro-plat/hydra/global/compatible"
	"github.com/micro-plat/lib4go/security/md5"

	xnet "github.com/micro-plat/lib4go/net"
)

var (
	localip = ""
)
var onceLock sync.Once

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
		if err := compatible.CheckPrivileges(); err != nil {
			return "", "", err
		}
	}
	if host == "" {
		return "0.0.0.0", port, nil
	}
	return host, port, nil
}

//LocalIP LocalIP
func LocalIP() string {
	onceLock.Do(func() {
		localip = xnet.GetLocalIPAddress(Def.IPMask)
	})
	return localip
}

var matchineCode = getMatchineCode()

//GetMatchineCode 获取机器码
func GetMatchineCode() string {
	return matchineCode
}

func getMatchineCode() string {
	if interfaces, err := net.Interfaces(); err == nil {
		mac := ""
		for _, inter := range interfaces {
			mac += inter.HardwareAddr.String()
		}
		return md5.Encrypt(mac)[:8]
	}
	return md5.Encrypt(LocalIP())[:8]
}
