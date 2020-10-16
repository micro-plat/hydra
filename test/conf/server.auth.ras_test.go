/*
author:taoshouyin
time:2020-10-16
*/

package conf

import (
	"reflect"
	"testing"

	"github.com/micro-plat/hydra/conf"
	"github.com/micro-plat/hydra/conf/server/auth/ras"
)

func TestNewRASAuth(t *testing.T) {
	tests := []struct {
		name string
		args []*ras.Auth
		want *ras.RASAuth
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := ras.NewRASAuth(tt.args...); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("NewRASAuth() = %v, want %v", got, tt.want)
			}
		})
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
