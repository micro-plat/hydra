package static

import (
	"embed"
	"fmt"
	"io/fs"
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

var defEmbedFs = &embedFs{}

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
}

func newEFS(name string, fs *embed.FS) *efs {
	return &efs{
		fs:   fs,
		name: name,
		hfs:  http.FS(fs),
	}
}

func (o *efs) ReadFile(name string) (http.FileSystem, string, error) { //http.FileServer(http.FS(embed.FS{}))
	return o.hfs, filepath.Join(o.name, name), nil
}

func (o *efs) GetRoot() string {
	return o.name
}

func (o *efs) Has(name string) bool {
	_, err := o.fs.ReadFile(filepath.Join(o.name, name))
	return err == nil
}

func (o *efs) GetDirEntrys(path string) (dirs []fs.DirEntry, err error) {
	return o.fs.ReadDir(filepath.Join(o.name, path))
}
