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
	p   []IFS
}

func newOSFS(dir string) *osfs {
	return &osfs{
		dir: dir,
		fs:  http.FS(localFs(dir)),
		p:   make([]IFS, 0, 0),
	}
}

func (o *osfs) ReadFile(name string) (http.FileSystem, string, error) {
	for _, p := range o.p {
		if _, ok := p.Has(name); ok {
			return p.ReadFile(name)
		}
	}
	if rname, ok := o.Has(name); ok {
		return o.fs, rname, nil
	}
	return o.fs, name, nil
}

func (o *osfs) Has(name string) (string, bool) {
	for _, p := range o.p {
		if v, ok := p.Has(name); ok {
			return v, ok
		}
	}
	names := getName(o.dir, name)
	for _, name := range names {
		info, err := os.Stat(name)
		if err == nil && !info.IsDir() {
			return name, true
		}
	}
	return "", false
}

func (o *osfs) Merge(p IFS) {
	if p == nil {
		return
	}
	o.p = append(o.p, p)
}

type localFs string

func (dir localFs) Open(name string) (fs.File, error) {
	f, err := os.Open(name)
	if err != nil {
		return nil, err
	}
	return f, nil

}
func getName(dir string, name string) []string {
	list := make([]string, 0, 2)
	list = append(list, filepath.Join(dir, name+".gz"))
	list = append(list, filepath.Join(dir, name))
	return list
}
