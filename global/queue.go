package global

import (
	"fmt"
	"strings"
)

// MQConf 消息队列全局配置
var MQConf = &messageQueueConf{
	platNameAsPrefix: true,
	encodingSend:     true,
	separate:         ":",
}

type messageQueueConf struct {
	platNameAsPrefix bool
	separate         string
	encodingSend     bool
}

// PlatNameAsPrefix 平台名称作为前缀
func (m *messageQueueConf) PlatNameAsPrefix(v bool) {
	m.platNameAsPrefix = v
}
func (m *messageQueueConf) EncodingSend(v bool) {
	m.encodingSend = v
}

// Separate 队列分段分隔符
func (m *messageQueueConf) Separate(v string) {
	m.separate = v
}

// NeedAddPrefix 是否需要增加前缀
func (m *messageQueueConf) NeedAddPrefix() bool {
	return m.platNameAsPrefix
}
func (m *messageQueueConf) IsEncodingSend() bool {
	return m.encodingSend
}

// GetQueueName 获取带有平台名称的完整队列名
func (m *messageQueueConf) GetQueueName(n string) string {
	if m.platNameAsPrefix {
		return fmt.Sprintf("%s%s%s", Def.PlatName, m.separate, n)
	}
	return n
}

// GetOriginalName 获取队列的原始队列名
func (m *messageQueueConf) GetOriginalName(n string) string {
	if m.platNameAsPrefix {
		n = strings.TrimPrefix(n, Def.PlatName+m.separate)
		//return fmt.Sprintf("%s%s%s", Def.PlatName, m.separate, n)
	}
	return n
}
