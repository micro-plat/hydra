package component

// import (
// 	"fmt"
// 	"os"
// 	"plugin"
// 	"sync"

// 	"github.com/micro-plat/lib4go/file"
// )

// type funcDecrypt func(string) (string, error)

// var decryptPlugins = make(map[string]*plugin.Plugin)
// var decryptFuncs = make(map[string]funcDecrypt)
// var mu sync.Mutex
// var dlName = "./libcrypto.so"
// var fNames = []string{"Decrypt", "RsaDecrypt", "DesDecrypt"}

// func init() {
// 	for _, name := range fNames {
// 		getDecrypt(name, dlName)
// 	}
// }
// func getFunc(name ...string) (funcDecrypt, bool) {
// 	for _, n := range name {
// 		if f, ok := decryptFuncs[n]; ok {
// 			return f, ok
// 		}
// 	}
// 	return nil, false
// }

// func getPlugin(dl string) (*plugin.Plugin, error) {
// 	mu.Lock()
// 	defer mu.Unlock()
// 	if p, ok := decryptPlugins[dl]; ok {
// 		return p, nil
// 	}
// 	path, err := file.GetAbs(dl)
// 	if err != nil {
// 		return nil, err
// 	}
// 	if _, err = os.Lstat(path); err != nil && os.IsNotExist(err) {
// 		return nil, fmt.Errorf("未找到插件:%s,err:%v", path, err)
// 	}
// 	pg, err := plugin.Open(path)
// 	if err != nil {
// 		return nil, fmt.Errorf("加载引擎插件失败:%s,err:%v", path, err)
// 	}
// 	decryptPlugins[dl] = pg
// 	return pg, nil
// }

// func getDecrypt(funcName string, dl string) (f funcDecrypt, err error) {
// 	mu.Lock()
// 	defer mu.Unlock()
// 	f, ok := decryptFuncs[funcName]
// 	if ok {
// 		return f, nil
// 	}
// 	pg, err := getPlugin(dl)
// 	if err != nil {
// 		return nil, err
// 	}
// 	work, err := pg.Lookup(funcName)
// 	if err != nil {
// 		return nil, fmt.Errorf("加载引擎插件%s失败未找到函数%s,err:%v", dl, funcName, err)
// 	}
// 	wkr, ok := work.(func(string) (string, error))

// 	if !ok || wkr == nil {
// 		return nil, fmt.Errorf("加载引擎插件%s失败 %s函数必须为func(string) (string, error)类型", dl, funcName)
// 	}
// 	decryptFuncs[funcName] = wkr
// 	return wkr, nil
// }

// //decrypt 解密数据
// func decrypt(i string) (string, error) {
// 	if f, ok := getFunc(fNames...); ok {
// 		return f(i)
// 	}
// 	return i, nil
// }
