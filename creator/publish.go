package creator

import (
	"encoding/json"
	"fmt"
	"path/filepath"

	"github.com/micro-plat/hydra/conf/server"
	varpub "github.com/micro-plat/hydra/conf/vars"
	"github.com/micro-plat/hydra/global"
	"github.com/micro-plat/hydra/registry"
)

//Pub 将配置发布到配置中心
func (c *conf) Pub(platName string, systemName string, clusterName string, registryAddr string, cover bool) error {
	if err := c.Load(); err != nil {
		return err
	}

	//本地文件系统则直接使用toml序列化方式进行发布
	proto := registry.GetProto(registryAddr)
	if proto == registry.FileSystem {
		return c.Encode2File(filepath.Join(registry.GetAddrs(registryAddr)[0], global.Def.LocalConfName), cover)
	}

	//创建注册中心，根据注册中心提供的接口进行配置发布
	r, err := registry.NewRegistry(registryAddr, global.Def.Log())
	if err != nil {
		return err
	}
	for tp, subs := range c.data {
		pub := server.NewServerPub(platName, systemName, tp, clusterName)
		if err := publish(r, pub.GetServerPath(), subs.Map()[ServerMainNodeName], cover); err != nil {
			return err
		}
		for name, value := range subs.Map() {
			if name == ServerMainNodeName {
				continue
			}
			if err := publish(r, pub.GetSubConfPath(name), value, cover); err != nil {
				return err
			}
		}
	}
	for tp, subs := range c.vars {
		pub := varpub.NewVarPub(platName)
		for k, v := range subs {
			if err := publish(r, pub.GetVarPath(tp, k), v, cover); err != nil {
				return err
			}
		}

	}
	return nil
}

func publish(r registry.IRegistry, path string, v interface{}, cover bool) error {

	value, err := getJSON(v)
	if err != nil {
		return fmt.Errorf("将%s配置信息转化为json时出错:%w", path, err)
	}
	if !cover {
		if b, _ := r.Exists(path); b {
			return nil
		}
	}
	if err := deleteAll(r, path); err != nil {
		return err
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
func getJSON(v interface{}) (value string, err error) {
	if x, ok := v.(string); ok {
		return x, nil
	}

	buff, err := json.Marshal(&v)
	if err != nil {
		return "", err
	}
	return string(buff), nil
}
