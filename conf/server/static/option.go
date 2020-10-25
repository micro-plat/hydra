package static

import "strings"

//Option jwt配置选项
type Option func(*Static)

//newStatic 构建Web服务静态文件配置
func newStatic() *Static {
	a := &Static{
		FileMap: map[string]FileInfo{},
	}
	a.Dir = "./src"
	a.FirstPage = "index.html"
	a.Rewriters = []string{"/", "index.htm", "default.html"}
	a.Exclude = []string{"/views/", ".exe", ".so"}
	a.Exts = []string{}
	return a
}

//WithImages 图片服务配置
func WithImages() Option {
	return func(s *Static) {
		s.Dir = "./src"
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
