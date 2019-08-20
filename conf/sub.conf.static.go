package conf

//Static 设置静态文件配置
type Static struct {
	Dir       string   `json:"dir,omitempty" valid:"ascii"`
	Archive   string   `json:"archive,omitempty" valid:"ascii"`
	Prefix    string   `json:"prefix,omitempty" valid:"ascii"`
	Exts      []string `json:"exts,omitempty" valid:"ascii"`
	Exclude   []string `json:"exclude,omitempty" valid:"ascii"`
	FirstPage string   `json:"first-page,omitempty" valid:"ascii"`
	Rewriters []string `json:"rewriters,omitempty" valid:"ascii"`
	Disable   bool     `json:"disable,omitempty"`
}

//NewWebServerStaticConf 构建Web服务静态文件配置
func NewWebServerStaticConf() *Static {
	return &Static{
		Dir:       "./static",
		Exts:      []string{".txt", ".jpg", ".jpeg", ".png", ".gif", ".ico", ".html", ".htm", ".js", ".css", ".map", ".ttf", ".woff", ".woff2", ".woff2"},
		FirstPage: "index.html",
		Rewriters: []string{"*"},
	}
}

//NewImageStaticConf 构建图片文件配置
func NewImageStaticConf() *Static {
	return &Static{
		Dir:  "./static",
		Exts: []string{".jpg", ".jpeg", ".png", ".gif", ".ico", ".tif", ".pcx", ".tga", ".exif", ".fpx", ".svg", ".psd", ".cdr", ".pcd", ".dxf", ".ufo", ".eps", ".ai", ".raw", ".WMF", ".webp"},
	}
}

//WithRoot 设置静态文件跟目录
func (s *Static) WithRoot(dir string) *Static {
	s.Dir = dir
	return s
}

//WithFirstPage 设置静首页地址
func (s *Static) WithFirstPage(firstPage string) *Static {
	s.FirstPage = firstPage
	return s
}

//WithExts 设置静态文件跟目录
func (s *Static) WithExts(exts ...string) *Static {
	s.Exts = exts
	return s
}

//WithArchive 设置静态文件跟目录
func (s *Static) WithArchive(archive string) *Static {
	s.Archive = archive
	return s
}

//AppendExts 设置静态文件跟目录
func (s *Static) AppendExts(exts ...string) *Static {
	s.Exts = append(s.Exts, exts...)
	return s
}

//WithPrefix 设置静态文件跟目录
func (s *Static) WithPrefix(prefix string) *Static {
	s.Prefix = prefix
	return s
}

//WithEnable 启用配置
func (s *Static) WithEnable() *Static {
	s.Disable = false
	return s
}

//WithDisable 禁用配置
func (s *Static) WithDisable() *Static {
	s.Disable = false
	return s
}
