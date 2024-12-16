package static

import (
	"errors"
	"io"
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
		fs:  ioFS{localFs(dir)},
		p:   make([]IFS, 0, 0),
	}
}

type ioFS struct {
	fsys fs.FS
}

type ioFile struct {
	file fs.File
}

func (f ioFS) Open(name string) (http.File, error) {
	file, err := f.fsys.Open(name)
	if err != nil {
		return nil, err
	}
	return ioFile{file}, nil
}

func (f ioFile) Close() error               { return f.file.Close() }
func (f ioFile) Read(b []byte) (int, error) { return f.file.Read(b) }
func (f ioFile) Stat() (fs.FileInfo, error) { return f.file.Stat() }

var errMissingSeek = errors.New("io.File missing Seek method")
var errMissingReadDir = errors.New("io.File directory missing ReadDir method")

func (f ioFile) Seek(offset int64, whence int) (int64, error) {
	s, ok := f.file.(io.Seeker)
	if !ok {
		return 0, errMissingSeek
	}
	return s.Seek(offset, whence)
}

func (f ioFile) ReadDir(count int) ([]fs.DirEntry, error) {
	d, ok := f.file.(fs.ReadDirFile)
	if !ok {
		return nil, errMissingReadDir
	}
	return d.ReadDir(count)
}

func (f ioFile) Readdir(count int) ([]fs.FileInfo, error) {
	d, ok := f.file.(fs.ReadDirFile)
	if !ok {
		return nil, errMissingReadDir
	}
	var list []fs.FileInfo
	for {
		dirs, err := d.ReadDir(count - len(list))
		for _, dir := range dirs {
			info, err := dir.Info()
			if err != nil {
				// Pretend it doesn't exist, like (*os.File).Readdir does.
				continue
			}
			list = append(list, info)
		}
		if err != nil {
			return list, err
		}
		if count < 0 || len(list) >= count {
			break
		}
	}
	return list, nil
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
