package queue

import (
	"fmt"

	"github.com/asaskevich/govalidator"
	"github.com/micro-plat/hydra/conf"
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

//Append 增加任务列表  @fix:存在的数据进行修改,不存在则添加 @hj
func (t *Queues) Append(queues ...*Queue) *Queues {
	keyMap := map[string]*Queue{}
	for _, v := range t.Queues {
		keyMap[v.Queue] = v
	}
	nonExistQueue := []*Queue{}
	for _, v := range queues {
		if queue, ok := keyMap[v.Queue]; ok {
			queue.Disable = v.Disable
			queue.Concurrency = v.Concurrency
			continue
		}
		nonExistQueue = append(nonExistQueue, v)
	}
	t.Queues = append(t.Queues, nonExistQueue...)
	return t
}

//GetConf 设置queue
func GetConf(cnf conf.IServerConf) (queues *Queues, err error) {
	queues = &Queues{}
	if _, err := cnf.GetSubObject("queue", queues); err != nil && err != conf.ErrNoSetting {
		return nil, fmt.Errorf("queues配置格式有误:%v", err)
	}

	for _, queue := range queues.Queues {
		if b, err := govalidator.ValidateStruct(queue); !b {
			return nil, fmt.Errorf("queue配置数据有误:%v", err)
		}
	}

	return queues, nil
}
