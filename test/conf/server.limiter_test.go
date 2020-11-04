package conf

import (
	"testing"

	"github.com/micro-plat/hydra/conf"
	"github.com/micro-plat/hydra/conf/server/limiter"
	"github.com/micro-plat/hydra/test/assert"
	"github.com/micro-plat/hydra/test/mocks"
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
		opts []limiter.Option
		want *limiter.Limiter
	}{
		{name: "初始化空对象", opts: []limiter.Option{}, want: &limiter.Limiter{Rules: []*limiter.Rule{}}},
		{name: "初始化disable对象", opts: []limiter.Option{limiter.WithDisable()}, want: &limiter.Limiter{Disable: true, Rules: []*limiter.Rule{}}},
		{name: "初始化enable对象", opts: []limiter.Option{limiter.WithEnable()}, want: &limiter.Limiter{Disable: false, Rules: []*limiter.Rule{}}},
		{name: "初始化Rules对象", opts: []limiter.Option{limiter.WithRuleList(limiter.NewRule("path", 1, limiter.WithMaxWait(3)))},
			want: &limiter.Limiter{Rules: []*limiter.Rule{limiter.NewRule("path", 1, limiter.WithMaxWait(3))}}},
	}
	for _, tt := range tests {
		got := limiter.New(tt.opts...)
		assert.Equal(t, tt.want.Disable, got.Disable, tt.name+",disable")
		assert.Equal(t, tt.want.Rules, got.Rules, tt.name+",Rules")
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
		{name: "空对象匹配", lmt: limiter.New(), path: "xx", want: false, want1: nil},
		{name: "精确对象匹配正确", lmt: limiter.New(limiter.WithRuleList(limiter.NewRule("/t1/t2", 1, limiter.WithMaxWait(3)))),
			path: "/t1/t2", want: true, want1: limiter.NewRule("/t1/t2", 1, limiter.WithMaxWait(3))},
		{name: "精确对象匹配失败", lmt: limiter.New(limiter.WithRuleList(limiter.NewRule("/t1/t2", 1, limiter.WithMaxWait(3)))), path: "/t1", want: false, want1: nil},
		{name: "单模糊对象匹配正确", lmt: limiter.New(limiter.WithRuleList(limiter.NewRule("/t1/*", 1, limiter.WithMaxWait(3)))), path: "/t1/tt", want: true, want1: limiter.NewRule("/t1/*", 1, limiter.WithMaxWait(3))},
		{name: "单模糊对象匹配失败", lmt: limiter.New(limiter.WithRuleList(limiter.NewRule("/t1/*", 1, limiter.WithMaxWait(3)))), path: "/t1/tt/ww", want: false, want1: nil},
		{name: "多模糊象匹配正确", lmt: limiter.New(limiter.WithRuleList(limiter.NewRule("/t1/t2/**", 1, limiter.WithMaxWait(3)))), path: "/t1/t2/ww/ss", want: true, want1: limiter.NewRule("/t1/t2/**", 1, limiter.WithMaxWait(3))},
		{name: "多模糊象匹配失败", lmt: limiter.New(limiter.WithRuleList(limiter.NewRule("/t1/t2/**", 1, limiter.WithMaxWait(3)))), path: "/t1/tt/ww/ss", want: false, want1: nil},
	}
	for _, tt := range tests {
		got, got1 := tt.lmt.GetLimiter(tt.path)
		assert.Equal(t, tt.want, got, tt.name+",bool")
		assert.Equal(t, tt.want1, got1, tt.name+",Path")
	}

	defer func() {
		e := recover()
		if e != nil {
			assert.Equal(t, "从缓存中未找到limite组件", e.(error).Error(), "从缓存中未找到limite组件")
		}
	}()

	//limiter是私有变脸不能测试分支panic
	tests = []test{}
}

func TestGetConf(t *testing.T) {
	type test struct {
		name    string
		cnf     conf.IServerConf
		want    *limiter.Limiter
		wantErr bool
	}
	conf := mocks.NewConfBy("hydra", "graytest")
	confB := conf.API(":8090")
	test1 := test{name: "限流节点不存在", cnf: conf.GetAPIConf().GetServerConf(), want: &limiter.Limiter{Disable: true}, wantErr: false}
	limiterObj, err := limiter.GetConf(test1.cnf)
	assert.Equal(t, test1.wantErr, (err != nil), test1.name)
	assert.Equal(t, test1.want, limiterObj, test1.name)

	confB.Limit(limiter.WithDisable())
	test2 := test{name: "限流节点存在,auths不存在", cnf: conf.GetAPIConf().GetServerConf(), want: &limiter.Limiter{Disable: true}, wantErr: false}
	limiterObj, err = limiter.GetConf(test2.cnf)
	assert.Equal(t, test2.wantErr, (err != nil), test2.name+",err")
	assert.Equal(t, test2.want, limiterObj, test2.name+",obj")

	confB.Limit(limiter.WithRuleList(limiter.NewRule("path", 1, limiter.WithAction("sss"))))
	test3 := test{name: "灰度节点存在,数据不合法", cnf: conf.GetAPIConf().GetServerConf(), want: nil, wantErr: true}
	limiterObj, err = limiter.GetConf(test3.cnf)
	assert.Equal(t, test3.wantErr, (err != nil), test3.name+",err")
	assert.Equal(t, test3.want, limiterObj, test3.name+",obj")

	confB.Limit(limiter.WithRuleList(limiter.NewRule("/path", 1, limiter.WithFallback(), limiter.WithMaxWait(3), limiter.WithAction("GET"), limiter.WithReponse(200, "success"))))
	test4 := test{name: "灰度节点存在,正确配置",
		cnf: conf.GetAPIConf().GetServerConf(), want: limiter.New(limiter.WithRuleList(limiter.NewRule("/path", 1, limiter.WithFallback(), limiter.WithMaxWait(3), limiter.WithAction("GET"), limiter.WithReponse(200, "success")))), wantErr: false}
	limiterObj, err = limiter.GetConf(test4.cnf)
	assert.Equal(t, test4.wantErr, (err != nil), test4.name+",err")
	assert.Equal(t, test4.want, limiterObj, test4.name+",obj")
}
