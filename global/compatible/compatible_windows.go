// Package compatible windows version
package compatible

import (
	"fmt"
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
	pid := syscall.Getpid()

	dll, err := syscall.LoadDLL("kernel32.dll")
	if err != nil {
		fmt.Println("LoadDLL:kernel32.dll", err)
		return
	}

	p, err := dll.FindProc("GenerateConsoleCtrlEvent")
	if err != nil {
		fmt.Println("FindProc:", err)
		return
	}

	r, _, err := p.Call(syscall.CTRL_BREAK_EVENT, uintptr(pid))
	if r == 0 {
		fmt.Println("GenerateConsoleCtrlEvent:", err)
		return
	}
}

const (
	SUCCESS = "\t\t\t\t\t[OK]"     // Show colored "OK"
	FAILED  = "\t\t\t\t\t[FAILED]" // Show colored "FAILED"
)
