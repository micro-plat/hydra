package lnfs

import (
	"fmt"
	"hash/crc64"
	"net/http"
	"path/filepath"
	"sync"
	"time"

	"github.com/fsnotify/fsnotify"
	"github.com/micro-plat/hydra/conf/server/nfs"
	"github.com/micro-plat/hydra/hydra/servers/pkg/nfs/infs"
	"github.com/micro-plat/hydra/hydra/servers/pkg/nfs/internal"
	"github.com/micro-plat/lib4go/errs"
	"github.com/micro-plat/lib4go/utility"
)

//Module 协调本地文件、本地指纹、远程指纹等处理
type Module struct {
	c         *nfs.NFS
	Local     *local
	remoting  *remoting
	async     *async
	once      sync.Once
	fsWatcher *fsnotify.Watcher
	checkChan chan struct{}
	prefix    string
	Done      bool
}

func newModule(c *nfs.NFS, prefix string) (m *Module) {
	m = &Module{
		c:         c,
		prefix:    prefix,
		Local:     newLocal(c.Local, c.Excludes, c.Includes),
		remoting:  newRemoting(),
		checkChan: make(chan struct{}, 1),
	}
	m.async = newAsync(m.Local, m.remoting)
	go m.watch()
	return m
}

//Update 更新环境配置
func (m *Module) Update(hosts []string, masterHost string, currentAddr string, isMaster bool) {
	m.remoting.Update(hosts, masterHost, currentAddr, isMaster, m.prefix)
	m.Local.Update(currentAddr)
	if isMaster {
		m.async.DoQuery()
	}
}

//checkAndDownload 判断整个集群是否存在文件
//1. 检查本地是否存在
//2. 从master获取指纹信息，看哪些服务器有此文件
//3. 向有此文件的服务器拉取文件
//4. 保存到本地
//5. 通知master我也有这个文件了,如果是master则告诉所有人我也有此文件了
func (m *Module) CheckAndDownload(name string) error {

	//从本地获取
	if m.Local.Has(name) {
		return nil
	}

	//从远程获取
	fp, err := m.remoting.GetFP(name)
	if err != nil {
		return err
	}

	//从远程拉取文件
	buff, err := m.remoting.Pull(fp)
	if err != nil {
		return err
	}

	//保存到本地
	fp, err = m.Local.SaveFile(name, buff, fp.Hosts...)
	if err != nil {
		return err
	}

	//上报给其它服务器
	m.async.DoReport(fp.GetMAP())
	return nil
}
func (m *Module) Get(name string) ([]byte, string, error) {
	if err := m.CheckAndDownload(name); err != nil {
		return nil, "", err
	}
	ctp := infs.GetContentType(name)
	buff, err := internal.ReadFile(filepath.Join(m.c.Local, name))
	return buff, ctp, err
}

//HasFile 本地是否存在文件
func (m *Module) HasFile(name string) error {
	if m.Local.Has(name) {
		return nil
	}
	return errs.NewErrorf(http.StatusNotFound, "文件%s不存在", name)
}

//HasFile 本地是否存在文件
func (m *Module) Exists(name string) bool {
	if m.HasFile(name) == nil {
		return true
	}
	return m.remoting.HasFile(name) == nil
}

//Save 保存新文件到本地
//1. 查询本地是否有此文件了,有则报错
//2. 保存到本地，返回指纹信息
//3. 通知master我也有这个文件了,如果是master则告诉所有人我也有此文件了
func (m *Module) Save(name string, buff []byte) (string, error) {
	//检查文件是否存在
	name = getFileName(name, m.c.Rename)
	if m.Local.Has(name) {
		return "", fmt.Errorf("文件名称重复:%s", name)
	}

	//保存到本地
	fp, err := m.Local.SaveFile(name, buff)
	if err != nil {
		return "", err
	}

	//远程通知
	m.async.DoReport(fp.GetMAP())
	return fp.Path, nil
}

//GetFP 获取本地的指纹信息，用于master对外提供服务
//1. 查询本地是否有文件的指纹信息
//2. 如果是master返回不存在
//3. 向master发起查询
func (m *Module) GetFP(name string) (*eFileFP, error) {
	//从本地文件获取
	if f, ok := m.Local.GetFP(name); ok {
		return f, nil
	}
	return nil, errs.NewError(http.StatusNotFound, "文件不存在")
}

//Query 获取本地包含有服务列表的指纹清单
func (m *Module) Query() EFileFPLists {
	return m.Local.GetFPs()
}

//GetFileList 获取文件列表
func (m *Module) GetFileList(path string, q string, all bool, index int, count int) infs.FileList {
	return m.Local.GetFileList(path, q, all, index, count)
}

//GetDirList 获取目录列表
func (m *Module) GetDirList(path string, deep int) infs.DirList {
	return m.Local.GetDirList(path, deep)
}

//RecvNotify 接收远程发过来的新文件通知
//1. 检查本地是否有些文件
//2. 文件不存在则自动下载
//3. 合并服务列表
func (m *Module) RecvNotify(f EFileFPLists) error {
	//处理本地新文件上报
	reports, downloads, err := m.Local.Merge(f)
	if err != nil {
		return err
	}

	//上报到服务
	m.async.DoReport(reports)
	for _, f := range downloads {
		m.async.DoDownload(f)
	}
	return m.Local.FPWrite(m.Local.FPS.Items())
}
func (c *Module) CreateDir(path string) error {
	return internal.CreateDir(filepath.Join(c.c.Local, path))
}

func (c *Module) Rename(oname string, nname string) error {
	return internal.Rename(filepath.Join(c.c.Local, oname), filepath.Join(c.c.Local, nname))
}
func (c *Module) GetScaleImage(path string, width int, height int, quality int) (buff []byte, ctp string, err error) {

	ctp = infs.GetContentType(path)
	buff, err = internal.ScaleImageByPath(c.c.Local, path, width, height, quality)
	if err == nil {
		return buff, ctp, err
	}
	buff, err = internal.ReadFile(filepath.Join(c.c.Local, path))
	return buff, ctp, err
}
func (c *Module) Conver2PDF(path string) (buff []byte, ctp string, err error) {
	buff, ctp, _, err = internal.Conver2PDF(c.c.Local, filepath.Join(c.c.Local, path))
	return buff, ctp, err
}

//Close 关闭服务
func (m *Module) Close() error {
	m.Done = true
	m.once.Do(func() {
		close(m.checkChan)
		m.Local.Close()
		m.async.Close()
		if m.fsWatcher != nil {
			m.fsWatcher.Close()
		}
	})
	return nil
}

func getCRC64(buff []byte) uint64 {
	return crc64.Checksum(buff, crc64.MakeTable(crc64.ISO))
}

func getFileName(name string, rename bool) string {
	if !rename {
		return filepath.Join(name)
	}
	ext := filepath.Ext(name)
	return filepath.Join(time.Now().Format("20060102"), fmt.Sprintf("%d%s", fnv32(utility.GetGUID()), ext))

}
func fnv32(key string) uint32 {
	hash := uint32(2166136261)
	const prime32 = uint32(16777619)
	for i := 0; i < len(key); i++ {
		hash *= prime32
		hash ^= uint32(key[i])
	}
	return hash
}
