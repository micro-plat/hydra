package rpc

import (
	"fmt"
	"strings"
)

//ResolvePath   解析注册中心地址
//domain:hydra,server:merchant_cron
//order.request@merchant_api.hydra 解析为:service: /order/request,server:merchant_api,domain:hydra
//order.request 解析为 service: /order/request,server:merchant_cron,domain:hydra
//order.request@merchant_rpc 解析为 service: /order/request,server:merchant_rpc,domain:hydra
func ResolvePath(address string, d string, s string) (service string, domain string, server string, err error) {
	raddress := strings.TrimRight(address, "@")
	addrs := strings.SplitN(raddress, "@", 2)
	if len(addrs) == 1 {
		if addrs[0] == "" {
			return "", "", "", fmt.Errorf("服务地址%s不能为空", address)
		}
		service = "/" + strings.Trim(strings.Replace(raddress, ".", "/", -1), "/")
		domain = d
		server = s
		return
	}
	if addrs[0] == "" {
		return "", "", "", fmt.Errorf("%s错误，服务名不能为空", address)
	}
	if addrs[1] == "" {
		return "", "", "", fmt.Errorf("%s错误，服务名，域不能为空", address)
	}
	service = "/" + strings.Trim(strings.Replace(addrs[0], ".", "/", -1), "/")
	raddr := strings.Split(strings.TrimRight(addrs[1], "."), ".")
	if len(raddr) >= 2 && raddr[0] != "" && raddr[1] != "" {
		domain = raddr[len(raddr)-1]
		server = strings.Join(raddr[0:len(raddr)-1], ".")
		return
	}
	if len(raddr) == 1 {
		if raddr[0] == "" {
			return "", "", "", fmt.Errorf("%s错误，服务器名称不能为空", address)
		}
		domain = d
		server = raddr[0]
		return
	}
	if raddr[0] == "" && raddr[1] == "" {
		return "", "", "", fmt.Errorf(`%s错误,未指定服务器名称和域名称`, addrs[1])
	}
	domain = raddr[1]
	server = s
	return
}

// func resolvePath(addr string, serverName string, platName string) (raddr string, faddr string, sName string, pName string, err error) {
// 	//拆分服务名，服务器名，平台名
// 	blocks := getBlocks(addr)
// 	if len(blocks) != 3 {
// 		err = fmt.Errorf("地址格式错误:%s", addr)
// 		return
// 	}
// 	//处理未匹配到服务名，平台名根据传入值进行设置
// 	sName = types.GetString(blocks[1], serverName)
// 	pName = types.GetString(blocks[2], platName)
// 	b, names := needCheckNames(blocks[0])
// 	if !b {
// 		faddr = blocks[0]
// 		raddr = blocks[0]
// 		return
// 	}

// 	//解析服务名为
// 	var raddrs = make([]string, 0, 3)
// 	var faddrs = make([]string, 0, 3)
// 	for _, name := range names {
// 		cnames := getNames(name)
// 		switch len(cnames) {
// 		case 1:
// 			faddrs = append(faddrs, cnames[0])
// 			raddrs = append(raddrs, cnames[0])
// 		case 2:
// 			faddrs = append(faddrs, cnames[0])
// 			raddrs = append(raddrs, cnames[1])
// 		}
// 	}
// 	faddr = "/" + strings.Join(faddrs, "/")
// 	raddr = "/" + strings.Join(raddrs, "/")
// 	return

// }

// func needCheckNames(request string) (bool, []string) {
// 	if !strings.Contains(request, ":") {
// 		return false, nil
// 	}
// 	lst := strings.Split(strings.Trim(request, "/"), "/")
// 	return true, lst
// }

// func getNames(addr string) []string {
// 	brackets := regexp.MustCompile(`^(:\w+)\[(\w+)\]$`)
// 	result := brackets.FindStringSubmatch(addr)
// 	if len(result) > 0 {
// 		return result[1:]
// 	}
// 	result = regexp.MustCompile(`^(\w+)$`).FindStringSubmatch(addr)
// 	if len(result) > 0 {
// 		return result[1:]
// 	}
// 	return nil
// }

// func getBlocks(addr string) []string {
// 	brackets := regexp.MustCompile(`^([^@]+)[@]?([\w]*)[.]?([\w]*)$`)
// 	result := brackets.FindStringSubmatch(addr)
// 	if len(result) == 0 {
// 		return result
// 	}
// 	return result[1:]
// }
