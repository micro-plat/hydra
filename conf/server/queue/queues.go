package queue

import (
	"fmt"

	"github.com/asaskevich/govalidator"
	"github.com/micro-plat/hydra/conf"
)

//TypeNodeName 分类节点名
const TypeNodeName = "queue"

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
func (q *Queues) Append(queues ...*Queue) (*Queues, []*Queue) {
	keyMap := map[string]*Queue{}
	for _, v := range q.Queues {
		keyMap[v.Queue] = v
	}
	notifyQueues := []*Queue{}
	for _, v := range queues {
		if queue, ok := keyMap[v.Queue]; ok {
			if queue.Disable != v.Disable || queue.Concurrency != v.Concurrency {
				notifyQueues = append(notifyQueues, v)
				queue.Disable = v.Disable
				queue.Concurrency = v.Concurrency
			}
			continue
		}
		notifyQueues = append(notifyQueues, v)
		q.Queues = append(q.Queues, v)
	}
	return q, notifyQueues
}

//GetConf 设置queue
func GetConf(cnf conf.IServerConf) (queues *Queues, err error) {
	queues = &Queues{}
	if _, err := cnf.GetSubObject(TypeNodeName, queues); err != nil && err != conf.ErrNoSetting {
		return nil, fmt.Errorf("queues配置格式有误:%v", err)
	}

	for _, queue := range queues.Queues {
		if b, err := govalidator.ValidateStruct(queue); !b {
			return nil, fmt.Errorf("queue配置数据有误:%v", err)
		}
	}

	return queues, nil
}
