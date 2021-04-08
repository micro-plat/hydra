package nfs

import (
	"encoding/json"
	"fmt"
	"io/fs"
	"os"
	"path/filepath"

	"github.com/micro-plat/lib4go/errs"
	"github.com/micro-plat/lib4go/types"
)

//local 本地文件管理
type local struct {
	path        string
	fpPath      string
	currentAddr string
	FPS         eFileFPLists
	NFPS        eFileFPLists
}

//newLocal 构建本地处理服务
func newLocal(path string) (*local, error) {
	l := &local{path: path, fpPath: filepath.Join(path, ".fp")}
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
	l.FPS = make(eFileFPLists, len(lst))
	l.NFPS = make(eFileFPLists, 0)
	for _, path := range lst {
		if v, ok := fps[path]; ok {
			v.AddHosts(l.currentAddr)
			l.FPS[path] = v
		} else {
			buff, err := l.Read(path)
			if err != nil {
				return err
			}
			fp := &eFileFP{Path: path, CRC64: getCRC64(buff), Hosts: []string{l.currentAddr}}
			l.FPS[path] = fp
			l.NFPS[path] = fp
		}
	}

	//更新数据
	if len(l.NFPS) > 0 {
		return l.FPWrite(l.FPS)
	}
	return nil
}

//MergeFPS 合并外部数据列表
func (l *local) MergeFPS(list eFileFPLists) {
	for _, fp := range list {
		if _, ok := l.FPS[fp.Path]; ok {
			fp.AddHosts(l.currentAddr)
			l.FPS[fp.Path] = fp
			l.NFPS[fp.Path] = fp
		}
	}
}

//SaveFile 保存文件
func (l *local) SaveFile(name string, buff []byte, hosts ...string) (f *eFileFP, err error) {
	//将文件写入本地
	path := filepath.Join(l.path, name)
	if err := l.Write(path, buff); err != nil {
		return nil, err
	}

	//生成crc64并
	fp := &eFileFP{Path: path, CRC64: getCRC64(buff)}
	fp.AddHosts(hosts...)
	fp.AddHosts(l.currentAddr)
	l.FPS[path] = fp
	return fp, l.FPWrite(l.FPS)
}

//GetFile 获取本地文件
func (l *local) GetFile(name string) ([]byte, error) {
	path := filepath.Join(l.path, name)
	if _, ok := l.FPS[path]; ok {
		return l.Read(path)
	}
	return nil, fmt.Errorf("文件%s %w", path, errs.ErrNotExist)
}

//FPHas 本地是否存在文件
func (l *local) Has(name string) bool {
	path := filepath.Join(l.path, name)
	if f, ok := l.FPS[path]; ok {
		return f.Has(l.currentAddr)
	}
	return false
}

//Has 本地是否存在文件
func (l *local) MergeHost(f *eFileFP) {
	if _, ok := l.FPS[f.Path]; !ok {
		l.FPS[f.Path] = f
		return
	}
	l.FPS[f.Path].AddHosts(f.Hosts...)
	return
}

//Read 读取文件，本地不存在
func (l *local) Read(name string) ([]byte, error) {
	return os.ReadFile(filepath.Join(l.path, name))
}

//Write 写入文件到本地
func (l *local) Write(name string, buff []byte) error {
	return os.WriteFile(filepath.Join(l.path, name), buff, 0666)
}

//List 文件清单
func (l *local) List(p ...string) ([]string, error) {
	path := types.GetStringByIndex(p, 0, l.path)
	dirEntity, err := os.ReadDir(path)
	if err != nil {
		return nil, fmt.Errorf("读取目录失败:%s %v", path, err)
	}
	list := make([]string, 0, len(dirEntity))
	for _, entity := range dirEntity {
		if entity.IsDir() {
			nlist, err := l.List(filepath.Join(path, entity.Name()))
			if err != nil {
				return nil, err
			}
			list = append(list, nlist...)
			continue
		}
		list = append(list, filepath.Join(path, entity.Name()))
	}
	return list, nil
}

//FPWrite 写入本地文件
func (l *local) FPWrite(content eFileFPLists) error {
	buff, err := json.Marshal(content)
	if err != nil {
		return err
	}
	return os.WriteFile(l.fpPath, buff, 0666)
}

//FPRead 读取指纹信息
func (l *local) FPRead() (eFileFPLists, error) {
	buff, err := os.ReadFile(l.fpPath)
	if err == os.ErrNotExist {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}
	list := make(eFileFPLists)
	err = json.Unmarshal(buff, &list)
	return list, err
}

//Open 读取文件
func (l *local) Open(name string) (fs.File, error) {
	return os.Open(filepath.Join(l.path, name))
}

//Close 将缓存数据写入本地文件
func (l *local) Close() error {
	return l.FPWrite(l.FPS)
}
