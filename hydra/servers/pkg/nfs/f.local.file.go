package nfs

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

//SaveFile 保存文件
func (l *local) SaveFile(name string, buff []byte, hosts ...string) (f *eFileFP, err error) {

	//将文件写入本地
	if err := l.FWrite(name, buff); err != nil {
		return nil, fmt.Errorf("保存文件失败:%w", err)
	}

	//生成crc64并
	fp := &eFileFP{
		Path: name,
		// CRC64: getCRC64(buff),
	}
	fp.MergeHosts(hosts...)
	fp.MergeHosts(l.currentAddr)
	l.FPS.Set(name, fp)
	return fp, l.FPWrite(l.FPS)
}

//FRead 读取文件，本地不存在
func (l *local) FRead(name string) ([]byte, error) {
	buff, err := os.ReadFile(filepath.Join(l.path, name))
	if err != nil {
		return nil, fmt.Errorf("读取文件失败:%w", err)
	}
	return buff, nil
}

//FWrite 写入文件到本地
func (l *local) FWrite(name string, buff []byte) error {
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

//FList 获取本地所有文件清单
func (l *local) FList(path string) ([]string, error) {

	//文件夹不存在时返回空
	dirEntity, err := os.ReadDir(path)
	if os.IsNotExist(err) {
		return nil, nil
	}
	if err != nil {
		return nil, fmt.Errorf("读取目录失败:%s %v", path, err)
	}

	//查找所有文件
	list := make([]string, 0, len(dirEntity))
	for _, entity := range dirEntity {
		if l.exclude(entity.Name()) {
			continue
		}
		//处理目录
		if entity.IsDir() {
			nlist, err := l.FList(filepath.Join(path, entity.Name()))
			if err != nil {
				return nil, err
			}
			list = append(list, nlist...)
			continue
		}

		//处理文件名称
		nname := filepath.Join(path, entity.Name())
		if strings.HasPrefix(nname, filepath.Join(l.path)) {
			nname = nname[len(filepath.Join(l.path))+1:]
		}
		list = append(list, nname)
	}
	return list, nil
}

var exclude = []string{".", "~"}

func (l *local) exclude(f string) bool {
	for _, ex := range exclude {
		if strings.HasPrefix(f, ex) {
			return true
		}
	}
	// if strings.HasPrefix(f, filepath.Join(l.path, time.Now().Format("20060102"))) {
	// 	return true
	// }
	return false
}
func (l *local) FindChange() bool {
	//获取本地文件列表
	change := false
	lst, err := l.FList(l.path)
	if err != nil {
		return false
	}

	//处理不一致数据
	for _, path := range lst {
		if ok := l.FPS.Has(path); !ok {
			fp := &eFileFP{
				Path:  path,
				Hosts: []string{l.currentAddr},
			}
			l.FPS.Set(path, fp)
			change = true
		}
	}
	return change
}
