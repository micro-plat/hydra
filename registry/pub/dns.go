package pub

import "encoding/json"

//DNSConf dns配置信息
type DNSConf struct {
	PlatName       string `json:"plat_name"`
	PlatCNName     string `json:"plat_cn_name"`
	ClusterName    string `json:"cluster_name"`
	SystemName     string `json:"system_name"`
	SystemCNName   string `json:"system_cn_name"`
	ServerType     string `json:"server_type"`
	ServerName     string `json:"server_name"`
	ServiceAddress string `json:"service_address"`
	Proto          string `json:"proto"`
	Host           string `json:"host"`
	Port           string `json:"port"`
	IPAddress      string `json:"ip"`
}

//GetDNSConf 获取DNS配置信息
func GetDNSConf(buff []byte) (*DNSConf, error) {
	var raw = DNSConf{}
	if err := json.Unmarshal(buff, &raw); err != nil {
		return nil, err
	}
	return &raw, nil
}
