// Package compatible linux version
package compatible

import (
	"errors"
	"os"
	"os/exec"
	"strconv"
	"strings"
	"syscall"
)

var errUnsupportedSystem = errors.New("Unsupported system")
var errRootPrivileges = errors.New("You must have root user privileges. Possibly using 'sudo' command should help")

//CheckPrivileges 检查是否有管理员权限
func CheckPrivileges() error {
	if output, err := exec.Command("id", "-g").Output(); err == nil {
		if gid, parseErr := strconv.ParseUint(strings.TrimSpace(string(output)), 10, 32); parseErr == nil {
			if gid == 0 {
				return nil
			}
			return errRootPrivileges
		}
	}
	return errUnsupportedSystem
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
