package metrics

import (
	"fmt"
	uurl "net/url"
	"strings"
	"time"

	"github.com/micro-plat/lib4go/influxdb"
	"github.com/micro-plat/lib4go/logger"
	"github.com/zkfy/cron"
	"github.com/zkfy/go-metrics"
)

type IReporter interface {
	Close() error
	Run()
}
type reporter struct {
	reg      metrics.Registry
	schedule cron.Schedule
	cron     string
	url      uurl.URL
	database string
	username string
	password string
	tags     map[string]string
	client   *influxdb.Client
	logger   *logger.Logger
	done     bool
}

const (
	WORKING       = "working"
	COUNTER       = "counter"
	GAUGE         = "gauge"
	GAUGEFLOAST64 = "gaugeFloat64"
	HISTOGRAM     = "histogram"
	METER         = "meter"
	TIMER         = "timer"
	QPS           = "qps"
)

// InfluxDB starts a InfluxDB reporter which will post the metrics from the given registry at each d interval.
func InfluxDB(r metrics.Registry, cron string, url, database, username, password string, logger *logger.Logger) (IReporter, error) {
	return InfluxDBWithTags(r, cron, url, database, username, password, nil, logger)
}

//MakeName 构建参数名称
func MakeName(name string, tp string, params ...string) string {
	if len(params)%2 != 0 {
		panic("MakeName params必须成对输入")
	}
	return name + "." + tp + ":>" + strings.Join(params, ":>")
}

//timer.merchant.api.request-server-192.168.0.240-client-127.0.0.1-url-/colin
func splitGroup(name string) (string, map[string]string) {
	names := strings.Split(name, ":>")
	tags := make(map[string]string)
	count := len(names)
	for i := 1; i < count; i++ {
		if i%2 == 1 && i+1 < count {
			tags[names[i]] = names[i+1]
		}
	}
	return names[0], tags
}

// InfluxDBWithTags starts a InfluxDB reporter which will post the metrics from the given registry at each d interval with the specified tags
func InfluxDBWithTags(r metrics.Registry, c string, url, database, username, password string, tags map[string]string, logger *logger.Logger) (IReporter, error) {
	sch, err := cron.ParseStandard(c)
	if err != nil {
		return nil, err
	}

	u, err := uurl.Parse(url)
	if err != nil {
		return nil, fmt.Errorf("unable to parse InfluxDB url %s. err=%v", url, err)
	}

	rep := &reporter{
		logger:   logger,
		cron:     c,
		reg:      r,
		schedule: sch,
		url:      *u,
		database: database,
		username: username,
		password: password,
		tags:     tags,
	}
	if err := rep.makeClient(); err != nil {
		return nil, fmt.Errorf("unable to make InfluxDB client. err=%v", err)
	}

	return rep, nil
}
func (r *reporter) Run() {
	r.run()
}
func (r *reporter) makeClient() (err error) {
	r.client, err = influxdb.NewClient(influxdb.Config{
		URL:      r.url,
		Username: r.username,
		Password: r.password,
	})

	return
}

func (r *reporter) run() {
	pingTicker := time.Tick(time.Second * 5)
	var intervalTicker int64
LOOP:
	for {
		select {
		case <-time.After(time.Second):
			if r.done {
				break LOOP
			}
			now := time.Now()
			if intervalTicker > now.Unix() {
				break
			}
			go func() {
				if err := r.send(); err != nil {
					r.logger.Errorf("unable to send metrics to InfluxDB. err=%v", err)
				}
			}()
			intervalTicker = r.schedule.Next(now).Unix()
		case <-pingTicker:
			_, _, err := r.client.Ping()
			if err != nil {
				r.logger.Errorf("got error while sending a ping to InfluxDB, trying to recreate client. err=%v", err)

				if err = r.makeClient(); err != nil {
					r.logger.Errorf("unable to make InfluxDB client. err=%v", err)
				}
			}
		}
	}
}

func (r *reporter) send() error {
	var pts []influxdb.Point
	r.reg.Each(func(name string, obj interface{}) {
		now := time.Now()
		rname, tags := splitGroup(name)
		switch metric := obj.(type) {
		case IQPS:
			metric.Mark(0)
			pts = append(pts, influxdb.Point{
				Measurement: rname,
				Tags:        tags,
				Fields: map[string]interface{}{
					"m1":  metric.M1(),
					"m5":  metric.M5(),
					"m15": metric.M15(),
				},
				Time: now,
			})
		case Counter:
			ms := metric.Snapshot()
			pts = append(pts, influxdb.Point{
				Measurement: rname,
				Tags:        tags,
				Fields: map[string]interface{}{
					"value": ms.Count(),
				},
				Time: now,
			})
		case Gauge:
			ms := metric.Snapshot()
			pts = append(pts, influxdb.Point{
				Measurement: rname,
				Tags:        tags,
				Fields: map[string]interface{}{
					"value": ms.Value(),
				},
				Time: now,
			})
		case GaugeFloat64:
			ms := metric.Snapshot()
			pts = append(pts, influxdb.Point{
				Measurement: rname,
				Tags:        tags,
				Fields: map[string]interface{}{
					"value": ms.Value(),
				},
				Time: now,
			})
		case Histogram:
			ms := metric.Snapshot()
			ps := ms.Percentiles([]float64{0.5, 0.75, 0.95, 0.99, 0.999, 0.9999})
			pts = append(pts, influxdb.Point{
				Measurement: rname,
				Tags:        tags,
				Fields: map[string]interface{}{
					"count":    ms.Count(),
					"max":      ms.Max(),
					"mean":     ms.Mean(),
					"min":      ms.Min(),
					"stddev":   ms.StdDev(),
					"variance": ms.Variance(),
					"p50":      ps[0],
					"p75":      ps[1],
					"p95":      ps[2],
					"p99":      ps[3],
					"p999":     ps[4],
					"p9999":    ps[5],
				},
				Time: now,
			})
		case Meter:
			ms := metric.Snapshot()
			pts = append(pts, influxdb.Point{
				Measurement: rname,
				Tags:        tags,
				Fields: map[string]interface{}{
					"count": ms.Count(),
					"m1":    ms.Rate1(),
					"m5":    ms.Rate5(),
					"m15":   ms.Rate15(),
					"mean":  ms.RateMean(),
				},
				Time: now,
			})
		case Timer:
			ms := metric.Snapshot()
			ps := ms.Percentiles([]float64{0.5, 0.75, 0.95, 0.99, 0.999, 0.9999})
			pts = append(pts, influxdb.Point{
				Measurement: rname,
				Tags:        tags,
				Fields: map[string]interface{}{
					"count":    ms.Count(),
					"max":      ms.Max(),
					"mean":     ms.Mean(),
					"min":      ms.Min(),
					"stddev":   ms.StdDev(),
					"variance": ms.Variance(),
					"p50":      ps[0],
					"p75":      ps[1],
					"p95":      ps[2],
					"p99":      ps[3],
					"p999":     ps[4],
					"p9999":    ps[5],
					"m1":       ms.Rate1(),
					"m5":       ms.Rate5(),
					"m15":      ms.Rate15(),
					"meanrate": ms.RateMean(),
				},
				Time: now,
			})
		}
	})

	bps := influxdb.BatchPoints{
		Points:   pts,
		Database: r.database,
	}
	_, err := r.client.Write(bps)
	return err
}
func (r *reporter) Close() error {
	r.done = true
	return nil
}
