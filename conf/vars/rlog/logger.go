package rlog

import (
	"encoding/json"
	"errors"
	"fmt"

	"github.com/asaskevich/govalidator"
	"github.com/micro-plat/hydra/conf"
	"github.com/micro-plat/hydra/registry"
	"github.com/micro-plat/lib4go/logger"
)

//LogName rlog 的日志名
const LogName = "rlog"

//TypeNodeName 分类节点名
const TypeNodeName = "app"

//Layout 日志配置
type Layout struct {
	Level   string `json:"level"  valid:"in(Off|Info|Warn|Error|Fatal|Debug|All)" toml:"level"`
	Service string `json:"service,omitempty" toml:"service"`
	Layout  string `json:"layout" toml:"layout"`
	Disable bool   `json:"disable,omitempty" toml:"disable,omitempty"`
}

const DefaultLayout = `{"server-ip":"%ip","time":"%datetime.%ms","level":"%level","session":"%session","content":"%content"}`

//New 初始化远程日志组件
func New(service string, opts ...Option) *Layout {
	l := &Layout{
		Service: service,
		Layout:  DefaultLayout,
		Level:   "Info",
	}
	for _, opt := range opts {
		opt(l)
	}
	return l
}

//ToLoggerLayout 转换为logger.Layout
func (l *Layout) ToLoggerLayout() *logger.Layout {
	return &logger.Layout{
		Type:   LogName,
		Level:  l.Level,
		Path:   l.Service,
		Layout: l.Layout,
	}
}

//GetConfByAddr 获取日志配置
func GetConfByAddr(r registry.IRegistry, platName string) (s *Layout, err error) {
	path := registry.Join(platName, "var", TypeNodeName, LogName)
	s = &Layout{}
	ok, err := r.Exists(path)
	if err != nil {
		return nil, fmt.Errorf("检查日志配置出错 %s %w", path, err)
	}
	if !ok {
		s.Disable = true
		return s, nil
	}

	buff, _, err := r.GetValue(path)
	if err != nil {
		return nil, fmt.Errorf("获取日志配置出错 %s %w", path, err)
	}
	if err := json.Unmarshal(buff, s); err != nil {
		err = fmt.Errorf("远程日志日志配置出错 %s %v", path, err)
		return nil, err
	}
	if b, err := govalidator.ValidateStruct(s); !b {
		panic(fmt.Errorf("./var/logger/rlog配置有误 %w", err))
	}
	return s, nil
}

//GetConf 获取主配置信息
func GetConf(cnf conf.IVarConf) (s *Layout, err error) {
	s = &Layout{
		Layout: DefaultLayout,
	}
	_, err = cnf.GetObject(TypeNodeName, LogName, s)
	if err != nil && !errors.Is(err, conf.ErrNoSetting) {
		return nil, fmt.Errorf("读取./var/%s/%s 配置发生错误 %w", TypeNodeName, LogName, err)
	}
	if errors.Is(err, conf.ErrNoSetting) {
		s.Disable = true
		return s, nil
	}
	if b, err := govalidator.ValidateStruct(s); !b {
		return nil, fmt.Errorf("./var/%s/%s 配置有误 %w", TypeNodeName, LogName, err)
	}
	return s, nil
}
