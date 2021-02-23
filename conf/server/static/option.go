package static

import "embed"

//DefaultSataticDir 默认静态文件存放路径
const DefaultSataticDir = "./static"

//DefaultHome 默认首页文件名
const DefaultHome = "index.html"

//DefaultExclude 默认需要排除的文件,扩展名,路径
var DefaultExclude = []string{".exe", ".so"}

//DefaultUnrewrite 默认不重写文件
var DefaultUnrewrite = []string{"/favicon.ico", "/robots.txt"}

//Option jwt配置选项
type Option func(*Static)

//WithExclude 排除配置
func WithExclude(excludes ...string) Option {
	return func(s *Static) {
		s.Excludes = excludes
	}
}

//WithUnrewrite 不重写列表
func WithUnrewrite(list ...string) Option {
	return func(s *Static) {
		s.Unrewrites = list
	}
}

//WithAssetsPath 设置资源地址
func WithAssetsPath(path string) Option {
	return func(s *Static) {
		s.Path = path
	}
}

//WithHomePage 设置静首页地址
func WithHomePage(homePage string) Option {
	return func(s *Static) {
		s.HomePage = homePage
	}
}

//WithEmbed 通过嵌入的方式指定压缩文件
func WithEmbed(root string, fs embed.FS) Option {
	return func(s *Static) {
		defEmbedFs.archive = &fs
		defEmbedFs.name = root
	}
}

//WithEmbedBytes 通过嵌入的方式指定压缩文件
func WithEmbedBytes(fileName string, bytes []byte) Option {
	return func(s *Static) {
		defEmbedFs.bytes = bytes
		defEmbedFs.name = fileName
	}
}

//WithAutoRewrite 设置为自动重写
func WithAutoRewrite() Option {
	return func(a *Static) {
		a.AutoRewrite = true
	}
}

//WithDisable 禁用配置
func WithDisable() Option {
	return func(a *Static) {
		a.Disable = true
	}
}

//WithEnable 启用配置
func WithEnable() Option {
	return func(a *Static) {
		a.Disable = false
	}
}
