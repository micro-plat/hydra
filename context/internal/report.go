package internal

// import (
// 	"fmt"
// 	"time"

// 	"github.com/SkyAPM/go2sky"
// 	"github.com/SkyAPM/go2sky/reporter"
// 	"github.com/micro-plat/hydra/conf/app"
// 	"github.com/micro-plat/hydra/conf/server/apm"
// 	"github.com/micro-plat/lib4go/concurrent/cmap"
// )

// var reporters = cmap.New(3)

// func getReporter(c app.IAPPConf) (go2sky.Reporter, error) {
// 	conf, err := c.GetAPMConf()
// 	if err != nil || conf.Disable {
// 		return nil, err
// 	}
// 	key := fmt.Sprintf("apm:%s", conf.Address)
// 	_, v, err := reporters.SetIfAbsentCb(key, func(v ...interface{}) (i interface{}, err error) {
// 		apm := v[0].(*apm.APM)
// 		report, err := reporter.NewGRPCReporter(apm.Address, reporter.WithCheckInterval(time.Second))
// 		return report, err

// 	}, conf)
// 	go remove(key)
// 	if err != nil || v == nil {
// 		return nil, err
// 	}

// 	return v.(go2sky.Reporter), nil
// }
// func remove(key string) {
// 	reporters.RemoveIterCb(func(k string, v interface{}) bool {
// 		return key != k
// 	})

// }
