package pkgs

import (
	"github.com/micro-plat/lib4go/net"
)

var (
	mask = ""
)

//LocalIP LocalIP
func LocalIP() string {
	return net.GetLocalIPAddress(mask)
}

//WithIPMask WithIPMask
func WithIPMask(val string) {
	mask = val
}
