// Package compatible darwin (mac os x) version
package compatible

import (
	"os"
	"syscall"
)

//CheckPrivileges 检查是否有管理员权限
func CheckPrivileges() error {
	return nil
}

//CmdsRunNotifySignals  hydra/cmds/run/notify.Signal
var CmdsRunNotifySignals = []os.Signal{os.Interrupt, os.Kill, syscall.SIGTERM, syscall.SIGUSR2}

//CmdsUpdateProcessSignal CmdsUpdateProcessSignal
var CmdsUpdateProcessSignal = syscall.SIGUSR2

//AppClose AppClose
func AppClose() {
	parent := syscall.Getppid()
	syscall.Kill(parent, syscall.SIGUSR2)
}
