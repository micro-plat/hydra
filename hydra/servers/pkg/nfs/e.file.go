package nfs

import "encoding/json"

type eFileFPLists = map[string]*eFileFP

//eRomotingFileFP 远程文件清单
type eFileFP struct {
	Path  string   `json:"path,omitempty"`
	CRC64 uint64   `json:"crc64,omitempty"`
	Hosts []string `json:"hosts,omitempty"`
}

//eFileEntity 文件实体
type eFileEntity struct {
	Path   string `json:"path,omitempty"`
	CRC64  uint64 `json:"crc64,omitempty"`
	Buffer []byte `json:"buffer,omitempty"`
}

//AddHosts 添加hosts,去除重复host
func (e *eFileFP) AddHosts(hosts ...string) {
	mp := make(map[string]interface{})
	nhost := make([]string, 0, len(hosts)+len(e.Hosts))
	for _, h := range hosts {
		if _, ok := mp[h]; !ok {
			mp[h] = 0
			nhost = append(nhost, h)
		}
	}
	for _, h := range e.Hosts {
		if _, ok := mp[h]; !ok {
			mp[h] = 0
			nhost = append(nhost, h)
		}
	}
	e.Hosts = nhost
}

//GetAliveHost 根据传入的Hosts检查所有hosts,只有在传入列表的才是可用的
func (e *eFileFP) GetAliveHost(hosts ...string) []string {
	mp := make(map[string]interface{})
	nhost := make([]string, 0, len(hosts)+len(e.Hosts))
	for _, h := range hosts {
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

//ExcludeHosts 排除指定的hosts
func (e *eFileFP) ExcludeHosts(hosts ...string) []string {
	mp := make(map[string]interface{})
	nhost := make([]string, 0, len(hosts)+len(e.Hosts))
	for _, h := range hosts {
		if _, ok := mp[h]; !ok {
			mp[h] = 0
		}
	}
	for _, h := range e.Hosts {
		if _, ok := mp[h]; !ok {
			mp[h] = 0
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
