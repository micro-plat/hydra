package nfs

import (
	"encoding/json"
	"fmt"
	"io/fs"
	"net/http"
	"os"
	"path/filepath"
	"strings"

	"github.com/micro-plat/lib4go/concurrent/cmap"
	"github.com/micro-plat/lib4go/errs"
	"github.com/micro-plat/lib4go/types"
)

var exclude = []string{".fp"}

//local 本地文件管理
type local struct {
	path        string
	fpPath      string
	currentAddr string
	FPS         cmap.ConcurrentMap
}

//newLocal 构建本地处理服务
func newLocal(path string, currentAddr string) (*local, error) {
	l := &local{
		path:        path,
		fpPath:      filepath.Join(path, ".fp"),
		currentAddr: currentAddr,
		FPS:         cmap.New(8),
	}
	return l, l.check()
}

//check 处理本地文件与指纹不一致，以文件为准
func (l *local) check() error {
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
			v.AddHosts(l.currentAddr)
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

//MergeFPSList 合并外部数据列表
func (l *local) MergeFPS(e *eFileFP) (eFileFPLists, eFileFPLists) {
	return l.MergeFPSList(map[string]*eFileFP{
		e.Path: e,
	})
}

//MergeFPSList 合并外部数据列表
func (l *local) MergeFPSList(list eFileFPLists) (eFileFPLists, eFileFPLists) {
	reports := make(eFileFPLists, 10)
	download := make(eFileFPLists, 10)
	for _, fp := range list {
		nlk, ok := l.FPS.Get(fp.Path)
		if !ok {
			download[fp.Path] = fp
			continue
		}
		lk := nlk.(*eFileFP)
		if fp.AddHosts(lk.Hosts...) {
			l.FPS.Set(fp.Path, fp)
			reports[fp.Path] = fp
		}
	}
	if len(reports) > 0 {
		l.FPWrite(l.FPS.Items())
	}
	return reports, download
}

//SaveFile 保存文件
func (l *local) SaveFile(name string, buff []byte, hosts ...string) (f *eFileFP, err error) {
	//将文件写入本地
	if err := l.Write(name, buff); err != nil {
		return nil, fmt.Errorf("保存文件失败:%w", err)
	}

	//生成crc64并
	fp := &eFileFP{Path: name, CRC64: getCRC64(buff)}
	fp.AddHosts(hosts...)
	fp.AddHosts(l.currentAddr)
	l.FPS.Set(name, fp)
	return fp, l.FPWrite(l.FPS)
}

//GetFile 获取本地文件
func (l *local) GetFile(name string) ([]byte, error) {
	if ok := l.FPS.Has(name); ok {
		return l.Read(name)
	}
	return nil, errs.NewErrorf(http.StatusNotFound, "文件%s,%w", name, errs.ErrNotExist)
}

//FPHas 本地是否存在文件
func (l *local) Has(name string) bool {
	if fx, ok := l.FPS.Get(name); ok {
		f := fx.(*eFileFP)
		return f.Has(l.currentAddr)
	}
	return false
}

//GetFP 获以FP配置
func (l *local) GetFP(name string) (*eFileFP, bool) {
	if fx, ok := l.FPS.Get(name); ok {
		f := fx.(*eFileFP)
		return f, true
	}
	return nil, false
}

//GetFPList 获以FP列表
func (l *local) GetFPList() eFileFPLists {
	list := make(eFileFPLists)
	for k, v := range l.FPS.Items() {
		list[k] = v.(*eFileFP)
	}
	return list
}

//Read 读取文件，本地不存在
func (l *local) Read(name string) ([]byte, error) {
	buff, err := os.ReadFile(filepath.Join(l.path, name))
	if err != nil {
		return nil, fmt.Errorf("读取文件失败:%w", err)
	}
	return buff, nil
}

//Write 写入文件到本地
func (l *local) Write(name string, buff []byte) error {
	err := os.WriteFile(filepath.Join(l.path, name), buff, 0666)
	if err != nil {
		return fmt.Errorf("写文件失败:%w", err)
	}
	return nil
}

//List 文件清单
func (l *local) List(p ...string) ([]string, error) {
	path := types.GetStringByIndex(p, 0, l.path)
	dirEntity, err := os.ReadDir(path)
	if os.IsNotExist(err) {
		return nil, nil
	}
	if err != nil {
		return nil, fmt.Errorf("读取目录失败:%s %v", path, err)
	}
	list := make([]string, 0, len(dirEntity))
	for _, entity := range dirEntity {

		if l.exclude(entity.Name()) {
			continue
		}
		if entity.IsDir() {
			nlist, err := l.List(filepath.Join(path, entity.Name()))
			if err != nil {
				return nil, err
			}
			list = append(list, nlist...)
			continue
		}
		nname := filepath.Join(path, entity.Name())
		if strings.HasPrefix(nname, filepath.Join(l.path)) {
			nname = nname[len(filepath.Join(l.path))+1:]
		}
		list = append(list, nname)
	}
	return list, nil
}

//FPWrite 写入本地文件
func (l *local) FPWrite(content interface{}) error {
	buff, err := json.Marshal(content)
	if err != nil {
		return err
	}
	return os.WriteFile(l.fpPath, buff, 0666)
}

//FPRead 读取指纹信息
func (l *local) FPRead() (eFileFPLists, error) {
	list := make(eFileFPLists)
	buff, err := os.ReadFile(l.fpPath)
	if os.IsNotExist(err) {
		return list, nil
	}
	if err != nil {
		return nil, fmt.Errorf("读取%s失败%w", l.fpPath, err)
	}
	err = json.Unmarshal(buff, &list)
	return list, err
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
func (l *local) exclude(f string) bool {
	for _, ex := range exclude {
		if ex == f {
			return true
		}
	}
	return false
}
