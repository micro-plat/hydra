//+build go1.8

package service

import (
	"fmt"
	"os"
	"path/filepath"
)

func (c *Config) execPath() (string, error) {
	if len(c.Executable) != 0 {
		fmt.Println("1")
		return filepath.Abs(c.Executable)
	}
	fmt.Println("2")
	return os.Executable()
}
