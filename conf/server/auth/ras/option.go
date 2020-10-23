package ras

//RASOption 配置选项
type RASOption func(*RASAuth)

// WithAuthList WithAuthList
func WithAuthList(list ...*Auth) RASOption {
	return func(a *RASAuth) {
		for _, item := range list {
			a.Auth = append(a.Auth, item)
		}
	}
}

//WithRASDisable 关闭
func WithRASDisable() RASOption {
	return func(a *RASAuth) {
		a.Disable = true
	}
}

//WithRASEnable 开启
func WithRASEnable() RASOption {
	return func(a *RASAuth) {
		a.Disable = false
	}
}
