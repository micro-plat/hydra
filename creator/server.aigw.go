package creator

import (
	"github.com/micro-plat/hydra/conf/server/acl/limiter"
	"github.com/micro-plat/hydra/conf/server/aigw"
	"github.com/micro-plat/hydra/conf/server/header"
	"github.com/micro-plat/hydra/conf/server/metric"
	"github.com/micro-plat/hydra/conf/server/processor"
)

type aigwBuilder struct {
	BaseBuilder
}

// newAIGW 构建AI网关配置生成器
func newAIGW(address string, opts ...aigw.Option) *aigwBuilder {
	b := &aigwBuilder{BaseBuilder: make(map[string]interface{})}
	b.BaseBuilder[ServerMainNodeName] = aigw.New(address, opts...)
	return b
}

// Load 加载路由
func (b *aigwBuilder) Load() {
}

// Header 响应头配置
func (b *aigwBuilder) Header(opts ...header.Option) *aigwBuilder {
	b.BaseBuilder[header.TypeNodeName] = header.New(opts...)
	return b
}

// Limit 服务器限流配置
func (b *aigwBuilder) Limit(opts ...limiter.Option) *aigwBuilder {
	path := limiter.ParNodeName + "/" + limiter.SubNodeName
	b.BaseBuilder[path] = limiter.New(opts...)
	return b
}

// Metric 服务监控配置
func (b *aigwBuilder) Metric(host string, db string, cron string, opts ...metric.Option) *aigwBuilder {
	b.BaseBuilder[metric.TypeNodeName] = metric.New(host, db, cron, opts...)
	return b
}

// Processor 构建Processor配置
func (b *aigwBuilder) Processor(opts ...processor.Option) *aigwBuilder {
	b.BaseBuilder[processor.TypeNodeName] = processor.New(opts...)
	return b
}
