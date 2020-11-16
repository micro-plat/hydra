package middleware

import (
	"fmt"
	"sync"

	"github.com/micro-plat/hydra/components/pkgs/metrics"
	"github.com/micro-plat/lib4go/logger"
	"github.com/micro-plat/lib4go/net"
)

//Metric 服务器处理能力统计
type Metric struct {
	reporter        metrics.IReporter
	logger          *logger.Logger
	currentRegistry metrics.Registry
	needCollect     bool
	once            sync.Once
	ip              string
}

//NewMetric new metric
func NewMetric() *Metric {
	return &Metric{}

}
func (m *Metric) onceDo(ctx IMiddleContext) {
	m.once.Do(func() {
		metric, err := ctx.APPConf().GetMetricConf()
		if err != nil {
			panic(fmt.Errorf("metric配置获取失败:%w", err))
		}

		if metric.Disable {
			return
		}

		m.currentRegistry = metrics.NewRegistry()
		m.ip = net.GetLocalIPAddress()
		m.logger = logger.New("metric")

		//2. 创建上报服务
		m.reporter, err = metrics.InfluxDB(m.currentRegistry,
			metric.Cron,
			metric.Host,
			metric.DataBase,
			metric.UserName,
			metric.Password, m.logger)
		if err != nil {
			panic(fmt.Errorf("初始化metric失败:%w", err))
		}
		m.needCollect = true
		//定时上报
		go m.reporter.Run()

	})
}

//Handle 处理请求
func (m *Metric) Handle() Handler {
	return func(ctx IMiddleContext) {

		//执行首次初始化
		m.onceDo(ctx)
		if !m.needCollect {
			ctx.Next()
			return
		}

		ctx.Response().AddSpecial("metric")

		//1. 初始化三类统计器---请求的QPS/正在处理的计数器/时间统计器
		url := ctx.Request().Path().GetRequestPath()
		conterName := metrics.MakeName(ctx.APPConf().GetServerConf().GetServerType()+".server.request", metrics.WORKING, "server", ctx.APPConf().GetServerConf().GetServerName(), "host", m.ip, "url", url) //堵塞计数
		timerName := metrics.MakeName(ctx.APPConf().GetServerConf().GetServerType()+".server.request", metrics.TIMER, "server", ctx.APPConf().GetServerConf().GetServerName(), "host", m.ip, "url", url)    //堵塞计数
		requestName := metrics.MakeName(ctx.APPConf().GetServerConf().GetServerType()+".server.request", metrics.QPS, "server", ctx.APPConf().GetServerConf().GetServerName(), "host", m.ip, "url", url)    //请求数

		//2. 对QPS进行计数
		metrics.GetOrRegisterQPS(requestName, m.currentRegistry).Mark(1)

		//3.对正在请求的服务进行计数
		counter := metrics.GetOrRegisterCounter(conterName, m.currentRegistry)
		counter.Inc(1)

		//4. 对服务处理时长进行统计
		metrics.GetOrRegisterTimer(timerName, m.currentRegistry).Time(func() {
			ctx.Next()
		})

		//5. 服务处理完成后进行减数
		counter.Dec(1)

		//6. 初始化第四类统计器----状态码上报
		statusCode, _ := ctx.Response().GetRawResponse()
		responseName := metrics.MakeName(ctx.APPConf().GetServerConf().GetServerType()+".server.response", metrics.METER, "server", ctx.APPConf().GetServerConf().GetServerName(), "host", m.ip,
			"url", url, "status", fmt.Sprintf("%d", statusCode)) //完成数

		//7. 对服务处理结果的状态码进行上报
		metrics.GetOrRegisterMeter(responseName, m.currentRegistry).Mark(1)
	}

}

//Stop stop metric
func (m *Metric) Stop() {
	if m.reporter != nil {
		m.reporter.Close()
		m.reporter = nil
	}
}
