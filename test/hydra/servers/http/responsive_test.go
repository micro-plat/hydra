package http

import (
	"fmt"
	"testing"

	"github.com/micro-plat/hydra/conf/app"
	"github.com/micro-plat/hydra/global"
	"github.com/micro-plat/hydra/hydra/servers/http"
	"github.com/micro-plat/hydra/services"
	"github.com/micro-plat/hydra/test/assert"
	"github.com/micro-plat/hydra/test/mocks"
)

func TestNewResponsive(t *testing.T) {
	confObj := mocks.NewConf() //构建对象
	confObj.API(":55003")      //初始化参数
	confObj.Web(":55002")      //初始化参数
	confObj.WS(":55001")       //初始化参数
	tests := []struct {
		name    string
		proto   string
		addr    string
		cnf     app.IAPPConf
		wantErr bool
	}{
		{name: "构建ws服务", addr: ":55001", proto: "ws", cnf: confObj.GetWSConf()},
		{name: "构建web服务", addr: ":55002", proto: "http", cnf: confObj.GetWebConf()},
		{name: "构建api服务", addr: ":55003", proto: "http", cnf: confObj.GetAPIConf()},
	}
	for _, tt := range tests {
		gotH, err := http.NewResponsive(tt.cnf)
		assert.Equal(t, nil, err, tt.name)
		addr := fmt.Sprintf("%s://%s%s", tt.proto, global.LocalIP(), tt.addr)
		assert.Equal(t, addr, gotH.Server.GetAddress(), tt.name)
	}
}

func xxTestResponsive_Start(t *testing.T) {
	confObj := mocks.NewConf() //构建对象
	confObj.API(":55004")      //初始化参数
	tests := []struct {
		name       string
		cnf        app.IAPPConf
		serverType string
		starting   func(app.IAPPConf) error
		wantErr    string
	}{
		{name: "starting报错", cnf: confObj.GetAPIConf(), serverType: "api", starting: func(app.IAPPConf) error { return fmt.Errorf("err") }, wantErr: "err"},
		{name: "starting不报错", cnf: confObj.GetAPIConf(), serverType: "api", starting: func(app.IAPPConf) error { return nil }, wantErr: "err"},
	}
	for _, tt := range tests {
		services.Def = services.New()
		services.Def.OnStarting(tt.starting)
		w, err := http.NewResponsive(tt.cnf)
		assert.Equal(t, nil, err, tt.name)
		err = w.Start()
		if tt.wantErr != "" {
			assert.Equal(t, tt.wantErr, err.Error(), tt.name)
		}
	}
}
