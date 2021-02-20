package static

import (
	"errors"
	"fmt"
	"net/http"

	"github.com/micro-plat/hydra/conf"
)

//TempDirName 临时目录创建名
const TempDirName = "hydra"

//TempArchiveName 临时压缩文件创建名
const TempArchiveName = "hydra*"

//TypeNodeName static分类节点名
const TypeNodeName = "static"

//IStatic 静态文件接口
type IStatic interface {
	GetConf() (*Static, bool)
}

//Static 设置静态文件配置
type Static struct {
	Path          string                `json:"path,omitempty" valid:"ascii" label:"静态文件路径或压缩包路径"`
	Excludes      []string              `json:"excludes,omitempty" valid:"ascii" label:"排除名称"`
	HomePage      string                `json:"homePath,omitempty" valid:"ascii" label:"静态文件首页"`
	AutoRewrite   bool                  `json:"autoRewrite,omitempty" valid:"ascii" label:"自动重写到首页"`
	Disable       bool                  `json:"disable,omitempty"`
	excludesMatch *conf.PathMatch       `json:"-"`
	fs            IFS                   `json:"-"`
	gzipfileMap   map[string]gzFileInfo `json:"-"`
}

//New 构建静态文件配置信息
func New(opts ...Option) *Static {
	a := &Static{HomePage: DefaultHome, Excludes: DefaultExclude, gzipfileMap: map[string]gzFileInfo{}}
	for _, opt := range opts {
		opt(a)
	}

	return a
}

//Get 获取文件内容//http.FileServer(http.FS(embed.FS{}))
func (s *Static) Get(name string) (http.FileSystem, string, error) {
	if s.fs == nil {
		return nil, "", nil
	}
	//排除内容
	if s.IsExclude(name) {
		return nil, "", nil
	}

	if !s.fs.Has(name) && s.AutoRewrite {
		return s.fs.ReadFile(s.HomePage)
	}
	return s.fs.ReadFile(name)
}

//IsExclude 是否是排除的文件
func (s *Static) IsExclude(rPath string) bool {
	if len(s.Excludes) == 0 {
		return false
	}
	ok, _ := s.excludesMatch.Match(rPath)
	return ok
}

//AllowRequest 是否是合适的请求
func (s *Static) AllowRequest(m string) bool {
	return m == http.MethodGet || m == http.MethodHead
}

//GetConf 设置static
func GetConf(cnf conf.IServerConf) (*Static, error) {
	static := New()
	_, err := cnf.GetSubObject(TypeNodeName, static)
	if err != nil && !errors.Is(err, conf.ErrNoSetting) {
		return nil, fmt.Errorf("static配置格式有误:%v", err)
	}
	static.excludesMatch = conf.NewPathMatch(static.Excludes...)

	//转换配置文件
	fs, err := static.check2fs()
	if err != nil {
		return nil, err
	}
	if fs != nil {
		static.fs = NewGzip(fs, static)
		return static, nil
	}

	//转换本地内嵌文件
	fs, err = defEmbedFs.check2FS()
	if err != nil {
		return nil, err
	}
	if fs != nil {
		static.fs = NewGzip(fs, static)
		return static, nil
	}
	return nil, conf.ErrNoSetting
}
