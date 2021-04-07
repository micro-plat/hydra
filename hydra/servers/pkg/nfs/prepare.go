package nfs

import (
	"fmt"
	"io/fs"
	"os"
	"path/filepath"
	"strings"

	"github.com/micro-plat/lib4go/types"
)

type NFS struct {
	unfinishedList map[string][]string
	fileList       types.XMap
}

//构建本地nfs文件系统
func newNFS(local string) (nfs *NFS, err error) {
	nfs = &NFS{}

	//读取文件
	nfs.fileList, err = readFiles(local)
	if err != nil {
		return nil, err
	}

	//构建未完成列表
	nfs.unfinishedList, err = getUnfinishedList(local)
	if err != nil {
		return nil, err
	}
	return nfs, nil
}

//WriteFile 写入文件到本地
func (n *NFS) WriteFile(name string, buff []byte) error {
	return nil
}

//ReadFile 读取文件，本地不存在时从远程读取
func (n *NFS) ReadFile(name string) (fs.FileInfo, error) {
	return nil, nil
}

//FindFile 查询文件是否存在
func (n *NFS) FindFile(name string) bool {
	return false
}

//RecvFile 接收集群内服务器发送过来的文件
func (n *NFS) RecvFile(name string, buff []byte) error {
	return nil
}

//ProvideFileList 提交本地文件清单
func (n *NFS) ProvideFileList() ([]string, error) {
	return nil, nil
}

//获取未完成列表(接收到新文件时写入，完成集群中所有机器同步后移除)
func getUnfinishedList(name string) (map[string][]string, error) {
	_, err := os.Stat(name)
	if err != nil {
		if err == os.ErrNotExist {
			return nil, nil
		}
		return nil, err
	}
	buff, err := os.ReadFile(name)
	if err != nil {
		return nil, err
	}
	files := strings.Split(strings.Trim(string(buff), ","), ",")
	list := make(map[string][]string, len(files))
	for _, f := range files {
		list[f] = make([]string, 0, 0)
	}
	return list, nil
}

//获取本地文件列表(用于本地下载和集群机器获取)
func readFiles(name string) (types.XMap, error) {
	if strings.HasPrefix(name, ".") {
		return nil, nil
	}
	dirEntity, err := os.ReadDir(name)
	if err != nil {
		return nil, fmt.Errorf("读取目录失败:%s %v", name, err)
	}
	list := types.NewXMap()
	for _, entity := range dirEntity {
		if entity.IsDir() {
			nlist, err := readFiles(filepath.Join(name, entity.Name()))
			if err != nil {
				return nil, err
			}
			list.Merge(nlist)
			continue
		}
		_, err := entity.Info()
		if err != nil {
			return nil, err
		}
		list.SetValue(filepath.Join(name, entity.Name()), entity.Name())
	}
	return list, nil
}
