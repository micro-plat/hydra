package component

import (
	"fmt"

	"github.com/asaskevich/govalidator"
	"github.com/micro-plat/hydra/conf"
	"github.com/micro-plat/hydra/registry"
	"github.com/micro-plat/lib4go/concurrent/cmap"
	"github.com/micro-plat/lib4go/influxdb"
)

//InfluxDBTypeNameInVar influxdb在var配置中的类型名称
const InfluxDBTypeNameInVar = "influxdb"

//InfluxDBNameInVar influxdb名称在var配置中的末节点名称
const InfluxDBNameInVar = "influxdb"

var _ IComponentInfluxDB = &StandardInfluxDB{}

//IComponentInfluxDB Component DB
type IComponentInfluxDB interface {
	GetRegularInflux(names ...string) (c influxdb.IInfluxClient)
	GetInflux(names ...string) (d influxdb.IInfluxClient, err error)
	GetInfluxBy(tpName string, name string) (c influxdb.IInfluxClient, err error)
	SaveInfluxObject(tpName string, name string, f func(c conf.IConf) (influxdb.IInfluxClient, error)) (bool, influxdb.IInfluxClient, error)
	Close() error
}

//StandardInfluxDB db
type StandardInfluxDB struct {
	IContainer
	name          string
	influxdbCache cmap.ConcurrentMap
}

//NewStandardInfluxDB 创建DB
func NewStandardInfluxDB(c IContainer, name ...string) *StandardInfluxDB {
	if len(name) > 0 {
		return &StandardInfluxDB{IContainer: c, name: name[0], influxdbCache: cmap.New(2)}
	}
	return &StandardInfluxDB{IContainer: c, name: InfluxDBNameInVar, influxdbCache: cmap.New(2)}
}

//GetRegularInflux 获取正式的没有异常Influx实例
func (s *StandardInfluxDB) GetRegularInflux(names ...string) (c influxdb.IInfluxClient) {
	c, err := s.GetInflux(names...)
	if err != nil {
		panic(err)
	}
	return c
}

//GetInflux get influxdb
func (s *StandardInfluxDB) GetInflux(names ...string) (influxdb.IInfluxClient, error) {
	name := s.name
	if len(names) > 0 {
		name = names[0]
	}
	return s.GetInfluxBy(InfluxDBTypeNameInVar, name)

}

//GetInfluxBy 根据类型获取缓存数据
func (s *StandardInfluxDB) GetInfluxBy(tpName string, name string) (c influxdb.IInfluxClient, err error) {
	_, c, err = s.SaveInfluxObject(tpName, name, func(jConf conf.IConf) (influxdb.IInfluxClient, error) {
		var metric conf.Metric
		if err := jConf.Unmarshal(&metric); err != nil {
			return nil, err
		}
		if b, err := govalidator.ValidateStruct(&metric); !b {
			return nil, err
		}
		return influxdb.NewInfluxClient(metric.Host, metric.DataBase, metric.UserName, metric.Password)
	})
	return c, err
}

//SaveInfluxObject 缓存对象
func (s *StandardInfluxDB) SaveInfluxObject(tpName string, name string, f func(c conf.IConf) (influxdb.IInfluxClient, error)) (bool, influxdb.IInfluxClient, error) {
	cacheConf, err := s.IContainer.GetVarConf(tpName, name)
	if err != nil {
		return false, nil, fmt.Errorf("%s %v", registry.Join("/", s.GetPlatName(), "var", tpName, name), err)
	}
	key := fmt.Sprintf("%s/%s:%d", tpName, name, cacheConf.GetVersion())
	ok, ch, err := s.influxdbCache.SetIfAbsentCb(key, func(input ...interface{}) (c interface{}, err error) {
		return f(cacheConf)
	})
	if err != nil {
		err = fmt.Errorf("创建influxdb失败:%s,err:%v", string(cacheConf.GetRaw()), err)
		return ok, nil, err
	}
	return ok, ch.(influxdb.IInfluxClient), err
}

//Close 释放所有缓存配置
func (s *StandardInfluxDB) Close() error {
	s.influxdbCache.RemoveIterCb(func(k string, v interface{}) bool {
		v.(*influxdb.InfluxClient).Close()
		return true
	})
	return nil
}
