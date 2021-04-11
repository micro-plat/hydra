package nfs

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/micro-plat/hydra"
	"github.com/micro-plat/lib4go/errs"
	"github.com/micro-plat/lib4go/types"
)

const (

	//获取远程文件的指纹信息
	rmt_fp_get = "/_/nfs/fp/get"

	//推送指纹数据
	rmt_fp_push = "/_/nfs/fp/push"

	//拉取指纹列表
	rmt_fp_list = "/_/nfs/fp/list"

	//获取远程文件数据
	rmt_file_pull = "/_/nfs/file/pull"
)

//remoting 远程文件管理
type remoting struct {
	isMaster    bool
	hosts       []string
	masterHost  string
	currentAddr string
}

func newRemoting(path string, hosts []string, masterHost string, currentAddr string, isMaster bool) *remoting {
	return &remoting{
		hosts:       hosts,
		masterHost:  masterHost,
		currentAddr: currentAddr,
		isMaster:    isMaster,
	}
}
func (r *remoting) Update(hosts []string, masterHost string, currentAddrs string, isMaster bool) {
	r.hosts = hosts
	r.masterHost = masterHost
	r.currentAddr = currentAddrs
	r.isMaster = isMaster
}

//GetFPFormMaster 查询某个文件是否存在，及所在的服务器
func (r *remoting) GetFPFormMaster(name string) (v *eFileFP, err error) {

	//查询远程服务
	v = &eFileFP{}
	input := types.XMap{"name": name}
	log := start(rmt_fp_get, name, r.masterHost)
	rpns, status, err := hydra.C.HTTP().GetRegularClient().Post(fmt.Sprintf("http://%s%s", r.masterHost, rmt_fp_get), input.ToKV())
	log.end(rmt_fp_get, name, r.masterHost, status, err)
	if status == http.StatusNoContent {
		return nil, errs.NewError(http.StatusNotFound, "文件不存在")
	}
	if err != nil {
		return nil, err
	}

	//处理参数合并
	if err = json.Unmarshal([]byte(rpns), v); err != nil {
		return nil, err
	}
	return v, nil
}

//Pull 从远程服务器拉取文件信息
func (r *remoting) Pull(name string, host []string) ([]byte, error) {
	input := types.XMap{"name": name}
	for _, host := range host {
		//查询远程服务
		log := start(rmt_file_pull, name, "to", host)
		rpns, status, err := hydra.C.HTTP().GetRegularClient().Post(fmt.Sprintf("http://%s%s", host, rmt_file_pull), input.ToKV())
		log.end(rmt_file_pull, name, status, err)
		if status == http.StatusNoContent {
			continue
		}
		if err != nil {
			return nil, err
		}
		return []byte(rpns), nil
	}
	return nil, errs.NewError(http.StatusNoContent, "未找到文件")
}

//Push 向集群推送文件指纹信息
func (r *remoting) Push(fp eFileFPLists) error {
	hosts := fexclude(r.currentAddr, r.getHosts()...)
	for _, host := range hosts {
		//查询远程服务
		log := start(rmt_fp_push, fp, "to", host)
		_, status, err := hydra.C.HTTP().GetRegularClient().Request("POST",
			fmt.Sprintf("http://%s%s", host, rmt_fp_push), fp.GetJSON(), "utf-8", http.Header{
				"Content-Type": []string{"application/json"},
			})
		log.end(rmt_fp_push, host, status, err)
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
		log := start(rmt_fp_list, "from", host)
		rpns, status, err := hydra.C.HTTP().GetRegularClient().Post(fmt.Sprintf("http://%s%s", host, rmt_fp_list), "")
		log.end(rmt_fp_list, "from", host, status, err)
		if err != nil {
			return nil, err
		}

		//处理参数合并
		nresult := make(eFileFPLists)
		if err = json.Unmarshal([]byte(rpns), &nresult); err != nil {
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
func fexclude(ex string, hosts ...string) []string {
	mp := make(map[string]interface{})
	nhost := make([]string, 0, len(hosts))
	for _, h := range hosts {
		if _, ok := mp[h]; !ok {
			mp[h] = 0
		}
	}
	for h := range mp {
		if h != ex {
			nhost = append(nhost, h)
		}
	}
	return nhost
}
