package creator

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"path/filepath"
	"strings"

	"github.com/micro-plat/hydra/component"
)

type Input struct {
	Name    string
	Desc    string
	Filters []func(string) (string, error)
}
type IBinder interface {
	GetMainConfNames(platName string, systemName string, tp string, clusterName string) []string
	GetSubConfNames(serverType string) []string
	GetVarConfNames() []string
	ScanMainConf(mainPath string, serverType string) error
	ScanSubConf(mainPath string, serverType string, subName string) error
	ScanVarConf(platName string, nodeName string) error

	GetMainConf(serverType string) string
	GetSubConf(serverType string, subName string) string
	GetVarConf(nodeName string) string

	GetMainConfScanNum(serverType string) int
	GetSubConfScanNum(serverType string, subName string) int
	GetVarConfScanNum(nodeName string) int
	GetInstallers(serverType string) []func(c component.IContainer) error
	GetSQL(dir string) ([]string, error)
	GetInput() map[string]*Input
	SetParam(k, v string)
	Print()
}
type Binder struct {
	API     *MainBinder
	RPC     *MainBinder
	WS      *MainBinder
	WEB     *MainBinder
	MQC     *MainBinder
	CRON    *MainBinder
	Plat    IPlatBinder
	binders map[string]*MainBinder
	params  map[string]string
	input   map[string]*Input
	show    bool
}

func NewBinder() *Binder {
	s := &Binder{params: make(map[string]string), input: make(map[string]*Input)}
	s.API = NewMainBinder(s.params)
	s.RPC = NewMainBinder(s.params)
	s.WS = NewMainBinder(s.params)
	s.WEB = NewMainBinder(s.params)
	s.MQC = NewMainBinder(s.params)
	s.CRON = NewMainBinder(s.params)
	s.Plat = NewPlatBinder(s.params)
	s.binders = map[string]*MainBinder{
		"api":  s.API,
		"rpc":  s.RPC,
		"web":  s.WEB,
		"mqc":  s.MQC,
		"cron": s.CRON,
		"ws":   s.WS,
	}
	return s
}
func (s *Binder) Print() {
	fmt.Println(s.binders)
}
func (s *Binder) SetParam(k, v string) {
	s.params[k] = v
}
func (s *Binder) GetInput() map[string]*Input {
	return s.input
}
func (s *Binder) SetInput(key, desc string, filters ...func(v string) (string, error)) {
	s.input[key] = &Input{
		Name:    key,
		Desc:    desc,
		Filters: filters,
	}
}
func (s *Binder) GetInstallers(serverType string) []func(c component.IContainer) error {
	return s.binders[serverType].GetInstallers()
}

//GetMainConfNames 获取已配置的主配置名称
func (s *Binder) GetMainConfNames(platName string, systemName string, tp string, clusterName string) []string {
	names := make([]string, 0, 1)
	//	binder := s.binders[tp]
	//	if v := binder.NeedScanCount(""); v > 0 {
	names = append(names, filepath.Join("/", platName, systemName, tp, clusterName, "conf"))
	//	}
	return names
}

//GetSubConfNames 获取已配置的主配置名称
func (s *Binder) GetSubConfNames(serverType string) []string {
	binder := s.binders[serverType]
	return binder.GetSubConfNames()
}

//GetVarConfNames 获取已配置的主配置名称
func (s *Binder) GetVarConfNames() []string {
	return s.Plat.GetVarNames()
}

//GetMainConfScanNum 获取主配置待扫描参数个数
func (s *Binder) GetMainConfScanNum(serverType string) int {
	binder := s.binders[serverType]
	return binder.NeedScanCount("")
}

//GetSubConfScanNum 获取子配置待扫描参数个数
func (s *Binder) GetSubConfScanNum(serverType string, subName string) int {
	binder := s.binders[serverType]
	return binder.NeedScanCount(subName)
}

//GetVarConfScanNum 获取var配置待扫描参数个数
func (s *Binder) GetVarConfScanNum(nodeName string) int {
	return s.Plat.NeedScanCount(nodeName)
}

//ScanMainConf 扫描主配置
func (s *Binder) ScanMainConf(mainPath string, serverType string) error {
	binder := s.binders[serverType]
	return binder.Scan(mainPath, "")
}

//ScanSubConf 扫描子配置
func (s *Binder) ScanSubConf(mainPath string, serverType string, subName string) error {
	binder := s.binders[serverType]
	return binder.Scan(mainPath, subName)
}

//ScanVarConf 扫描平台配置
func (s *Binder) ScanVarConf(platName string, nodeName string) error {
	return s.Plat.Scan(platName, nodeName)
}

//GetMainConf 获取主配置信息
func (s *Binder) GetMainConf(serverType string) string {
	binder := s.binders[serverType]
	return binder.GetNodeConf("")
}

//GetSubConf 获取子配置信息
func (s *Binder) GetSubConf(serverType string, subName string) string {
	binder := s.binders[serverType]
	return binder.GetNodeConf(subName)
}

//GetVarConf 获取平台配置信息
func (s *Binder) GetVarConf(nodeName string) string {
	return s.Plat.GetNodeConf(nodeName)
}

//GetSQL 获取指定目录下所有.sql文件中的SQL语句，并用分号拆分
func (s *Binder) GetSQL(dir string) ([]string, error) {
	files, err := filepath.Glob(filepath.Join(dir, "*.sql"))
	if err != nil {
		return nil, err
	}
	buff := bytes.NewBufferString("")
	for _, f := range files {
		buf, err := ioutil.ReadFile(f)
		if err != nil {
			return nil, err
		}
		_, err = buff.Write(buf)
		if err != nil {
			return nil, err
		}
		buff.WriteString(";")
	}
	tables := make([]string, 0, 8)
	tbs := strings.Split(buff.String(), ";")
	for _, t := range tbs {
		if tb := strings.TrimSpace(t); len(tb) > 0 {
			tables = append(tables, Translate(tb, s.params))
		}
	}
	tables = append(tables, Translate(buff.String(), s.params))
	return tables, nil
}
