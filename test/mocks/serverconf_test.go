package mocks

import (
	"testing"
)

func TestGetConf(t *testing.T) {
	conf := NewConf() //构建对象

	conf.API(":8080") //初始化参数

	server := conf.GetAPIConf() //获取配置

	if server.GetMainConf().GetRootConf().GetString("address") != ":8080" {
		t.Error("端口号获取失败")
	}

}
