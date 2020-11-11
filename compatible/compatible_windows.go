// Package compatible windows version
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
var CmdsRunNotifySignals = []os.Signal{os.Interrupt, os.Kill, syscall.SIGTERM, syscall.SIGINT}

//CmdsUpdateProcessSignal CmdsUpdateProcessSignal
var CmdsUpdateProcessSignal = syscall.SIGINT

//AppClose AppClose
func AppClose() {
	syscall.Exit(0)
}
