package queue

import (
	"fmt"

	"github.com/asaskevich/govalidator"
	"github.com/micro-plat/hydra/conf"
)

//Queues queue任务
type Queues struct {
	Queues []*Queue `json:"queues"`
}

//Queue 任务项
type Queue struct {
	*option
}

//NewQueues 构建Queues
func NewQueues(v Option, opts ...Option) *Queues {
	q := &Queues{Queues: make([]*Queue, 0, 1)}
	fq := &Queue{option: &option{}}
	v(fq.option)
	q.Queues = append(q.Queues, fq)
	for _, opt := range opts {
		oq := &Queue{option: &option{}}
		opt(oq.option)
		q.Queues = append(q.Queues, oq)
	}
	return q
}

//GetQueues 设置queue
func GetQueues(cnf conf.IServerConf) (queues *conf.Queues, err error) {
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
