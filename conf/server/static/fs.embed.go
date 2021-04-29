package static

import (
	"embed"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"path/filepath"

	"github.com/mholt/archiver"
)

//处理内嵌文件，压缩包则解压并保存到本地，否则通过内存映射

type embedFs struct {
	name    string
	archive *embed.FS
	bytes   []byte
}

var defEmbedFs = map[string]*embedFs{}

//check2FS 检查并转换为fs类型
func (e *embedFs) getFileEmbed() (IFS, error) {
	if len(e.bytes) > 0 {
		return newEmbedFile(e.name, e.bytes)
	}
	if e.archive != nil {
		return newEFS(e.name, e.archive), nil
	}
	return nil, nil
}

func newEmbedFile(name string, buff []byte) (*osfs, error) {
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
	defer func(fileName string) {
		os.Remove(fileName)
	}(file.Name())
	return unarchive(file.Name())
}

//efs 扩展embed.fs
type efs struct {
	fs   *embed.FS
	hfs  http.FileSystem
	name string
	p    []IFS
}

func newEFS(name string, fs *embed.FS) *efs {
	return &efs{
		fs:   fs,
		name: name,
		hfs:  http.FS(fs),
		p:    make([]IFS, 0, 0),
	}
}

func (o *efs) ReadFile(name string) (http.FileSystem, string, error) {
	for _, p := range o.p {
		if _, ok := p.Has(name); ok {
			return p.ReadFile(name)
		}
	}
	return o.hfs, filepath.Join(o.name, name), nil
}

func (o *efs) Has(name string) (string, bool) {
	for _, p := range o.p {
		if v, ok := p.Has(name); ok {
			return v, true
		}
	}
	path := filepath.Join(o.name, name)
	_, err := o.fs.ReadFile(path)
	return path, err == nil
}

func (o *efs) Merge(p IFS) {
	if p == nil {
		return
	}
	o.p = append(o.p, p)
}
