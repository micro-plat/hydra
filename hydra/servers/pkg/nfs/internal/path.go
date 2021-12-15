package internal

import (
	"bytes"
	"os"
	"os/exec"
)

func interactiveToexec(commandName string, params []string) (string, bool) {
	cmd := exec.Command(commandName, params...)
	buf, err := cmd.Output()
	w := bytes.NewBuffer(nil)
	cmd.Stderr = w
	if err != nil {
		return "", false
	} else {
		return string(buf), true
	}
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
