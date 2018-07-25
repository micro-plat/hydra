package file

import (
	"os"
	"path/filepath"
)

//Exists 检查文件或路径是否存在
func Exists(path string) bool {
	_, err := os.Stat(path)
	return err == nil || os.IsExist(err)
}

//GetAbs 获取文件绝对路径
func GetAbs(path string) (string, error) {
	absPath, err := filepath.Abs(path)
	if err != nil {
		return "", err
	}
	return absPath, nil
}
