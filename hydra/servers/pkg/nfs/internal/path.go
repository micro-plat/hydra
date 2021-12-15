package internal

import (
	"bytes"
	"fmt"
	"os"
	"os/exec"
	"strings"
)

func interactiveToexec(commandName string, params []string) (string, error) {
	cmd := exec.Command(commandName, params...)
	buf, err := cmd.Output()
	w := bytes.NewBuffer(nil)
	cmd.Stderr = w
	if err != nil {
		return "", fmt.Errorf("文档转换失败：%w %s %s", err, commandName, strings.Join(params, " "))
	}
	return string(buf), nil

}

func pathExists(path string) (bool, error) {
	_, err := os.Stat(path)
	if err == nil {
		return true, nil
	}
	if os.IsNotExist(err) {
		return false, nil
	}
	return false, err
}
