// Package global linux version
package global

import (
	"errors"
	"os/exec"
	"strconv"
	"strings"
)

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

var errUnsupportedSystem = errors.New("Unsupported system")
var errRootPrivileges = errors.New("You must have root user privileges. Possibly using 'sudo' command should help")
