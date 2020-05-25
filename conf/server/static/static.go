package static

import (
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"

	"github.com/asaskevich/govalidator"
	"github.com/micro-plat/hydra/conf"
	"github.com/micro-plat/lib4go/archiver"
)

//IStatic 静态文件接口
type IStatic interface {
	GetConf() (*Static, bool)
}

//Static 设置静态文件配置
type Static struct {
	*option
}

//New 构建静态文件配置信息
func New(opts ...Option) *Static {
	s := &Static{option: newOption()}
	for _, opt := range opts {
		opt(s.option)
	}
	return s
}

//AllowRequest 是否是合适的请求
func (s *Static) AllowRequest(m string) bool {
	return m == "GET" || m == "HEAD"
}

//GetConf 设置static
func GetConf(cnf conf.IMainConf) (static *Static) {
	//设置静态文件路由
	_, err := cnf.GetSubObject("static", &static)
	if err != nil && err != conf.ErrNoSetting {
		panic(fmt.Errorf("static配置有误:%v", err))
	}
	if err == conf.ErrNoSetting {
		static = New(WithDisable())
		return
	}
	if b, err := govalidator.ValidateStruct(&static); !b {
		panic(fmt.Errorf("static配置有误:%v", err))
	}
	static.Dir, err = unarchive(static.Dir, static.Archive) //处理归档文件
	return
}

var waitRemoveDir = make([]string, 0, 1)

func unarchive(dir string, path string) (string, error) {
	if path == "" {
		return dir, nil
	}
	archive := archiver.MatchingFormat(path)
	if archive == nil {
		return "", fmt.Errorf("指定的文件不是归档文件:%s", path)
	}
	tmpDir, err := ioutil.TempDir("", "hydra")
	if err != nil {
		return "", fmt.Errorf("创建临时文件失败:%v", err)
	}
	reader, err := os.Open(path)
	if err != nil {
		return "", fmt.Errorf("无法打开文件:%s(%v)", path, err)
	}
	defer reader.Close()
	err = archive.Read(reader, tmpDir)
	if err != nil {
		return "", fmt.Errorf("读取归档文件失败:%v", err)
	}
	ndir := filepath.Join(tmpDir, dir)
	waitRemoveDir = append(waitRemoveDir, tmpDir)
	return ndir, nil
}
