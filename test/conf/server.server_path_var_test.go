package conf

import (
	"testing"

	"github.com/micro-plat/hydra/conf/vars"
	"github.com/micro-plat/hydra/test/assert"
)

func TestNewVarPub(t *testing.T) {
	tests := []struct {
		name     string
		platName string
		want     string
	}{
		{name: "1. Conf-NewVarPub-初始化平台名为空", platName: "", want: ""},
		{name: "2. Conf-NewVarPub-初始化平台名1", platName: "nil", want: "nil"},
		{name: "3. Conf-NewVarPub-初始化平台名2", platName: "paltname", want: "paltname"},
	}
	for _, tt := range tests {
		got := vars.NewVarPub(tt.platName)
		assert.Equal(t, tt.want, got.GetPlatName(), tt.name)
	}
}

func Test_varPath_GetVarPath(t *testing.T) {
	tests := []struct {
		name     string
		platName string
		tp       []string
		want     string
	}{
		{name: "1. Conf-varPathGetVarPath-入参为nil", platName: "pt1", tp: nil, want: "/pt1/var"},
		{name: "2. Conf-varPathGetVarPath-入参为空", platName: "pt1", tp: []string{}, want: "/pt1/var"},
		{name: "3. Conf-varPathGetVarPath-入参一段", platName: "pt1", tp: []string{"p1"}, want: "/pt1/var/p1"},
		{name: "4. Conf-varPathGetVarPath-入参二段", platName: "pt1", tp: []string{"p1", "p2"}, want: "/pt1/var/p1/p2"},
		{name: "5. Conf-varPathGetVarPath-入参三段", platName: "pt1", tp: []string{"p1", "p2", "p3"}, want: "/pt1/var/p1/p2/p3"},
		{name: "6. Conf-varPathGetVarPath-入参六段", platName: "pt1", tp: []string{"p1", "p2", "p3", "p4", "p5", "p6"}, want: "/pt1/var/p1/p2/p3/p4/p5/p6"},
	}
	for _, tt := range tests {
		c := vars.NewVarPub(tt.platName)
		got := c.GetVarPath(tt.tp...)
		assert.Equal(t, tt.want, got, tt.name)
	}
}

func Test_varPath_GetRLogPath(t *testing.T) {
	tests := []struct {
		name     string
		platName string
		want     string
	}{
		{name: "1. Conf-varPathGetRLogPath-rlog路径获取1", platName: "pt1", want: "/pt1/var/app/rlog"},
		{name: "2. Conf-varPathGetRLogPath-rlog路径获取2", platName: "pt2", want: "/pt2/var/app/rlog"},
		{name: "3. Conf-varPathGetRLogPath-rlog路径获取3", platName: "pt3", want: "/pt3/var/app/rlog"},
	}
	for _, tt := range tests {
		c := vars.NewVarPub(tt.platName)
		got := c.GetRLogPath()
		assert.Equal(t, tt.want, got, tt.name)
	}
}
