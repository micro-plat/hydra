package conf

type RemoteAuth struct {
	RPCServiceName string   `json:"rpc-service" valid:"required"`
	Include        []string `json:"include" valid:"required"`
	Disable        bool     `json:"disable,omitempty"`
}

//WithRemoteAuth 添加固定签名认证
func (a *Authes) WithRemoteAuth(auth *RemoteAuth) *Authes {
	a.RemoteAuth = auth
	return a
}

//NewRemoteAuth 创建固定Secret签名认证
func NewRemoteAuth(rpcService string, path ...string) *RemoteAuth {
	ninclude := []string{"*"}
	if len(path) > 0 {
		ninclude = path
	}
	return &RemoteAuth{
		RPCServiceName: rpcService,
		Include:        ninclude,
	}
}

//Contains 检查指定的路径是否允许签名
func (a *RemoteAuth) Contains(p string) bool {
	if len(a.Include) == 0 {
		return true
	}
	for _, i := range a.Include {
		if i == "*" || i == p {
			return true
		}
	}
	return false
}

//WithInclude 设置include的请求服务路径
func (a *RemoteAuth) WithInclude(path ...string) *RemoteAuth {
	if len(path) > 0 {
		a.Include = path
	}
	return a

//WithDisable 禁用配置
func (a *RemoteAuth) WithDisable() *RemoteAuth {
	a.Disable = true
	return a
}

//WithEnable 启用配置
func (a *RemoteAuth) WithEnable() *RemoteAuth {
	a.Disable = false
	return a
}
