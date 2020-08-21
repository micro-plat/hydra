package apm

import (
	"context"
	"fmt"
)

//IAPM 缓存接口
type IAPM interface {
	CreateTracer(service string) (tracer Tracer, err error)
}

//Resover 定义配置文件转换方法
type Resover interface {
	Resolve(instance, conf string) (IAPM, error)
}

var apmResolvers = make(map[string]Resover)

//Register 注册配置文件适配器
func Register(apmtype string, resolver Resover) {
	if resolver == nil {
		panic("apm: Register adapter is nil")
	}
	if _, ok := apmResolvers[apmtype]; ok {
		panic("apm: Register called twice for adapter " + apmtype)
	}
	apmResolvers[apmtype] = resolver
}

//New 根据适配器名称及参数返回配置处理器
func New(apmtype string, instance, conf string) (IAPM, error) {
	resolver, ok := apmResolvers[apmtype]
	if !ok {
		return nil, fmt.Errorf("cache: unknown adapter name %q (forgotten import?)", apmtype)
	}
	return resolver.Resolve(instance, conf)
}

type APMInfo struct {
	Tracer  Tracer
	RootCtx context.Context
}
