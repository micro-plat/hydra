package limiter

import "golang.org/x/time/rate"

//Option 配置选项
type Option func(*Limiter)

// WithRuleList WithRuleList
func WithRuleList(list ...*Rule) Option {
	return func(a *Limiter) {
		for _, rule := range list {
			rule.limiter = rate.NewLimiter(rate.Limit(rule.MaxAllow), rule.MaxAllow)
			a.Rules = append(a.Rules, rule)
		}
	}
}

//WithDisable 关闭
func WithDisable() Option {
	return func(a *Limiter) {
		a.Disable = true
	}
}

//WithEnable 开启
func WithEnable() Option {
	return func(a *Limiter) {
		a.Disable = false
	}
}

//RuleOption Rule配置选项
type RuleOption func(*Rule)

//WithMaxWait 设置请求允许等待的时长
func WithMaxWait(second int) RuleOption {
	return func(a *Rule) {
		a.MaxWait = second
	}
}

//WithFallback 启用服务降级处理
func WithFallback() RuleOption {
	return func(a *Rule) {
		a.Fallback = true
	}
}

//WithReponse 设置响应内容
func WithReponse(status int, content string) RuleOption {
	return func(a *Rule) {
		a.Resp = &Resp{Status: status, Content: content}
	}
}
