package nfs

import (
	"time"

	"github.com/micro-plat/lib4go/types"
)

type eFileFPLists map[string]*eFileFP

func (e eFileFPLists) Merge(list eFileFPLists) {
	for k, v := range list {
		if _, ok := e[k]; !ok {
			e[k] = v
			continue
		}
		e[k].MergeHosts(v.Hosts...)
	}
}

//GetAlives 获取可用于通知所有alive的tp列表
func (e eFileFPLists) GetAlives(allAliveHosts []string) map[string]eFileFPLists {
	if len(e) == 0 {
		return nil
	}
	list := make(map[string]eFileFPLists, len(allAliveHosts))
	for _, v := range allAliveHosts {
		list[v] = e
	}
	return list
}

//eRomotingFileFP 远程文件清单
type eFileFP struct {
	Path string `json:"path,omitempty" valid:"required" `

	ModTime time.Time `json:"modTime,omitempty" valid:"required"`

	// CRC64 uint64   `json:"crc64,omitempty" valid:"required" `
	Hosts []string `json:"hosts,omitempty" valid:"required" `

	//文件大小
	Size int64 `json:"size,omitempty"`
}

//eFileEntity 文件实体
type eFileEntity struct {
	Path    string    `json:"path,omitempty" valid:"required"`
	ModTime time.Time `json:"modTime,omitempty" valid:"required"`
	Size    int64     `json:"buffer,omitempty" valid:"required"`
}

type eFileEntityList []*eFileEntity

func (l eFileEntityList) GetMap() map[string]*eFileEntity {
	mp := make(map[string]*eFileEntity)
	for _, v := range l {
		mp[v.Path] = v
	}
	return mp
}

//MergeHosts 合并hosts
func (e *eFileFP) MergeHosts(hosts ...string) bool {
	hasChange := false
	mp := make(map[string]interface{})
	nhost := make([]string, 0, len(hosts)+len(e.Hosts))
	for _, h := range e.Hosts {
		if _, ok := mp[h]; !ok && h != "" && h != ":" {
			mp[h] = 0
			nhost = append(nhost, h)
		}
	}

	for _, h := range hosts {
		if _, ok := mp[h]; !ok && h != "" && h != ":" {
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
		if _, ok := mp[h]; ok && h != "" {
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
func (e *eFileFP) String() string {
	return types.ToJSON(e)
}
func (e *eFileFP) GetMAP() eFileFPLists {
	return map[string]*eFileFP{
		e.Path: e,
	}
}
