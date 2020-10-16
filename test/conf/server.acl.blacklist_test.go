/*
author:taoshouyin
time:2020-10-15
*/

package conf

import (
	"fmt"
	"strings"
	"testing"
	"time"

	"github.com/micro-plat/hydra/conf/server/acl/blacklist"
	"github.com/micro-plat/hydra/conf/server/api"
	_ "github.com/micro-plat/hydra/registry/registry/localmemory"
	"github.com/micro-plat/hydra/test/assert"
	"github.com/micro-plat/hydra/test/mocks"
)

func TestBlackListGetConf(t *testing.T) {
	conf := mocks.NewConf() //构建对象
	confB := conf.API(":8081", api.WithHeaderReadTimeout(30), api.WithTimeout(30, 30))
	bobj := blacklist.GetConf(conf.GetAPIConf().GetMainConf())
	assert.Equal(t, bobj, &blacklist.BlackList{Disable: true}, "节点不存在,获取默认对象")

	confB.BlackList(blacklist.WithIP("192.168.0.*", "192.168.1.2"))
	bobj = blacklist.GetConf(conf.GetAPIConf().GetMainConf())
	assert.Equal(t, bobj, blacklist.New(blacklist.WithIP("192.168.0.*", "192.168.1.2")), "正常对象获取")

	//获取时json对象不合法
	path := conf.GetAPIConf().GetMainConf().GetSubConfPath("acl", "black.list")
	defer func() {
		if e := recover(); e != nil {
			if !strings.Contains(e.(string), fmt.Sprintf("获取%s配置失败", path)) {
				t.Error("json错误,返回了未知的错误信息")
			}
		}
	}()
	conf.Registry.Update(path, "错误的json字符串")
	ch, _ := conf.Registry.WatchValue(path)
	select {
	case <-time.After(3 * time.Second):
		return
	case <-ch:
		bobj = blacklist.GetConf(conf.GetAPIConf().GetMainConf())
		t.Errorf("%v", bobj)
	}
}
