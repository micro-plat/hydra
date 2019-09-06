package conf

type RemoteAuth struct {
	RPCServiceName string `json:"rpc-service" valid:"required"`
	Disable        bool   `json:"disable,omitempty"`
}

//WithRemoteAuth 添加固定签名认证
func (a *Authes) WithRemoteAuth(auth *RemoteAuth) *Authes {
	a.RemoteAuth = auth
	return a
}

//NewRemoteAuth 创建固定Secret签名认证
func NewRemoteAuth(rpcService string) *RemoteAuth {
	return &RemoteAuth{
		RPCServiceName: rpcService,
	}
}

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
