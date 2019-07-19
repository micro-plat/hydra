package engines

import (
	"fmt"
	"os"
	"plugin"
	"reflect"
	"sync"

	"github.com/micro-plat/lib4go/file"
)

var components = make(map[string]ServiceLoader)
var mu sync.Mutex

func getComponent(p string) (f ServiceLoader, err error) {
	path, err := file.GetAbs(p)
	if err != nil {
		return
	}
	if p, ok := components[path]; ok {
		return p, nil
	}
	mu.Lock()
	defer mu.Unlock()
	if p, ok := components[path]; ok {
		return p, nil
	}
	if _, err = os.Lstat(path); err != nil && os.IsNotExist(err) {
		return nil, nil
	}
	pg, err := plugin.Open(path)
	if err != nil {
		return nil, fmt.Errorf("加载引擎插件失败:%s,err:%v", path, err)
	}
	work, err := pg.Lookup("GetLoader")
	if err != nil {
		return nil, fmt.Errorf("加载引擎插件%s失败未找到函数GetLoader,err:%v", path, err)
	}
	wkr, ok := work.(ServiceLoader)
	if !ok {
		return nil, fmt.Errorf("加载引擎插件%s失败 GetLoader函数必须为ServiceLoader类型", path)
	}
	components[p] = wkr
	return components[p], nil

}

//LoadComponents 加载所有插件
func (r *ServiceEngine) LoadComponents(files ...string) error {
	for _, file := range files {
		//根据加载的文件名，获取组件
		comp, err := getComponent(file)
		if err != nil {
			return err
		}
		if comp == nil || reflect.ValueOf(comp).IsNil() {
			continue
		}
		if err = comp(r.StandardComponent, r); err != nil {
			return err
		}
	}
	return nil
}
