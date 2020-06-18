package global

import "fmt"

//MQConf 消息队列全局配置
var MQConf = &messageQueueConf{
	platNameAsPrefix: true,
	separate:         ":",
}

type messageQueueConf struct {
	platNameAsPrefix bool
	separate         string
}

//PlatNameAsPrefix 平台名称作为前缀
func (m *messageQueueConf) PlatNameAsPrefix(v bool) {
	m.platNameAsPrefix = v
}

//Separate 队列分段分隔符
func (m *messageQueueConf) Separate(v string) {
	m.separate = v
}

//NeedAddPrefix 是否需要增加前缀
func (m *messageQueueConf) NeedAddPrefix() bool {
	return m.platNameAsPrefix
}

//GetQueueName 获取队列名称
func (m *messageQueueConf) GetQueueName(n string) string {
	if m.platNameAsPrefix {
		return fmt.Sprintf("%s%s%s", Def.PlatName, m.separate, n)
	}
	return n
}
