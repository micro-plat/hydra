package nfs

import (
	"encoding/json"
	"fmt"
	"os"
)

//FPHas 本地是否存在文件
func (l *local) Has(name string) bool {
	return l.FPS.Has(name)
}

//GetFP 获以FP配置
func (l *local) GetFP(name string) (*eFileFP, bool) {
	if fx, ok := l.FPS.Get(name); ok {
		f := fx.(*eFileFP)
		return f, true
	}
	return nil, false
}

//GetFPs 获以FP列表
func (l *local) GetFPs() eFileFPLists {
	l.nfsChecker.Wait()
	list := make(eFileFPLists)
	for k, v := range l.FPS.Items() {
		list[k] = v.(*eFileFP)
	}
	return list
}

//FPWrite 写入本地文件
func (l *local) FPWrite(content interface{}) error {
	buff, err := json.Marshal(content)
	if err != nil {
		return err
	}
	_, s := os.Stat(l.path)
	if os.IsNotExist(s) {
		os.MkdirAll(l.path, 0777)
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
	if len(buff) == 0 {
		return list, nil
	}
	err = json.Unmarshal(buff, &list)
	return list, err
}
