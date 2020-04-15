package queue

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
