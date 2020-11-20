package rpc

import (
	"fmt"
	"testing"

	"github.com/micro-plat/hydra/conf"
	"github.com/micro-plat/hydra/conf/app"
	"github.com/micro-plat/hydra/global"
	"github.com/micro-plat/hydra/hydra/servers/rpc"
	"github.com/micro-plat/hydra/registry/pub"
	"github.com/micro-plat/hydra/services"
	"github.com/micro-plat/hydra/test/assert"
	"github.com/micro-plat/hydra/test/mocks"
	"github.com/micro-plat/lib4go/logger"
)

func TestNewResponsive(t *testing.T) {
	confObj := mocks.NewConfBy("rpcserver_resserivece_test", "testrpacsdf") //构建对象
	confObj.RPC(":35211")
	tests := []struct {
		name    string
		cnf     app.IAPPConf
		proto   string
		addr    string
		wantErr bool
	}{
		{name: "构建rpc服务", addr: ":35211", proto: "tcp", cnf: confObj.GetRPCConf()},
	}
	for _, tt := range tests {
		gotH, err := rpc.NewResponsive(tt.cnf)
		assert.Equal(t, nil, err, tt.name)
		addr := fmt.Sprintf("%s://%s%s", tt.proto, global.LocalIP(), tt.addr)
		assert.Equal(t, addr, gotH.Server.GetAddress(), tt.name)
	}
}

func TestResponsive_Start(t *testing.T) {
	confObj := mocks.NewConfBy("rpcserver_resserivece_test1", "testrpacsdf1") //构建对象
	confObj.RPC(":35212")
	// global.Def.ServerTypes = []string{global.RPC}
	tests := []struct {
		name       string
		cnf        app.IAPPConf
		serverType string
		starting   func(app.IAPPConf) error
		wantErr    string
	}{
		{name: "starting报错", cnf: confObj.GetRPCConf(), serverType: global.RPC, starting: func(app.IAPPConf) error { return fmt.Errorf("err") }, wantErr: "err"},
		{name: "starting不报错", cnf: confObj.GetRPCConf(), serverType: global.RPC, starting: func(app.IAPPConf) error { return nil }, wantErr: "err"},
	}
	for _, tt := range tests {
		services.Def = services.New()
		services.Def.OnStarting(tt.starting, tt.serverType)
		w, err := rpc.NewResponsive(tt.cnf)
		assert.Equal(t, nil, err, tt.name)
		err = w.Start()
		if tt.wantErr != "" {
			assert.Equal(t, tt.wantErr, err.Error(), tt.name)
		}
	}
}

func TestResponsive_Notify(t *testing.T) {
	type fields struct {
		Server   *rpc.Server
		conf     app.IAPPConf
		comparer conf.IComparer
		pub      pub.IPublisher
		log      logger.ILogger
		first    bool
	}
	type args struct {
		c app.IAPPConf
	}
	tests := []struct {
		name       string
		fields     fields
		args       args
		wantChange bool
		wantErr    bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			w := &rpc.Responsive{
				Server: tt.fields.Server,
			}
			gotChange, err := w.Notify(tt.args.c)
			if (err != nil) != tt.wantErr {
				t.Errorf("Responsive.Notify() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if gotChange != tt.wantChange {
				t.Errorf("Responsive.Notify() = %v, want %v", gotChange, tt.wantChange)
			}
		})
	}
}
