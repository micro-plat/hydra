package queue

import (
	"fmt"

	"github.com/asaskevich/govalidator"
	"github.com/micro-plat/hydra/registry/conf"
)

//Queues queue任务
type Queues struct {
	Queues []*Queue `json:"queues" toml:"queues,omitempty"`
}

//NewEmptyQueues 构建空的queues
func NewEmptyQueues() *Queues {
	return &Queues{
		Queues: make([]*Queue, 0),
	}
}

//NewQueues 构建Queues
func NewQueues(queue ...*Queue) *Queues {
	q := NewEmptyQueues()
	q.Queues = append(q.Queues, queue...)
	return q
}

//Append 增加任务列表
func (t *Queues) Append(queues ...*Queue) *Queues {
	for _, q := range queues {
		t.Queues = append(t.Queues, q)
	}
	return t
}

//GetConf 设置queue
func GetConf(cnf conf.IMainConf) (queues *Queues, err error) {
	if _, err = cnf.GetSubObject("queue", &queues); err != nil && err != conf.ErrNoSetting {
		return nil, fmt.Errorf("queue:%v", err)
	}
	if len(queues.Queues) > 0 {
		if b, err := govalidator.ValidateStruct(&queues); !b {
			return nil, fmt.Errorf("queue配置有误:%v", err)
		}
	}
	return queues, nil
}
