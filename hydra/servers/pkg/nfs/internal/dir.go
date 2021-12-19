package internal

import (
	"os"
	"path/filepath"
)

func CreateDir(dir string, p string) error {
	path := filepath.Join(dir, p)
	ok, err := pathExists(path)
	if err != nil {
		return err
	}
	if ok {
		return nil
	}
	return os.MkdirAll(path, 0755)
}
func Rename(dir string, p string, x string) error {
	npath := filepath.Join(dir, x)
	ok, err := pathExists(npath)
	if err != nil {
		return err
	}
	if ok {
		return nil
	}
	return os.Rename(filepath.Join(dir, p), npath)
}
