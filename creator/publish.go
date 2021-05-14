package creator

import (
	"encoding/json"
	"fmt"
	"reflect"

	"github.com/manifoldco/promptui"
	"github.com/micro-plat/hydra/conf/pkgs/security"
	"github.com/micro-plat/hydra/conf/server"
	varpub "github.com/micro-plat/hydra/conf/vars"
	"github.com/micro-plat/hydra/global"
	"github.com/micro-plat/hydra/registry"
	"github.com/micro-plat/lib4go/types"
)

//Pub 将配置发布到配置中心
func (c *conf) Pub(platName string, systemName string, clusterName string, registryAddr string, input types.XMap) error {

	if err := c.Load(); err != nil {
		return err
	}

	//检查导入配置
	tps := make([]string, 0, len(c.data))
	for k := range c.data {
		tps = append(tps, k)
	}
	if err := checkImportConf(tps, platName, systemName, clusterName, input); err != nil {
		return err
	}

	//创建注册中心，根据注册中心提供的接口进行配置发布
	r, err := registry.GetRegistry(registryAddr, global.Def.Log())
	if err != nil {
		return err
	}

	confs := make(map[string]interface{})
	cache := types.XMap{}
	//加入server配置
	for tp, subs := range c.data {
		pub := server.NewServerPub(platName, systemName, tp, clusterName)
		path := pub.GetServerPath()
		value, err := getValue(path, subs.Map()[ServerMainNodeName], input)
		if err != nil {
			return err
		}
		//先发布main节点配置
		if err := publish(r, path, value, cache); err != nil {
			return err
		}
		confs[path] = value

		//发布子配置
		for name, value := range subs.Map() {
			path := pub.GetSubConfPath(name)
			if name == ServerMainNodeName {
				continue
			}
			value, err := getValue(path, value, input)
			if err != nil {
				return err
			}
			if err := publish(r, path, value, cache); err != nil {
				return err
			}
			confs[path] = value
		}
	}

	//加入var配置
	for tp, subs := range c.vars {
		pub := varpub.NewVarPub(platName)
		for k, v := range subs {
			path := pub.GetVarPath(tp, k)
			value, err := getValue(path, v, input)
			if err != nil {
				return err
			}
			confs[path] = value
			if err := publish(r, path, value, cache); err != nil {
				return err
			}
		}
	}

	//加入项目未配置的导入配置项
	for k, v := range input {
		if _, ok := confs[k]; !ok {
			if err := publish(r, k, v, cache); err != nil {
				return err
			}
		}
	}

	return nil
}

func publish(r registry.IRegistry, path string, v interface{}, input types.XMap) error {
	value, err := getJSON(path, v, input)
	if err != nil {
		return err
	}
	if b, _ := r.Exists(path); b {
		buff, _, err := r.GetValue(path)
		if err != nil {
			return err
		}
		if !checkCover(string(buff)) { //不覆盖配置则退出
			return nil
		}
		if err := deleteAll(r, path); err != nil {
			return err
		}
	}

	if err := r.CreatePersistentNode(path, value); err != nil {
		return fmt.Errorf("创建配置节点%s %s出错:%w", path, value, err)
	}
	return nil
}

func deleteAll(r registry.IRegistry, path string) error {
	if b, err := r.Exists(path); err != nil || !b {
		return err
	}
	list, err := getAllPath(r, path)
	if err != nil {
		return err
	}
	for _, v := range list {
		if err := r.Delete(v); err != nil {
			return err
		}
	}
	return nil

}

func getAllPath(r registry.IRegistry, path string) ([]string, error) {
	child, _, err := r.GetChildren(path)
	if err != nil {
		return nil, err
	}
	list := make([]string, 0, len(child))
	for _, c := range child {
		npath := registry.Join(path, c)
		nlist, err := getAllPath(r, npath)
		if err != nil {
			return nil, err
		}
		list = append(list, nlist...)
	}
	list = append(list, path)
	return list, nil

}

//getJSON 将对象序列化为json字符串
func getJSON(path string, v interface{}, input types.XMap) (value string, err error) {
	if err := checkAndInput(path, reflect.ValueOf(v), []string{}, input); err != nil {
		return "", err
	}
	if x, ok := v.(string); ok {
		return x, nil
	}
	buff, err := json.Marshal(v)
	if err != nil {
		return "", err
	}

	switch en := v.(type) {
	case security.IEncrypt:
		return en.Encrypt(buff), nil
	default:
		return string(buff), nil
	}
}

func checkCover(v string) bool {
	y := "Yes,覆盖，使用当前配置覆盖已有配置"
	n := "No,跳过，不覆盖已有配置"
	prompt := promptui.Select{
		Label: fmt.Sprintf("注册中心已存在配置%s 是否覆盖?", v),
		Items: []string{y, n},
	}
	_, result, err := prompt.Run()
	return err == nil && result == y
}
