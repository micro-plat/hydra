package conf

import (
	"fmt"
	"reflect"
	"testing"

	"github.com/micro-plat/hydra/conf"
	"github.com/micro-plat/hydra/conf/server/limiter"
	"github.com/micro-plat/hydra/test/assert"
)

func TestNewLimmiterRule(t *testing.T) {
	type args struct {
		path  string
		allow int
		opts  []limiter.RuleOption
	}
	tests := []struct {
		name string
		args args
		want *limiter.Rule
	}{
		{name: "默认对象初始化", args: args{path: "path", allow: 1, opts: []limiter.RuleOption{}}, want: &limiter.Rule{Path: "path", MaxAllow: 1}},
		{name: "设置action对象初始化", args: args{path: "path", allow: 1, opts: []limiter.RuleOption{limiter.WithAction("get", "post")}}, want: &limiter.Rule{Action: []string{"get", "post"}, Path: "path", MaxAllow: 1}},
		{name: "设置maxwait对象初始化", args: args{path: "path", allow: 1, opts: []limiter.RuleOption{limiter.WithMaxWait(3)}}, want: &limiter.Rule{MaxWait: 3, Path: "path", MaxAllow: 1}},
		{name: "设置fallback对象初始化", args: args{path: "path", allow: 1, opts: []limiter.RuleOption{limiter.WithFallback()}}, want: &limiter.Rule{Fallback: true, Path: "path", MaxAllow: 1}},
		{name: "设置Resp对象初始化", args: args{path: "path", allow: 1, opts: []limiter.RuleOption{limiter.WithReponse(205, "success")}}, want: &limiter.Rule{Resp: &limiter.Resp{Status: 205, Content: "success"}, Path: "path", MaxAllow: 1}},
	}
	for _, tt := range tests {
		got := limiter.NewRule(tt.args.path, tt.args.allow, tt.args.opts...)
		assert.Equal(t, tt.want.Path, got.Path, tt.name+",Path")
		assert.Equal(t, tt.want.Action, got.Action, tt.name+",Action")
		assert.Equal(t, tt.want.MaxAllow, got.MaxAllow, tt.name+",MaxAllow")
		assert.Equal(t, tt.want.MaxWait, got.MaxWait, tt.name+",MaxWait")
		assert.Equal(t, tt.want.Fallback, got.Fallback, tt.name+",Fallback")
		assert.Equal(t, tt.want.Resp, got.Resp, tt.name+",Resp")
	}
}

func TestLimiterNew(t *testing.T) {
	tests := []struct {
		name string
		auth *limiter.Rule
		opts []limiter.Option
		want *limiter.Limiter
	}{
		{name: "初始化空对象", auth: &limiter.Rule{}, opts: []limiter.Option{}, want: &limiter.Limiter{Rules: []*limiter.Rule{&limiter.Rule{}}}},
		{name: "初始化disable对象", auth: &limiter.Rule{}, opts: []limiter.Option{limiter.WithDisable()}, want: &limiter.Limiter{Disable: true, Rules: []*limiter.Rule{&limiter.Rule{}}}},
		{name: "初始化enable对象", auth: &limiter.Rule{}, opts: []limiter.Option{limiter.WithEnable()}, want: &limiter.Limiter{Disable: false, Rules: []*limiter.Rule{&limiter.Rule{}}}},
		{name: "初始化Rules对象", auth: limiter.NewRule("path", 1, limiter.WithMaxWait(3)),
			opts: []limiter.Option{limiter.WithRuleList(limiter.NewRule("path", 1, limiter.WithFallback()))},
			want: &limiter.Limiter{Rules: []*limiter.Rule{limiter.NewRule("path", 1, limiter.WithMaxWait(3)), limiter.NewRule("path", 1, limiter.WithFallback())}}},
	}
	for _, tt := range tests {
		got := limiter.New(tt.auth, tt.opts...)
		assert.Equal(t, tt.want.Disable, got.Disable, tt.name+",disable")
		assert.Equal(t, tt.want.Rules, got.Rules, tt.name+",Rules")
	}
}

func TestGetConf(t *testing.T) {
	type args struct {
		cnf conf.IMainConf
	}
	tests := []struct {
		name    string
		args    args
		want    *limiter.Limiter
		wantErr bool
	}{
		// TODO: Add test cases.

	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := limiter.GetConf(tt.args.cnf)
			if (err != nil) != tt.wantErr {
				t.Errorf("GetConf() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("GetConf() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestLimiter_GetLimiter(t *testing.T) {
	type test struct {
		name  string
		lmt   *limiter.Limiter
		path  string
		want  bool
		want1 *limiter.Rule
	}
	tests := []test{
		//{name: "空对象匹配", lmt: limiter.New(&limiter.Rule{}), path: "xx", want: false, want1: nil},
		{name: "精确对象匹配正确", lmt: limiter.New(limiter.NewRule("/t1/t2", 1, limiter.WithMaxWait(3))), path: "/t1/t2", want: true, want1: limiter.NewRule("/t1/t2", 1, limiter.WithMaxWait(3))},
		//{name: "精确对象匹配失败", lmt: limiter.New(limiter.NewRule("/t1/t2", 1, limiter.WithMaxWait(3))), path: "/t1", want: false, want1: nil},
		//{name: "单模糊对象匹配正确", lmt: limiter.New(limiter.NewRule("/t1/*", 1, limiter.WithMaxWait(3))), path: "/t1/tt", want: true, want1: limiter.NewRule("/t1/*", 1, limiter.WithMaxWait(3))},
		//{name: "单模糊对象匹配失败", lmt: limiter.New(limiter.NewRule("/t1/*", 1, limiter.WithMaxWait(3))), path: "/t1/tt/ww", want: false, want1: nil},
		//{name: "多模糊象匹配正确", lmt: limiter.New(limiter.NewRule("/t1/t2/**", 1, limiter.WithMaxWait(3))), path: "/t1/t2/ww/ss", want: true, want1: limiter.NewRule("/t1/t2/**", 1, limiter.WithMaxWait(3))},
		//{name: "多模糊象匹配失败", lmt: limiter.New(limiter.NewRule("/t1/t2/**", 1, limiter.WithMaxWait(3))), path: "/t1/tt/ww/ss", want: false, want1: nil},
	}
	for _, tt := range tests {
		fmt.Println(tt.name)
		got, got1 := tt.lmt.GetLimiter(tt.path)
		assert.Equal(t, tt.want, got, tt.name+",bool")
		assert.Equal(t, tt.want1, got1, tt.name+",Path")
		// assert.Equal(t, tt.want1.Path, got1.Path, tt.name+",Path")
		// assert.Equal(t, tt.want1.Action, got1.Action, tt.name+",Action")
		// assert.Equal(t, tt.want1.MaxAllow, got1.MaxAllow, tt.name+",MaxAllow")
		// assert.Equal(t, tt.want1.MaxWait, got1.MaxWait, tt.name+",MaxWait")
		// assert.Equal(t, tt.want1.Fallback, got1.Fallback, tt.name+",Fallback")
		// assert.Equal(t, tt.want1.Resp, got1.Resp, tt.name+",Resp")
	}

	defer func() {
		e := recover()
		assert.Equal(t, "从缓存中未找到limite组件", e.(error).Error(), "从缓存中未找到limite组件")
	}()
	tests = []test{}
}
