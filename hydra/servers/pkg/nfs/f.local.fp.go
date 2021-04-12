package nfs

import (
	"encoding/json"
	"fmt"
	"os"
)

//GetFPByName 获以FP配置
func (l *local) GetFPByName(name string) (*eFileFP, bool) {
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
	err = json.Unmarshal(buff, &list)
	return list, err
}
