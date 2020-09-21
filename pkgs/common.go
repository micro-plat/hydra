package pkgs

import (
	"sync"

	"github.com/micro-plat/lib4go/net"
)

var (
	mask    = ""
	localip = ""
)
var onceLock sync.Once

//LocalIP LocalIP
func LocalIP() string {
	onceLock.Do(func() {
		localip = net.GetLocalIPAddress(mask)
	})
	return localip
}

//WithIPMask WithIPMask
func WithIPMask(val string) {
	mask = val
}
