package static

//本地文件系统,内嵌等文件提供统一的读取接口

import (
	"io/fs"
	"net/http"
)

//IFS 文件管理器
type IFS interface {
	Has(name string) bool
	GetRoot() string
	GetDirEntrys(path string) (dirs []fs.DirEntry, err error)
	ReadFile(name string) (fs http.FileSystem, realPath string, err error)
}
