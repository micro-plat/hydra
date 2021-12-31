package nfs

import (
	"fmt"
	"strings"
)

// Save 保存文件到本地NFS服务路径
func Save(name string, buff []byte) (string, error) {
	if currentModule == nil {
		return "", fmt.Errorf("本地nfs未初始化或已关闭")
	}
	fp, domain, err := currentModule.Save("", name, buff)
	if err != nil {
		return "", err
	}
	rpath := fmt.Sprintf("%s/%s", strings.Trim(domain, "/"), strings.Trim(fp, "/"))
	return rpath, nil
}

//Exists 检查本地和远程文件是否存在
func Exists(name string) bool {
	if currentModule == nil {
		return false
	}
	return currentModule.Exists(name)
}

//Download 下载文件
func Download(name string) ([]byte, error) {
	if currentModule == nil {
		return nil, fmt.Errorf("本地nfs未初始化或已关闭")
	}
	buff, _, err := currentModule.Get(name)
	return buff, err
}
