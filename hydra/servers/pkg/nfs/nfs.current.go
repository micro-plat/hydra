package nfs

import (
	"fmt"
	"strings"
)

//保存文件到本地NFS服务路径
func Save(name string, buff []byte) (string, error) {
	if currentModule == nil || currentModule.done {
		return "", fmt.Errorf("本地nfs未初始化或已关闭")
	}
	fp, domain, err := currentModule.SaveNewFile("", name, buff)
	if err != nil {
		return "", err
	}
	rpath := fmt.Sprintf("%s/%s", strings.Trim(domain, "/"), strings.Trim(fp.Path, "/"))
	return rpath, nil
}
