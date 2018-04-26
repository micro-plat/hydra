package cron

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

//ITasks 设置tasks
type ITasks interface {
	SetTasks(string, []*conf.Task) error
}

//SetTasks 设置tasks
func SetTasks(engine servers.IRegistryEngine, set ITasks, cnf conf.IServerConf, ext map[string]interface{}) (enable bool, err error) {
	reidsConf, err := cnf.GetSubConf("redis")
	if err != nil && err != conf.ErrNoSetting {
		return false, err
	}

	var tasks conf.Tasks
	if _, err = cnf.GetSubObject("task", &tasks); err == conf.ErrNoSetting {
		err = fmt.Errorf("task:%v", err)
		return false, err
	}
	if err != nil {
		return false, err
	}
	if len(tasks.Tasks) == 0 {
		err = errors.New("task:未配置")
		return false, err
	}

	if b, err := govalidator.ValidateStruct(&tasks); !b {
		err = fmt.Errorf("task配置有误:%v", err)
		return false, err
	}
	ntasks := make([]*conf.Task, 0, len(tasks.Tasks))
	for _, task := range tasks.Tasks {
		if task.Disable {
			continue
		}
		if task.Name == "" {
			task.Name = task.Service
		}
		if task.Setting == nil {
			task.Setting = make(map[string]string)
		}
		for k, v := range tasks.Setting {
			if _, ok := task.Setting[k]; !ok {
				task.Setting[k] = v
			}
		}
		task.Handler = middleware.ContextHandler(engine, task.Name, task.Engine, task.Service, task.Setting, ext)
		ntasks = append(ntasks, task)
	}
	var raw string
	if reidsConf != nil {
		raw = string(reidsConf.GetRaw())
	}
	if err = set.SetTasks(raw, ntasks); err != nil {
		return false, err
	}
	return len(ntasks) > 0, nil
}
