/*
author:taoshouyin
time:2020-10-16
*/

package conf

import (
	"reflect"
	"testing"

	"github.com/micro-plat/hydra/test/assert"

	"github.com/micro-plat/hydra/conf"
	"github.com/micro-plat/hydra/conf/server/auth/ras"
)

func TestNewAuth(t *testing.T) {
	tests := []struct {
		name    string
		service string
		opts    []ras.Option
		want    *ras.Auth
	}{
		{name: "设置默认对象", service: "", opts: []ras.Option{}, want: &ras.Auth{Service: "", Requests: []string{"*"}, Connect: &ras.Connect{},
			Params: make(map[string]interface{}), Required: make([]string, 0, 1), Alias: make(map[string]string), Decrypt: make([]string, 0, 1)}},
		{name: "设置service对象", service: "test-tsy", opts: []ras.Option{}, want: &ras.Auth{Service: "test-tsy", Requests: []string{"*"}, Connect: &ras.Connect{},
			Params: make(map[string]interface{}), Required: make([]string, 0, 1), Alias: make(map[string]string), Decrypt: make([]string, 0, 1)}},
		{name: "设置Requests对象", service: "", opts: []ras.Option{ras.WithRequest("/t1/t2")}, want: &ras.Auth{Service: "", Requests: []string{"/t1/t2"}, Connect: &ras.Connect{},
			Params: make(map[string]interface{}), Required: make([]string, 0, 1), Alias: make(map[string]string), Decrypt: make([]string, 0, 1)}},
	}
	for _, tt := range tests {
		got := ras.New(tt.service, tt.opts...)
		assert.Equal(t, tt.want, got, tt.name)
	}
}

func TestNewRASAuth(t *testing.T) {
	tests := []struct {
		name    string
		service string
		opts    []ras.Option
		want    *ras.Auth
	}{
		{name: "设置默认对象", service: "", opts: []ras.Option{}, want: &ras.Auth{Service: "", Requests: []string{"*"}, Connect: &ras.Connect{},
			Params: make(map[string]interface{}), Required: make([]string, 0, 1), Alias: make(map[string]string), Decrypt: make([]string, 0, 1)}},
	}
	for _, tt := range tests {
		got := ras.New(tt.service, tt.opts...)
		assert.Equal(t, tt.want, got, tt.name)
	}
}

func TestRASAuth_Match(t *testing.T) {
	tests := []struct {
		name    string
		disable bool
		auth    []*ras.Auth
		args    string
		want    bool
		want1   *ras.Auth
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			a := ras.RASAuth{
				Disable: tt.disable,
				Auth:    tt.auth,
			}
			got, got1 := a.Match(tt.args)
			if got != tt.want {
				t.Errorf("RASAuth.Match() got = %v, want %v", got, tt.want)
			}
			if !reflect.DeepEqual(got1, tt.want1) {
				t.Errorf("RASAuth.Match() got1 = %v, want %v", got1, tt.want1)
			}
		})
	}
}

func TestAuthRASGetConf(t *testing.T) {
	tests := []struct {
		name      string
		args      conf.IMainConf
		wantAuths *ras.RASAuth
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if gotAuths := ras.GetConf(tt.args); !reflect.DeepEqual(gotAuths, tt.wantAuths) {
				t.Errorf("GetConf() = %v, want %v", gotAuths, tt.wantAuths)
			}
		})
	}
}
