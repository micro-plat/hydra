package middleware

import (
	"fmt"
	"sync"

	"github.com/gin-gonic/gin"
	"github.com/micro-plat/hydra/conf"
	"github.com/micro-plat/hydra/servers/pkg/timer"
	"github.com/micro-plat/lib4go/concurrent/cmap"
	"github.com/micro-plat/lib4go/logger"
	"github.com/micro-plat/lib4go/metrics"
	"github.com/micro-plat/lib4go/net"
	"github.com/micro-plat/lib4go/xsync"
)

type reporter struct {
	reporter metrics.IReporter
	Host     string
	Database string
	username string
	password string
	cron     string
}

//Metric 服务器处理能力统计
type Metric struct {
	logger          *logger.Logger
	reporter        *reporter
	registry        cmap.ConcurrentMap
	ticket          *xsync.Ticket
	mu              sync.Mutex
	currentRegistry metrics.Registry
	conf            *conf.MetadataConf
	ip              string
	timer           *timer.Timer
	done            bool
	closeChan       chan struct{}
}

//NewMetric new metric
func NewMetric(conf *conf.MetadataConf) *Metric {
	return &Metric{
		conf:            conf,
		currentRegistry: metrics.NewRegistry(),
		ip:              net.GetLocalIPAddress(),
		closeChan:       make(chan struct{}),
		ticket:          xsync.Sequence.Get(),
	}
}

//Stop stop metric
func (m *Metric) Stop() {
	m.mu.Lock()
	defer m.mu.Unlock()
	m.ticket.Done()
	if !m.done {
		close(m.closeChan)
	}
	m.done = true
	if m.reporter != nil && m.reporter.reporter != nil {
		m.reporter.reporter.Close()
	}
	if m.timer != nil {
		m.timer.Close()
	}
}

//Restart restart metric
func (m *Metric) Restart(host string, dataBase string, userName string, password string, c string,
	lg *logger.Logger) (err error) {
	m.Stop()

	m.done = false
	m.closeChan = make(chan struct{})
	m.timer, err = timer.NewTimer(c)
	if err != nil {
		return err
	}
	m.logger = lg
	m.reporter = &reporter{Host: host, Database: dataBase, username: userName, password: password, cron: c}
	m.reporter.reporter, err = metrics.InfluxDB(m.currentRegistry,
		c,
		m.reporter.Host, m.reporter.Database,
		m.reporter.username,
		m.reporter.password, m.logger)
	if err != nil {
		return
	}

	go m.reporter.reporter.Run()
	//	go m.collectSys()
	m.timer.Start()
	return nil
}
func (m *Metric) collectSys() {
	if !m.ticket.Wait() {
		return
	}
	go m.loopCollectCPU()
	go m.loopCollectDisk()
	go m.loopCollectMem()
	go m.loopNetConnCount()
}

//Handle 处理请求
func (m *Metric) Handle() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		url := ctx.Request.URL.Path
		conterName := metrics.MakeName(m.conf.Type+".server.request", metrics.WORKING, "name", m.conf.Name, "host", m.ip, "url", url) //堵塞计数
		timerName := metrics.MakeName(m.conf.Type+".server.request", metrics.TIMER, "name", m.conf.Name, "host", m.ip, "url", url)    //堵塞计数
		requestName := metrics.MakeName(m.conf.Type+".server.request", metrics.QPS, "name", m.conf.Name, "host", m.ip, "url", url)    //请求数
		metrics.GetOrRegisterQPS(requestName, m.currentRegistry).Mark(1)

		counter := metrics.GetOrRegisterCounter(conterName, m.currentRegistry)
		counter.Inc(1)
		metrics.GetOrRegisterTimer(timerName, m.currentRegistry).Time(func() { ctx.Next() })
		counter.Dec(1)

		statusCode := ctx.Writer.Status()
		responseName := metrics.MakeName(m.conf.Type+".server.response", metrics.METER, "name", m.conf.Name, "host", m.ip,
			"url", url, "status", fmt.Sprintf("%d", statusCode)) //完成数
		metrics.GetOrRegisterMeter(responseName, m.currentRegistry).Mark(1)
	}
}
