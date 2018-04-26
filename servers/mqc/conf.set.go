package mqc

import (
	"errors"
	"fmt"

	"github.com/asaskevich/govalidator"
	"github.com/micro-plat/hydra/conf"
	"github.com/micro-plat/hydra/servers"
	"github.com/micro-plat/hydra/servers/pkg/middleware"
)

type ISetMetric interface {
	SetMetric(*conf.Metric) error
}

//SetMetric 设置metric
func SetMetric(set ISetMetric, cnf conf.IServerConf) (enable bool, err error) {
	//设置静态文件路由
	var metric conf.Metric
	_, err = cnf.GetSubObject("metric", &metric)
	if err != nil && err != conf.ErrNoSetting {
		return false, err
	}
	if err == conf.ErrNoSetting {
		metric.Disable = true
	} else {
		if b, err := govalidator.ValidateStruct(&metric); !b {
			err = fmt.Errorf("metric配置有误:%v", err)
			return false, err
		}
	}
	err = set.SetMetric(&metric)
	return !metric.Disable && err == nil, err
}

//IQueues 设置queue
type IQueues interface {
	SetQueues(string, string, []*conf.Queue) error
}

//SetQueues 设置queue
func SetQueues(engine servers.IRegistryEngine, set IQueues, cnf conf.IServerConf, ext map[string]interface{}) (enable bool, err error) {

	serverConf, err := cnf.GetSubConf("server")
	if err == conf.ErrNoSetting {
		err = fmt.Errorf("server节点:%v", err)
		return false, err
	}
	if err != nil {
		return false, err
	}
	var server conf.Server
	if err = serverConf.Unmarshal(&server); err != nil {
		return false, err
	}
	if b, err := govalidator.ValidateStruct(&server); !b {
		err = fmt.Errorf("server配置有误:%v", err)
		return false, err
	}
	var queues conf.Queues
	if _, err = cnf.GetSubObject("queue", &queues); err == conf.ErrNoSetting {
		err = fmt.Errorf("queue:%v", err)
		return false, err
	}
	if err != nil {
		return false, err
	}
	if len(queues.Queues) == 0 {
		return false, errors.New("queue:未配置")
	}
	if b, err := govalidator.ValidateStruct(&queues); !b {
		err = fmt.Errorf("queue配置有误:%v", err)
		return false, err
	}
	nqueues := make([]*conf.Queue, 0, len(queues.Queues))
	for _, queue := range queues.Queues {
		if queue.Disable {
			continue
		}
		if queue.Name == "" {
			queue.Name = queue.Service
		}
		if queue.Setting == nil {
			queue.Setting = make(map[string]string)
		}
		for k, v := range queues.Setting {
			if _, ok := queue.Setting[k]; !ok {
				queue.Setting[k] = v
			}
		}
		queue.Handler = middleware.ContextHandler(engine, queue.Name, queue.Engine, queue.Service, queue.Setting, ext)
		nqueues = append(nqueues, queue)
	}
	if err = set.SetQueues(server.Proto, string(serverConf.GetRaw()), nqueues); err != nil {
		return false, err
	}
	return len(nqueues) > 0, nil
}
