package whitelist

//NewIPList 构建IPLIST
func NewIPList(request []string, ip ...string) *IPList {
	iplist := &IPList{
		Requests: request,
	}
	iplist.IPS = append(iplist.IPS, ip...)
	return iplist
}
