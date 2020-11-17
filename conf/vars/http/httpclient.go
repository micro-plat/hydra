package http

const (
	//typeNode DB在var配置中的类型名称
	HttpTypeNode = "http"

	//nameNode DB名称在var配置中的末节点名称
	HttpNameNode = "client"
)

//HTTPConf http客户端配置对象
type HTTPConf struct {
	RequestID         string   `json:"-"`
	ConnectionTimeout int      `json:"ctime"`
	RequestTimeout    int      `json:"rtime"`
	Certs             []string `json:"certs"`
	Ca                string   `json:"ca"`
	Proxy             string   `json:"proxy"`
	Keepalive         bool     `json:"keepalive"`
	Trace             bool     `json:"trace"`
}

//New 构建http 客户端配置信息
func New(opts ...Option) *HTTPConf {
	httpConf := &HTTPConf{
		ConnectionTimeout: 30,
		RequestTimeout:    10,
	}
	for _, opt := range opts {
		opt(httpConf)
	}

	return httpConf
}
