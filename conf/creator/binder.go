package creator

import (
	"fmt"
	"path/filepath"
)

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
	Print()
}
type Binder struct {
	API     IMainBinder
	RPC     IMainBinder
	WEB     IMainBinder
	MQC     IMainBinder
	CRON    IMainBinder
	Plat    IPlatBinder
	binders map[string]IMainBinder
	show    bool
}

func NewBinder() *Binder {
	s := &Binder{}
	s.API = NewMainBinder()
	s.RPC = NewMainBinder()
	s.WEB = NewMainBinder()
	s.MQC = NewMainBinder()
	s.CRON = NewMainBinder()
	s.Plat = NewPlatBinder()
	s.binders = map[string]IMainBinder{
		"api":  s.API,
		"rpc":  s.RPC,
		"web":  s.WEB,
		"mqc":  s.MQC,
		"cron": s.CRON,
	}
	return s
}
func (s *Binder) Print() {
	fmt.Println(s.binders)
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
