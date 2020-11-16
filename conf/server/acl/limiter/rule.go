package limiter

import (
	"net/http"
	"time"

	"golang.org/x/time/rate"
)

//Resp 限流器
type Resp struct {
	Status  int    `json:"status" valid:"required" toml:"status,omitempty"`
	Content string `json:"content" valid:"required" toml:"content,omitempty"`
}

//Rule 按请求设定的限流器
type Rule struct {
	Path     string `json:"path" valid:"ascii,required" toml:"path,omitempty"`
	MaxAllow int    `json:"maxAllow"  toml:"maxAllow,omitempty"`
	MaxWait  int    `json:"maxWait,omitempty"  toml:"maxWait,omitempty"`
	Fallback bool   `json:"fallback,omitempty"  toml:"fallback,omitempty"`
	Resp     *Resp  `json:"resp,omitempty" valid:"required" toml:"resp,omitempty"`
	limiter  *rate.Limiter
}

//NewRule 构建限流规则
func NewRule(path string, allow int, opts ...RuleOption) *Rule {
	r := &Rule{
		Path:     path,
		MaxAllow: allow,
	}
	for _, opt := range opts {
		opt(r)
	}

	r.limiter = rate.NewLimiter(rate.Limit(r.MaxAllow), r.MaxAllow)
	return r
}

//GetLimiter 获取限流器
func (l *Rule) GetLimiter() *rate.Limiter {
	return l.limiter
}

//GetDelay 获取延迟等待时长
func (l *Rule) GetDelay() time.Duration {
	return time.Second * time.Duration(l.MaxWait)
}

//GetResponse 获取响应信息
func (l *Rule) GetResponse() (int, string) {
	if l.Resp == nil {
		return http.StatusTooManyRequests, ""
	}
	return l.Resp.Status, l.Resp.Content
}
