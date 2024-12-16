package internal

import (
	"os"
)

func CreateDir(path string) error {
	ok, err := pathExists(path)
	if err != nil {
		return err
	}
	if ok {
		return nil
	}
	return os.MkdirAll(path, 0755)
}
func Rename(p string, x string) error {
	ok, err := pathExists(p)
	if err != nil {
		return err
	}
	if ok {
		return nil
	}
	return os.Rename(p, x)
}
