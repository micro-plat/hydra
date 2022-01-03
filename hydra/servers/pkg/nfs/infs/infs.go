package infs

import (
	"strings"
)

type Infs interface {
	Start() error
	Close() error
	Exists(string) bool

	GetDirList(string, int) DirList
	GetFileList(path string, q string, all bool, index int, count int) FileList

	Save(string, string, []byte) (string, string, error)
	Get(string) ([]byte, string, error)
	CreateDir(string, string) error
	Rename(string, string, string) error
	GetScaleImage(root string, path string, width int, height int, quality int) (buff []byte, ctp string, err error)
	Conver2PDF(root string, path string) (buff []byte, ctp string, err error)
	Registry(tp string)
}

//处理多级目录
func MultiPath(path string) string {
	return strings.Trim(strings.ReplaceAll(path, "|", "/"), "/")
}

//处理多级目录
func UnMultiPath(path string) string {
	return strings.Trim(strings.ReplaceAll(path, "/", "|"), "|")
}
