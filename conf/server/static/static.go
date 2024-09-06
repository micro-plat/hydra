package static

import (
	"errors"
	"fmt"
	"net/http"
	"strings"

	"github.com/micro-plat/hydra/conf"
	"github.com/micro-plat/hydra/conf/pkgs/security"
)

// TempDirName 临时目录创建名
const TempDirName = "hydra"

// TempArchiveName 临时压缩文件创建名
const TempArchiveName = "hydra*"

// TypeNodeName static分类节点名
const TypeNodeName = "static"

// IStatic 静态文件接口
type IStatic interface {
	GetConf() (*Static, bool)
}

// Static 设置静态文件配置
type Static struct {
	security.ConfEncrypt
	Path           string          `json:"path,omitempty"  label:"静态文件路径或压缩包路径"`
	Excludes       []string        `json:"excludes,omitempty"  label:"排除名称"`
	HomePage       string          `json:"homePath,omitempty" valid:"ascii" label:"静态文件首页"`
	AutoRewrite    bool            `json:"autoRewrite,omitempty"  label:"自动重写到首页"`
	Unrewrites     []string        `json:"unrewrite,omitempty"  label:"不重写列表"`
	Disable        bool            `json:"disable,omitempty"`
	unrewriteMatch *conf.PathMatch `json:"-"`
	fs             IFS             `json:"-"`
	serverType     string
}

// New 构建静态文件配置信息
func New(serverType string, opts ...Option) *Static {
	a := &Static{
		serverType: serverType,
		HomePage:   DefaultHome,
		Excludes:   DefaultExclude,
		Unrewrites: DefaultUnrewrite,
	}
	for _, opt := range opts {
		opt(a)
	}

	return a
}

// Has 检查文件是否存在
func (s *Static) Has(name string) bool {
	if s.fs == nil {
		return false
	}
	//排除内容
	if s.IsExclude(name) {
		return false
	}
	if _, ok := s.fs.Has(name); ok {
		return true
	}
	if s.AutoRewrite && !s.IsUnrewrite(name) {
		return true
	}
	return false
}

// OptionsCheck 检查文件是否存在
func (s *Static) OptionsCheck(name string) bool {
	if s.fs == nil {
		return false
	}
	//排除内容
	if s.IsExclude(name) {
		return false
	}
	if _, ok := s.fs.Has(name); ok {
		return true
	}
	return false
}

// Get 获取文件内容//http.FileServer(http.FS(embed.FS{}))
func (s *Static) Get(name string) (http.FileSystem, string, error) {
	if s.fs == nil {
		return nil, "", nil
	}
	//排除内容
	if s.IsExclude(name) {
		return nil, "", nil
	}

	//文件不存在
	if _, ok := s.fs.Has(name); !ok {
		//是否是不重写文件
		if s.IsUnrewrite(name) {
			return nil, "", nil
		}
		//是否自动重写
		if s.AutoRewrite {
			return s.fs.ReadFile(s.HomePage)
		}
		return nil, "", nil
	}
	return s.fs.ReadFile(name)
}

// IsExclude 是否是排除的文件
// @todo:该方法不是非常合适，需要修改匹配算法
func (s *Static) IsExclude(rPath string) bool {
	if len(s.Excludes) == 0 {
		return false
	}

	for i := range s.Excludes {
		if strings.Contains(rPath, s.Excludes[i]) {
			return true
		}
	}

	return false
}

// IsUnrewrite 是否是非重写文件
func (s *Static) IsUnrewrite(rPath string) bool {
	if len(s.Unrewrites) == 0 {
		return false
	}
	ok, _ := s.unrewriteMatch.Match(rPath)
	return ok
}

// AllowRequest 是否是合适的请求
func (s *Static) AllowRequest(m string) bool {
	return m == http.MethodGet || m == http.MethodHead
}

// GetConf 设置static
func GetConf(cnf conf.IServerConf) (*Static, error) {
	static := New(cnf.GetServerType())
	_, err := cnf.GetSubObject(TypeNodeName, static)
	if err != nil {
		if errors.Is(err, conf.ErrNoSetting) {
			static.Disable = true
			return static, nil
		}
		return nil, fmt.Errorf("static配置格式有误:%v", err)
	}

	//禁用靜态文件配置
	if static.Disable {
		return static, nil
	}
	static.unrewriteMatch = conf.NewPathMatch(static.Unrewrites...)

	//转换配置文件
	fs, err := static.getFileOS()
	if err != nil {
		return nil, err
	}

	//转换本地内嵌文件
	if nfs, ok := defEmbedFs[static.serverType]; ok {
		efs, err := nfs.getFileEmbed()
		if err != nil {
			return nil, err
		}
		efs.Merge(fs)
		fs = efs
	}
	if fs != nil {
		static.fs = fs
		return static, nil
	}
	return nil, fmt.Errorf("%s %w", "static", conf.ErrNoSetting)
}
