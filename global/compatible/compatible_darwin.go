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
	parent := syscall.Getpid()
	syscall.Kill(parent, syscall.SIGTERM)
}

const (
	SUCCESS = "\033[32m\t\t\t\t\t[OK]\033[0m"     // Show colored "OK"
	FAILED  = "\033[31m\t\t\t\t\t[FAILED]\033[0m" // Show colored "FAILED"
)
