package db

//Option 配置选项
type Option func(*DB)

//WithConnect 设置数据库连接信息：最大打开数，空闲数，连接超时时长
func WithConnect(maxOpen int, maxIdle int, lifeTime int) Option {
	return func(a *DB) {
		a.MaxOpen = maxOpen
		a.MaxIdle = maxIdle
		a.LifeTime = lifeTime
	}
}

//WithEnableEncryption 启用加密设置
func WithEnableEncryption() Option {
	return func(a *DB) {
		a.EnableEncryption = true
	}
}

//WithTenantMode 设置多租户模式：shared_schema | shared_rls | indep_schema
func WithTenantMode(mode string) Option {
	return func(a *DB) {
		a.TenantMode = mode
	}
}
