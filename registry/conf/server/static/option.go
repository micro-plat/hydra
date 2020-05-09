package static

//option 配置参数
type option struct {
	Dir       string   `json:"dir,omitempty" valid:"ascii" toml:"dir,omitempty"`
	Archive   string   `json:"archive,omitempty" valid:"ascii" toml:"archive,omitempty"`
	Prefix    string   `json:"prefix,omitempty" valid:"ascii" toml:"prefix,omitempty"`
	Exts      []string `json:"exts,omitempty" valid:"ascii" toml:"exts,omitempty"`
	Exclude   []string `json:"exclude,omitempty" valid:"ascii" toml:"exclude,omitempty"`
	FirstPage string   `json:"first-page,omitempty" valid:"ascii" toml:"first-page,omitempty"`
	Rewriters []string `json:"rewriters,omitempty" valid:"ascii" toml:"rewriters,omitempty"`
	Disable   bool     `json:"disable,omitempty" toml:"disable,omitempty"`
}

//Option jwt配置选项
type Option func(*option)

//WithExts 构建Web服务静态文件配置
func newOption() *option {
	a := &option{}
	a.Dir = "./static"
	a.FirstPage = "index.html"
	a.Rewriters = []string{"/", "index.htm", "default.html"}
	a.Exclude = []string{"/views/", ".exe", ".so"}
	a.Exts = []string{".txt", ".html", ".htm", ".js", ".css", ".map", ".ttf", ".woff", ".woff2", ".woff2", ".jpg", ".jpeg", ".png", ".gif", ".ico", ".tif", ".pcx", ".tga", ".exif", ".fpx", ".svg", ".psd", ".cdr", ".pcd", ".dxf", ".ufo", ".eps", ".ai", ".raw", ".WMF", ".webp"}
	return a
}

//WithImages 图片服务配置
func WithImages() Option {
	return func(s *option) {
		s.Dir = "./static"
		s.Exts = []string{".jpg", ".jpeg", ".png", ".gif", ".ico", ".tif", ".pcx", ".tga", ".exif", ".fpx", ".svg", ".psd", ".cdr", ".pcd", ".dxf", ".ufo", ".eps", ".ai", ".raw", ".WMF", ".webp"}
	}
}

//WithRoot 设置静态文件跟目录
func WithRoot(dir string) Option {
	return func(s *option) {
		s.Dir = dir
	}
}

//WithFirstPage 设置静首页地址
func WithFirstPage(firstPage string) Option {
	return func(s *option) {
		s.FirstPage = firstPage
	}
}

//WithExts 设置静态文件跟目录
func WithExts(exts ...string) Option {
	return func(s *option) {
		s.Exts = exts
	}
}

//WithArchive 设置静态文件跟目录
func WithArchive(archive string) Option {
	return func(s *option) {
		s.Archive = archive
	}
}

//AppendExts 设置静态文件跟目录
func AppendExts(exts ...string) Option {
	return func(s *option) {
		s.Exts = append(s.Exts, exts...)
	}
}

//WithPrefix 设置静态文件跟目录
func WithPrefix(prefix string) Option {
	return func(s *option) {
		s.Prefix = prefix
	}
}

//WithDisable 禁用配置
func WithDisable() Option {
	return func(a *option) {
		a.Disable = true
	}
}

//WithEnable 启用配置
func WithEnable() Option {
	return func(a *option) {
		a.Disable = false
	}
}
