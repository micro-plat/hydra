package lmq

//LMQ
type LMQ struct {
	Proto string `json:"proto,omitempty"`
}

//New 构建mqtt配置
func New() *LMQ {
	return &LMQ{Proto: "lmq"}
}
