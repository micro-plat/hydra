package rpc

//RPCTypeNode rpc在var配置中的类型名称
const RPCTypeNode = "rpc"

//RPCNameNode rpc名称在var配置中的末节点名称
const RPCNameNode = "rpc"

//LocalFirst LocalFirst
const LocalFirst = "localfirst"

//RoundRobin RoundRobin
const RoundRobin = "round_robin"

//RPCConf http客户端配置对象
type RPCConf struct {
	ConntTimeout int      `json:"connectionTimeout"`
	Log          string   `json:"log"`
	SortPrefix   string   `json:"sortPrefix"`
	Tls          []string `json:"tls"`
	Balancer     string   `json:"balancer"` //负载类型 localfirst:本地服务优先  round_robin:论寻负载
}

//New 构建http 客户端配置信息
func New(opts ...Option) *RPCConf {
	rpcConf := &RPCConf{
		ConntTimeout: 3,
		Balancer:     LocalFirst,
	}
	for _, opt := range opts {
		opt(rpcConf)
	}

	return rpcConf
}
