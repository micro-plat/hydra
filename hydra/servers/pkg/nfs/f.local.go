package nfs

import (
	"io/fs"
	"os"
	"path/filepath"
	"sync"

	"github.com/micro-plat/lib4go/concurrent/cmap"
)

var exclude = []string{".fp"}

//local 本地文件管理
type local struct {
	path        string
	fpPath      string
	currentAddr string
	FPS         cmap.ConcurrentMap
	once        sync.Once
	readyChan   chan struct{}
}

//newLocal 构建本地处理服务
func newLocal(path string) *local {
	l := &local{
		path:      path,
		fpPath:    filepath.Join(path, ".fp"),
		readyChan: make(chan struct{}),
		FPS:       cmap.New(8),
	}
	return l
}

func (l *local) Update(path string, currentAddr string) {
	needCheck := l.path != path || l.currentAddr != currentAddr
	l.path = path
	l.currentAddr = currentAddr
	l.fpPath = filepath.Join(path, ".fp")
	if needCheck {
		l.check()
	}
	return
}

//check 处理本地文件与指纹不一致，以文件为准
func (l *local) check() error {
	defer close(l.readyChan)
	//读取本地指纹
	fps, err := l.FPRead()
	if err != nil {
		return err
	}

	//获取本地文件列表
	lst, err := l.List()
	if err != nil {
		return err
	}

	//处理不一致数据
	for _, path := range lst {
		if v, ok := fps[path]; ok {
			v.MergeHosts(l.currentAddr)
			l.FPS.Set(path, v)
			continue
		}
		buff, err := l.Read(path)
		if err != nil {
			return err
		}
		fp := &eFileFP{Path: path, CRC64: getCRC64(buff), Hosts: []string{l.currentAddr}}
		l.FPS.Set(path, fp)
	}
	//更新数据
	return l.FPWrite(l.FPS.Items())
}

//MergeLocal 合并到本地列表
func (l *local) MergeLocal(list eFileFPLists) (reports eFileFPLists, download eFileFPLists, err error) {
	reports = make(eFileFPLists, 10)
	download = make(eFileFPLists, 10)
	for _, fp := range list {
		nlk, ok := l.FPS.Get(fp.Path)
		if !ok {
			download[fp.Path] = fp
			continue
		}
		lk := nlk.(*eFileFP)
		v0 := fp.MergeHosts(l.currentAddr)
		v1 := lk.MergeHosts(fp.Hosts...)
		v2 := fp.MergeHosts(lk.Hosts...)
		if v0 || v1 || v2 {
			reports[fp.Path] = fp
		}
	}
	if len(reports) > 0 {
		err = l.FPWrite(l.FPS.Items())
	}
	return reports, download, err
}

func GetAllNotify(fps eFileFPLists, allAliveHosts []string) map[string]eFileFPLists {
	if len(fps) == 0 {
		return nil
	}
	list := make(map[string]eFileFPLists, len(allAliveHosts))
	for _, v := range allAliveHosts {
		list[v] = fps
	}
	return list
}

//GetFile 获取本地文件
func (l *local) GetFile(name string) ([]byte, error) {
	return l.Read(name)
}

//FPHas 本地是否存在文件
func (l *local) Has(name string) bool {
	_, ok := l.FPS.Get(name)
	return ok
}

//Open 读取文件
func (l *local) Open(name string) (fs.File, error) {
	return os.Open(filepath.Join(l.path, name))
}

//Close 将缓存数据写入本地文件
func (l *local) Close() error {
	l.FPS = nil
	return nil
}
