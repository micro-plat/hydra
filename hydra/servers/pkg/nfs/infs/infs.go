package infs

import (
	"path"
	"path/filepath"
	"strings"
)

type Infs interface {
	Start() error
	Close() error
	Exists(string) bool

	GetDirList(string, int) DirList
	GetFileList(path string, q string, all bool, index int, count int) FileList

	Save(string, []byte) (string, error)
	Get(string) ([]byte, string, error)
	CreateDir(string) error
	Rename(string, string) error
	GetScaleImage(path string, width int, height int, quality int) (buff []byte, ctp string, err error)
	Conver2PDF(path string) (buff []byte, ctp string, err error)
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

func GetFileType(p string) string {
	return strings.ToLower(strings.Trim(filepath.Ext(p), "."))
}
func GetFullFileName(p string) string {
	return path.Base(p)
}
func GetFileName(p string) string {
	filenameall := GetFullFileName(p)
	ext := path.Ext(p)
	return filenameall[0 : len(filenameall)-len(ext)]
}

var exclude = []string{".", "~"}

func Exclude(p string, excludes []string, includes ...string) (v bool) {
	f := GetFullFileName(p)
	for _, ex := range exclude {
		if strings.HasPrefix(f, ex) {
			return true
		}
	}
	for _, v := range excludes {
		if strings.Contains(p, v) {
			return true
		}
	}
	if len(includes) > 0 {
		for _, v := range includes {
			if strings.Contains(p, v) {
				return false
			}
		}
		return true
	}

	return false
}
