package balancer

import "github.com/micro-plat/lib4go/metrics"
import "github.com/micro-plat/lib4go/concurrent/cmap"

//Limiter 限流组件,用于限制客户端每分钟请求服务器的数量,当超过限制数量时不再选择该服务提供者
type Limiter struct {
	settings       cmap.ConcurrentMap //限流规则配置
	currentService string             //当前服务名称
	metricRegistry metrics.Registry   //metric数据存储
}

//NewLimiter 创建限流组件
func NewLimiter(service string, lt map[string]int) *Limiter {
	m := &Limiter{currentService: service}
	m.metricRegistry = metrics.NewRegistry()
	m.settings = cmap.New(8)
	for k, v := range lt {
		m.settings.Set(k, float64(v))
	}
	return m
}

//Update 更新限流规则
func (m *Limiter) Update(lt map[string]int) {
	m.settings.Clear()
	for k, v := range lt {
		m.settings.Set(k, float64(v))
	}
}

//Check 检查当前服务IP是否达到限流条件
func (m *Limiter) Check(ip string) bool {
	if count, ok := m.settings.Get("*"); ok {
		limiterName := metrics.MakeName(".limiter", metrics.QPS, "service", m.currentService)
		meter := metrics.GetOrRegisterQPS(limiterName, m.metricRegistry)
		if meter.M1() >= count.(int32) {
			return false
		}
		meter.Mark(1)
	}
	if count, ok := m.settings.Get(m.currentService); ok {
		limiterName := metrics.MakeName(".limiter", metrics.QPS, "service", m.currentService, "ip", ip)
		meter := metrics.GetOrRegisterQPS(limiterName, m.metricRegistry)
		if meter.M1() >= count.(int32) {
			return false
		}
		meter.Mark(1)
	}
	return true
}
