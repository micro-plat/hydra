package builder

import (
	"bytes"

	"github.com/BurntSushi/toml"
	"github.com/micro-plat/hydra/registry/conf/server/api"
	"github.com/micro-plat/hydra/registry/conf/server/cron"
)

//Conf 配置服务
var Conf = &conf{data: make(map[string]map[string]interface{})}

type conf struct {
	funcs []func() error
	data  map[string]map[string]interface{}
}

//Ready 注册配置准备函数
func (c *conf) Ready(fs ...interface{}) {
	for _, fn := range fs {
		if f, ok := fn.(func()); ok {
			c.funcs = append(c.funcs, func() error {
				f()
				return nil
			})
			continue
		}
		if f, ok := fn.(func() error); ok {
			c.funcs = append(c.funcs, f)
			continue
		}
		panic("函数签名格式不正确，支持的格式有func()、func()error")
	}
}

//Load 加载所有配置
func (c *conf) Load() error {
	c.Decode("./conf.toml")
	for _, f := range c.funcs {
		if err := f(); err != nil {
			return err
		}
	}
	return nil
}

//API api服务器配置
func (c *conf) API(address string, opts ...api.Option) apiBuilder {
	api := NewAPI(address, opts...)
	c.data["api"] = api
	return api
}

//CRON cron服务器配置
func (c *conf) CRON(opts ...cron.Option) cronBuilder {
	cron := newCron(opts...)
	c.data["cron"] = cron
	return cron
}

//Encode 将当前配置序列化为toml格式
func (c *conf) Encode() (string, error) {
	var buffer bytes.Buffer
	encoder := toml.NewEncoder(&buffer)
	err := encoder.Encode(&c.data)
	return buffer.String(), err
}

//Decode 从配置文件中读取配置信息
func (c *conf) Decode(f string) error {
	_, err := toml.DecodeFile(f, &c.data)
	return err
}
