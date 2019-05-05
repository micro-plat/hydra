package zk

//Option 配置选项
type Option func(*ZookeeperClient)

func WithdDigest(u, p string) Option {
	return func(o *ZookeeperClient) {
		if u != "" {
			o.userName = u
			o.password = p
			o.digest = true
		}

	}
}
