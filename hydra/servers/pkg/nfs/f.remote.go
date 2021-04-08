package nfs

import (
	"fmt"
	"net/http"

	"github.com/micro-plat/hydra"
	"github.com/micro-plat/lib4go/errs"
)

const (

	//获取远程文件的指纹信息
	rpc_fp_get = "/_/nfs/fp/get@%s"

	//推送指纹数据
	rpc_fp_push = "/_/nfs/fp/push@%s"

	//拉取指纹列表
	rpc_fp_list = "/_/nfs/fp/list@%s"

	//获取远程文件数据
	rpc_file_pull = "/_/nfs/file/pull@%s"
)

//remoting 远程文件管理
type remoting struct {
	isMaster   bool
	hosts      []string
	masterHost string
}

func newRemoting(path string, hosts []string, masterHost string, isMaster bool) *remoting {
	return &remoting{
		hosts:      hosts,
		masterHost: masterHost,
		isMaster:   isMaster,
	}
}
func (r *remoting) Update(hosts []string, masterHost string, isMaster bool) {
	r.hosts = hosts
	r.masterHost = masterHost
	r.isMaster = isMaster
}

//GetFP 查询某个文件是否存在，及所在的服务器
func (r *remoting) GetFPFormMaster(name string) (*eFileFP, error) {
	//查询远程服务
	rpns, err := hydra.C.RPC().GetRegularRPC().Request(fmt.Sprintf(rpc_fp_get, r.masterHost), map[string]interface{}{
		"path": name,
	})
	if rpns.GetStatus() == http.StatusNoContent {
		return nil, errs.NewError(http.StatusNotFound, "文件不存在")
	}
	if err != nil {
		return nil, err
	}

	//处理参数合并
	rCache := &eFileFP{}
	if err = rpns.Bind(&rCache); err != nil {
		return nil, err
	}
	return rCache, nil

}

//Pull 从远程服务器拉取文件信息
func (r *remoting) Pull(name string, host []string) ([]byte, error) {
	for _, host := range host {
		//查询远程服务
		rpns, err := hydra.C.RPC().GetRegularRPC().Request(fmt.Sprintf(rpc_file_pull, host), map[string]interface{}{
			"name": name,
		})
		if rpns.GetStatus() == http.StatusNoContent {
			continue
		}
		if err != nil {
			return nil, err
		}
		buff, ok := rpns.GetMap().Get(name)
		if ok {
			return []byte(buff.(string)), nil
		}
	}
	return nil, errs.NewError(http.StatusNoContent, "未找到文件")

}

//Push 向集群推送文件指纹信息
func (r *remoting) Push(fp *eFileFP) error {
	hosts := fp.ExcludeHosts(r.getHosts()...)
	for _, host := range hosts {
		//查询远程服务
		_, err := hydra.C.RPC().GetRegularRPC().Request(fmt.Sprintf(rpc_fp_push, host), fp)
		if err != nil {
			return err
		}
	}
	return nil
}

//Query 向集群机器获取Cache列表,master向所有机器获取,slave启动时向master获取
func (r *remoting) Query() (eFileFPLists, error) {
	//查询数据
	result := make(eFileFPLists)
	for _, host := range r.getHosts() {

		//查询远程服务
		rpns, err := hydra.C.RPC().GetRegularRPC().Request(fmt.Sprintf(rpc_fp_list, host), nil)
		if err != nil {
			return nil, err
		}

		//处理参数合并
		nresult := make(eFileFPLists)
		if err = rpns.Bind(nresult); err != nil {
			return nil, err
		}
		for k, v := range nresult {
			if _, ok := result[k]; !ok {
				result[k] = v
				continue
			}
			result[k].Hosts = append(result[k].Hosts, v.Hosts...)
		}
	}
	return result, nil
}

func (r *remoting) getHosts() []string {
	if r.isMaster {
		return r.hosts
	}
	return []string{r.masterHost}

}
