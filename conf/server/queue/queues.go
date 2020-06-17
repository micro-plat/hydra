package queue

import (
	"fmt"

	"github.com/asaskevich/govalidator"
	"github.com/micro-plat/hydra/conf"
)

//VarRootName 在var中的跟路径
const VarRootName = "queue"

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

type ConfHandler func(cnf conf.IMainConf) *Queues

func (h ConfHandler) Handle(cnf conf.IMainConf) interface{} {
	return h(cnf)
}

//GetConf 设置queue
func GetConf(cnf conf.IMainConf) (queues *Queues) {
	queues = &Queues{}
	if _, err := cnf.GetSubObject("queue", queues); err != nil && err != conf.ErrNoSetting {
		panic(err)
	}
	if len(queues.Queues) > 0 {
		if b, err := govalidator.ValidateStruct(queues); !b {
			panic(fmt.Errorf("queue配置有误:%v", err))
		}
	}
	return queues
}
