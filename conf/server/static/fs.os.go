package static

import (
	"io/fs"
	"net/http"
	"os"
	"path/filepath"
)

//osfs 本地文件系统
type osfs struct {
	dir string
	fs  http.FileSystem
}

func newOSFS(dir string) *osfs {
	return &osfs{
		dir: dir,
		fs:  http.FS(os.DirFS(dir)),
	}
}

func (o *osfs) ReadFile(name string) (http.FileSystem, string, error) {
	return o.fs, filepath.Join(o.dir, name), nil
}

func (o *osfs) GetRoot() string {
	return o.dir
}

func (o *osfs) Has(name string) bool {
	_, err := os.Stat(name)
	return err == nil
}

func (o *osfs) GetDirEntrys(path string) (dirs []fs.DirEntry, err error) {
	return os.ReadDir(filepath.Join(o.dir, path))
}
