package nfs

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/micro-plat/lib4go/types"
)

//SaveFile 保存文件
func (l *local) SaveFile(name string, buff []byte, hosts ...string) (f *eFileFP, err error) {
	//将文件写入本地
	if err := l.Write(name, buff); err != nil {
		return nil, fmt.Errorf("保存文件失败:%w", err)
	}

	//生成crc64并
	fp := &eFileFP{Path: name, CRC64: getCRC64(buff)}
	fp.MergeHosts(hosts...)
	fp.MergeHosts(l.currentAddr)
	l.FPS.Set(name, fp)
	return fp, l.FPWrite(l.FPS)
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
	rpath := filepath.Join(l.path, name)

	//处理目录
	dir := filepath.Dir(rpath)
	_, err := os.Stat(dir)
	if os.IsNotExist(err) {
		os.MkdirAll(dir, 0777)
	}

	//生成文件
	err = os.WriteFile(rpath, buff, 0666)
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
func (l *local) exclude(f string) bool {
	for _, ex := range exclude {
		if ex == f {
			return true
		}
	}
	return false
}
