package middleware

import (
	"fmt"

	"github.com/micro-plat/hydra/registry/conf/server/metric"
	"github.com/micro-plat/hydra/servers/pkg/swap"
	"github.com/micro-plat/lib4go/logger"
	"github.com/micro-plat/lib4go/metrics"
	"github.com/micro-plat/lib4go/net"
)

//Metric 服务器处理能力统计
type Metric struct {
	reporter        metrics.IReporter
	logger          *logger.Logger
	currentRegistry metrics.Registry
	name            string
	serverType      string
	ip              string
}

//NewMetric new metric
func NewMetric(name string, serverType string, f metric.IMetric) (*Metric, error) {

	//1. 检查metric配置
	metric, ok := f.GetConf()
	if !ok || metric.Disable {
		return &Metric{}, nil
	}

	//2.创建配置信息
	m := &Metric{
		name:            name,
		serverType:      serverType,
		currentRegistry: metrics.NewRegistry(),
		ip:              net.GetLocalIPAddress(),
		logger:          logger.New("metric"),
	}

	//3. 创建上报服务
	var err error
	m.reporter, err = metrics.InfluxDB(m.currentRegistry,
		metric.Cron,
		metric.Host, metric.DataBase,
		metric.UserName,
		metric.Password, m.logger)
	if err != nil {
		return nil, err
	}

	//定时上报
	go m.reporter.Run()

	return m, nil
}

//Handle 处理请求
func (m *Metric) Handle() swap.Handler {
	return func(r swap.IContext) {

		//1. 初始化三类统计器---请求的QPS/正在处理的计数器/时间统计器
		url := ctx.Request().GetService()
		conterName := metrics.MakeName(m.serverType+".server.request", metrics.WORKING, "server", m.name, "host", m.ip, "url", url) //堵塞计数
		timerName := metrics.MakeName(m.serverType+".server.request", metrics.TIMER, "server", m.name, "host", m.ip, "url", url)    //堵塞计数
		requestName := metrics.MakeName(m.serverType+".server.request", metrics.QPS, "server", m.name, "host", m.ip, "url", url)    //请求数

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
		statusCode := ctx.GetStatusCode()
		responseName := metrics.MakeName(m.serverType+".server.response", metrics.METER, "server", m.name, "host", m.ip,
			"url", url, "status", fmt.Sprintf("%d", statusCode)) //完成数

		//7. 对服务处理结果的状态码进行上报
		metrics.GetOrRegisterMeter(responseName, m.currentRegistry).Mark(1)
	}

}

//Stop stop metric
func (m *Metric) Stop() {
	if m.reporter != nil {
		m.reporter.Close()
	}
}
