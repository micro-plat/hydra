package static

//本地文件系统,内嵌等文件提供统一的读取接口

import (
	"embed"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"path/filepath"

	"github.com/mholt/archiver"
)

//IFS 文件管理器
type IFS interface {
	ReadFile(name string) (http.FileSystem, string, error)
	Has(name string) bool
}

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
func (o *osfs) Has(name string) bool {
	_, err := os.Stat(name)
	return err == nil
}
func newEmbedFile(name string, buff []byte) (IFS, error) {
	//单个文件,直接返回
	if _, err := archiver.ByExtension(name); err != nil {
		return nil, err
	}

	//压缩文件，进行解压等返回
	rootPath := filepath.Dir(os.Args[0])
	path := fmt.Sprintf("%s%s", TempArchiveName, name)
	file, err := ioutil.TempFile(rootPath, path)
	if err != nil {
		return nil, err
	}
	_, err = file.Write(buff)
	if err != nil {
		return nil, err
	}
	if err = file.Close(); err != nil {
		return nil, err
	}
	return unarchive(file.Name())
}

//efs 扩展embed.fs
type efs struct {
	fs   embed.FS
	hfs  http.FileSystem
	name string
}

func newEFS(name string, fs embed.FS) *efs {
	return &efs{
		fs:   fs,
		name: name,
		hfs:  http.FS(fs),
	}
}

func (o *efs) ReadFile(name string) (http.FileSystem, string, error) { //http.FileServer(http.FS(embed.FS{}))
	return o.hfs, filepath.Join(o.name, name), nil
}

func (o *efs) Has(name string) bool {
	_, err := o.fs.ReadFile(filepath.Join(o.name, name))
	return err == nil
}
