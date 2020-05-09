package builder

import (
	"encoding/json"
	"fmt"

	"github.com/micro-plat/hydra/registry"
	"github.com/micro-plat/hydra/registry/conf/server"
)

//Pub 将本地配置发布到配置中心
func (c *conf) Pub(platName string, systemName string, clusterName string, r registry.IRegistry) error {
	for tp, subs := range c.data {
		pub := server.NewPub(platName, systemName, tp, clusterName)
		if err := publish(r, pub.GetMainPath(), subs["main"]); err != nil {
			return err
		}
		for name, value := range subs {
			if name == "main" {
				continue
			}
			if err := publish(r, pub.GetSubConfPath(name), value); err != nil {
				return err
			}
		}
	}
	return nil
}

func publish(r registry.IRegistry, path string, v interface{}) error {
	value, err := getJSON(&v)
	if err != nil {
		return fmt.Errorf("将%s配置信息转化为json时出错:%w", path, err)
	}
	if err := r.CreatePersistentNode(path, value); err != nil {
		return fmt.Errorf("创建配置节点%s %s出错:%w", path, value, err)
	}
	return nil
}

//getJSON 将对象序列化为json字符串
func getJSON(v interface{}) (value string, err error) {
	buff, err := json.Marshal(&v)
	if err != nil {
		return "", err
	}
	return string(buff), nil
}
