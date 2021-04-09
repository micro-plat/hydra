package nfs

import (
	"encoding/json"
)

type eFileFPLists map[string]*eFileFP

func (ex eFileFPLists) GetJSON() string {
	buff, _ := json.Marshal(ex)
	return string(buff)
}

//eRomotingFileFP 远程文件清单
type eFileFP struct {
	Path  string   `json:"path,omitempty" valid:"required" `
	CRC64 uint64   `json:"crc64,omitempty" valid:"required" `
	Hosts []string `json:"hosts,omitempty" valid:"required" `
}

//eFileEntity 文件实体
type eFileEntity struct {
	Path   string `json:"path,omitempty" valid:"required"`
	CRC64  uint64 `json:"crc64,omitempty" valid:"required"`
	Buffer []byte `json:"buffer,omitempty" valid:"required"`
}

//AddHosts 添加hosts,去除重复host
func (e *eFileFP) AddHosts(hosts ...string) bool {
	hasChange := false
	mp := make(map[string]interface{})
	nhost := make([]string, 0, len(hosts)+len(e.Hosts))
	for _, h := range e.Hosts {
		if _, ok := mp[h]; !ok && h != "" {
			mp[h] = 0
			nhost = append(nhost, h)
		}
	}

	for _, h := range hosts {
		if _, ok := mp[h]; !ok && h != "" {
			mp[h] = 0
			nhost = append(nhost, h)
			hasChange = true
		}
	}

	e.Hosts = nhost
	return hasChange
}

//GetAliveHost 根据传入的Hosts检查所有hosts,只有在传入列表的才是可用的
func (e *eFileFP) GetAliveHost(aliveHosts ...string) []string {
	mp := make(map[string]interface{})
	nhost := make([]string, 0, len(aliveHosts)+len(e.Hosts))
	for _, h := range aliveHosts {
		if _, ok := mp[h]; !ok {
			mp[h] = 0
		}
	}
	for _, h := range e.Hosts {
		if _, ok := mp[h]; ok {
			nhost = append(nhost, h)
		}
	}
	return nhost
}

func (e *eFileFP) Has(host string) bool {
	for _, h := range e.Hosts {
		if h == host {
			return true
		}
	}
	return false
}
func (e *eFileFP) GetJSON() string {
	buff, _ := json.Marshal(e)
	return string(buff)
}
func (e *eFileFP) String() string {
	return e.GetJSON()
}
func (e *eFileFP) GetMAP() eFileFPLists {
	return map[string]*eFileFP{
		e.Path: e,
	}
}
