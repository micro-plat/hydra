package static

import "strings"

//DefaultSataticDir 默认静态文件存放路径
const DefaultSataticDir = "./src"

//DefaultFirstPage 默认首页文件名
const DefaultFirstPage = "index.html"

//DefaultRewriters 默认需要重写的路径
var DefaultRewriters = []string{"/", "index.htm", "default.html", "default.htm"}

//DefaultExclude 默认需要排除的文件,扩展名,路径
var DefaultExclude = []string{"/view/", "/views/", "/web/", ".exe", ".so"}

//Option jwt配置选项
type Option func(*Static)

//newStatic 构建Web服务静态文件配置
func newStatic() *Static {
	a := &Static{
		FileMap: map[string]FileInfo{},
	}
	a.Dir = DefaultSataticDir
	a.FirstPage = DefaultFirstPage
	a.Rewriters = DefaultRewriters
	a.Exclude = DefaultExclude
	a.Exts = []string{}
	return a
}

//WithImages 图片服务配置
func WithImages() Option {
	return func(s *Static) {
		s.Dir = DefaultSataticDir
		s.Exts = []string{}
	}
}

//WithRewriters 图片服务配置
func WithRewriters(rewriters ...string) Option {
	return func(s *Static) {
		s.Rewriters = rewriters
	}
}

//WithExclude 图片服务配置
func WithExclude(exclude ...string) Option {
	return func(s *Static) {
		s.Exclude = exclude
	}
}

//WithRoot 设置静态文件跟目录
func WithRoot(dir string) Option {
	return func(s *Static) {
		s.Dir = dir
	}
}

//WithFirstPage 设置静首页地址
func WithFirstPage(firstPage string) Option {
	return func(s *Static) {
		s.FirstPage = firstPage
	}
}

//WithExts 设置静态文件跟目录
func WithExts(exts ...string) Option {
	return func(s *Static) {
		s.Exts = exts
	}
}

//WithArchive 设置静态文件跟目录
func WithArchive(archive string) Option {
	return func(s *Static) {
		if !strings.Contains(archive, ".") {
			s.Archive = archive + ".zip"
			return
		}
		s.Archive = archive
	}
}

//AppendExts 设置静态文件跟目录
func AppendExts(exts ...string) Option {
	return func(s *Static) {
		s.Exts = append(s.Exts, exts...)
	}
}

//WithPrefix 设置静态文件跟目录
func WithPrefix(prefix string) Option {
	return func(s *Static) {
		s.Prefix = prefix
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
