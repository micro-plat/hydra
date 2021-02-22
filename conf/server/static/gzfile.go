package static

import (
	"io/fs"
	"net/http"
	"path"
)

//FileInfo 压缩文件保存
type gzFileInfo struct {
	GzFile string
	HasGz  bool
}

//GetGzFile 获取gz 压缩包
func (s *Static) getGzFile(rPath string) string {
	fi, ok := s.gzipfileMap[rPath]
	if !ok {
		return ""
	}
	if !fi.HasGz {
		return ""
	}
	return fi.GzFile
}

//GetGzip GetGzip
func (s *Static) GetGzip() interface{} {
	return s.gzipfileMap
}

func (s *Static) refreshGzip() {
	//root := s.fs.GetRoot()
	root := ""
	entrys, err := s.fs.GetDirEntrys(root)
	if err != nil {
		return
	}
	for _, e := range entrys {
		s.deepSearch(root, e)
	}
}

func (s *Static) deepSearch(parent string, entry fs.DirEntry) {
	cur := path.Join(parent, entry.Name())
	if entry.IsDir() {
		list, err := s.fs.GetDirEntrys(cur)
		if err != nil {
			return
		}
		for _, item := range list {
			s.deepSearch(cur, item)
		}
		return
	}
	gzfileName := cur + ".gz"
	if s.fs.Has(gzfileName) {
		s.gzipfileMap["/"+cur] = gzFileInfo{
			GzFile: gzfileName,
			HasGz:  true,
		}
	}
}

type gzipFSWrapper struct {
	orgFs  IFS
	static *Static
}

//NewGzip NewGzip
func NewGzip(org IFS, static *Static) IFS {
	static.fs = org
	static.refreshGzip()
	return &gzipFSWrapper{
		orgFs:  org,
		static: static,
	}
}

func (w *gzipFSWrapper) ReadFile(name string) (fs http.FileSystem, realPath string, err error) {
	fs, realPath, err = w.orgFs.ReadFile(name)
	if err != nil {
		return
	}
	if gzfile := w.static.getGzFile(name); gzfile != "" {
		realPath = path.Join(w.GetRoot(), gzfile)
	}
	return
}

func (w *gzipFSWrapper) Has(name string) bool {
	return w.orgFs.Has(name)
}

func (w *gzipFSWrapper) GetRoot() string {
	return w.orgFs.GetRoot()
}

func (w *gzipFSWrapper) GetDirEntrys(path string) (dirs []fs.DirEntry, err error) {
	return w.orgFs.GetDirEntrys(path)
}
