package creator

import (
	"bytes"
	"fmt"
	"os"

	"github.com/BurntSushi/toml"
	"github.com/micro-plat/hydra/conf/server/api"
	"github.com/micro-plat/hydra/conf/server/cron"
	"github.com/micro-plat/hydra/global"
	"github.com/micro-plat/hydra/hydra/servers"
)

//IConf 配置注册管理
type IConf interface {

	//OnReady 系统准备好后触发
	OnReady(fs ...interface{})

	//API api服务器配置
	API(address string, opts ...api.Option) httpBuilder

	//Custome 自定义服务器配置
	Custome(tp string, s ...interface{}) customerBuilder
}

//Conf 配置服务
var Conf = &conf{data: make(map[string]map[string]interface{})}

type conf struct {
	funcs []func() error
	data  map[string]map[string]interface{}
}

//OnReady 注册配置准备函数
func (c *conf) OnReady(fs ...interface{}) {
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
	for _, f := range c.funcs {
		if err := f(); err != nil {
			return err
		}
	}
	types := servers.GetServerTypes()
	for _, t := range types {
		if _, ok := c.data[t]; !ok {
			switch t {
			case global.API:
				c.data[global.API] = newHTTP(":8080").load()
			case global.Web:
				c.data[global.Web] = newHTTP(":8089").load()
			default:
				c.data[t] = map[string]interface{}{
					"main": map[string]interface{}{},
				}
			}

		}
	}
	//添加其它服务器
	return nil
}

//API api服务器配置
func (c *conf) API(address string, opts ...api.Option) httpBuilder {
	api := newHTTP(address, opts...)
	c.data[global.API] = api
	return api
}

//CRON cron服务器配置
func (c *conf) CRON(opts ...cron.Option) cronBuilder {
	cron := newCron(opts...)
	c.data[global.CRON] = cron
	return cron
}

//Custome 用户自定义配置服务
func (c *conf) Custome(tp string, s ...interface{}) customerBuilder {
	if _, ok := c.data[tp]; ok {
		panic(fmt.Sprintf("不能重复注册%s", tp))
	}
	customer := newCustomerBuilder(s...)
	c.data[tp] = customer
	return customer
}

//Encode 将当前配置序列化为toml格式
func (c *conf) Encode() (string, error) {
	var buffer bytes.Buffer
	encoder := toml.NewEncoder(&buffer)
	err := encoder.Encode(&c.data)
	return buffer.String(), err
}

//Encode2File 将当前配置内容保存到文件中
func (c *conf) Encode2File(path string, cover bool) error {
	if !cover {
		if _, err := os.Stat(path); err == nil || os.IsExist(err) {
			return fmt.Errorf("配置文件已存在 %s", path)
		}
	}
	f, err := os.OpenFile(path, os.O_RDWR|os.O_CREATE, 0755)
	if err != nil {
		return fmt.Errorf("无法打开文件:%s %w", path, err)
	}
	encoder := toml.NewEncoder(f)
	err = encoder.Encode(&c.data)
	if err != nil {
		return err
	}
	if err := f.Close(); err != nil {
		return err
	}
	return nil
}

//Decode 从配置文件中读取配置信息
func (c *conf) Decode(f string) error {
	_, err := toml.DecodeFile(f, &c.data)
	return err
}
