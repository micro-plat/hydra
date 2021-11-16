package nfs

import (
	"fmt"
	"io/ioutil"
	"net/http"
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

//下载文件
func Download(name string) ([]byte, error) {
	if err := currentModule.checkAndDownload(name); err != nil {
		return nil, err
	}
	fs, err := http.FS(currentModule.local).Open(name)
	if err != nil {
		return nil, err
	}
	defer fs.Close()
	return ioutil.ReadAll(fs)
}
