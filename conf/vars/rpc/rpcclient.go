package rpc

//RPCTypeNode rpc在var配置中的类型名称
const RPCTypeNode = "rpc"

//RPCNameNode rpc名称在var配置中的末节点名称
const RPCNameNode = "rpc"

//RPCConf http客户端配置对象
type RPCConf struct {
	ConntTimeout int      `json:"connectionTimeout"`
	Log          string   `json:"log"`
	LocalIP      string   `json:"localIP"`
	SortPrefix   string   `json:"sortPrefix"`
	Tls          []string `json:"tls"`
	LocalFirst   bool     `json:"localFirst"` // 本地优先负载
	RoundRobin   bool     `json:"roundRobin"` // 论寻负载
}

//New 构建http 客户端配置信息
func New(opts ...Option) *RPCConf {
	rpcConf := &RPCConf{
		ConntTimeout: 3,
		RoundRobin:   true, //默认论寻负载
	}
	for _, opt := range opts {
		opt(rpcConf)
	}

	return rpcConf
}
