package global

import (
	xnet "net"

	"github.com/micro-plat/lib4go/security/md5"
)

var matchineCode = getMatchineCode()

//GetMatchineCode 获取机器码
func GetMatchineCode() string {
	return matchineCode
}

func getMatchineCode() string {
	if interfaces, err := xnet.Interfaces(); err == nil {
		mac := ""
		for _, inter := range interfaces {
			mac += inter.HardwareAddr.String()
		}
		return md5.Encrypt(mac)[:8]
	}
	return md5.Encrypt(LocalIP())[:8]
}
