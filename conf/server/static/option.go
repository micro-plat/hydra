package static

import "embed"

//DefaultSataticDir 默认静态文件存放路径
const DefaultSataticDir = "./static"

//DefaultHome 默认首页文件名
const DefaultHome = "index.html"

//DefaultExclude 默认需要排除的文件,扩展名,路径
var DefaultExclude = []string{".exe", ".so"}

//Option jwt配置选项
type Option func(*Static)

//WithExclude 图片服务配置
func WithExclude(exclude ...string) Option {
	return func(s *Static) {
		s.Excludes = exclude
	}
}

//WithHomePage 设置静首页地址
func WithHomePage(firstPage string) Option {
	return func(s *Static) {
		s.HomePage = firstPage
	}
}

//WithEmbed 通过嵌入的方式指定压缩文件
func WithEmbed(root string, fs embed.FS) Option {
	return func(s *Static) {
		defEmbedFs.archive = fs
		defEmbedFs.name = root
	}
}

//WithDisable 禁用配置
func WithDisable() Option {
	return func(a *Static) {
		a.Disable = true
	}
}

//WithRewrite 设置为自动重写
func WithRewrite() Option {
	return func(a *Static) {
		a.AutoRewrite = true
	}
}

//WithEnable 启用配置
func WithEnable() Option {
	return func(a *Static) {
		a.Disable = false
	}
}
