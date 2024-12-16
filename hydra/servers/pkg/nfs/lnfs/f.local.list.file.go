package lnfs

import (
	"path/filepath"
	"sort"
	"strings"

	"github.com/micro-plat/hydra/hydra/servers/pkg/nfs/infs"
)

//GetFileList 获取文件列表
func (l *local) GetFileList(path string, q string, all bool, index int, count int) infs.FileList {
	list := l.getFileList(path, q, all)
	total := index + count
	if index >= list.Len() {
		return make(infs.FileList, 0)
	}
	if total > list.Len() {
		return list[index:]
	}
	return list[index:total]
}

//GetFileList 获以文件列表
func (l *local) getFileList(path string, q string, all bool) infs.FileList {
	list := make(infs.FileList, 0, 1)
	fps := l.GetFPs()
	for k, v := range fps {

		if all || !all && k == filepath.Join(path, v.Name) {

			if !strings.Contains(v.Name, q) {
				continue
			}

			list = append(list, &infs.FileInfo{
				Type:    infs.GetFileType(k),
				Path:    k,
				ModTime: v.ModTime.Format("2006/01/02 15:04:05"),
				Size:    v.Size,
				DPath:   infs.UnMultiPath(k),
				Name:    v.Name,
			})
		}
	}
	sort.Sort(list)
	return list
}
