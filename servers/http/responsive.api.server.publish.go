package http

import (
	"fmt"
	"net"
	"path"
	"strings"
	"time"

	"github.com/micro-plat/hydra/registry"
	"github.com/micro-plat/lib4go/jsons"
)

//publish 将当前服务器的节点信息发布到注册中心
func (w *ApiResponsiveServer) publish() (err error) {
	if err = w.pubServerNode(); err != nil {
		return err
	}
	if err = w.pubServiceNode(); err != nil {
		return
	}
	if err = w.pubDNSNode(); err != nil {
		return
	}
	go w.publishCheck()
	return
}
func (w *ApiResponsiveServer) pubServerNode() error {
	addr := w.server.GetAddress()
	ipPort := strings.Split(addr, "://")[1]
	rname := fmt.Sprintf("%s_%s", ipPort, w.currentConf.GetClusterID())
	pubPath := registry.Join(w.currentConf.GetServerPubRootPath(), rname)
	data := map[string]string{
		"service":    addr,
		"cluster_id": w.currentConf.GetClusterID(),
	}
	jsonData, _ := jsons.Marshal(data)
	nodeData := string(jsonData)
	if b, _ := w.engine.GetRegistry().Exists(pubPath); b {
		w.engine.GetRegistry().Delete(pubPath)
	}
	err := w.engine.GetRegistry().CreateTempNode(pubPath, nodeData)
	if err != nil {
		err = fmt.Errorf("服务发布失败:(%s)[%v]", pubPath, err)
		return err
	}
	w.pubs[pubPath] = nodeData
	return nil
}
func (w *ApiResponsiveServer) pubServiceNode() error {
	addr := w.server.GetAddress(w.currentConf.GetString("dn"))
	ipPort := strings.Split(addr, "://")[1]
	data := map[string]string{
		"service":    addr,
		"cluster_id": w.currentConf.GetClusterID(),
	}
	jsonData, _ := jsons.Marshal(data)
	nodeData := string(jsonData)
	names := w.currentConf.GetStrings("host")
	if len(names) == 0 {
		names = append(names, w.currentConf.GetSysName())
	}
	for _, host := range names {
		servicePath := path.Join(w.currentConf.GetServicePubRootPath(host), ipPort+"_")
		rservicePath, err := w.engine.GetRegistry().CreateSeqNode(servicePath, nodeData)
		if err != nil {
			err = fmt.Errorf("服务发布失败:(%s)[%v]", servicePath, err)
			return err
		}
		w.pubs[rservicePath] = nodeData
	}
	return nil
}
func (w *ApiResponsiveServer) pubDNSNode() error {
	names := w.currentConf.GetStrings("host")
	if len(names) == 0 {
		return nil
	}
	addr := w.server.GetAddress(w.currentConf.GetString("dn"))
	ipPort := strings.Split(addr, "://")[1]
	ip, _, _ := net.SplitHostPort(ipPort)
	data := map[string]string{
		"service":    addr,
		"cluster_id": w.currentConf.GetClusterID(),
	}
	jsonData, _ := jsons.Marshal(data)
	nodeData := string(jsonData)
	for _, host := range names {
		servicePath := path.Join(w.currentConf.GetDNSPubRootPath(host), ip+"_")
		rservicePath, err := w.engine.GetRegistry().CreateSeqNode(servicePath, nodeData)
		if err != nil {
			err = fmt.Errorf("服务发布失败:(%s)[%v]", servicePath, err)
			return err
		}
		w.pubs[rservicePath] = nodeData
	}
	return nil
}

//publishCheck 定时检查节点数据是否存在
func (w *ApiResponsiveServer) publishCheck() {
LOOP:
	for {
		select {
		case <-w.closeChan:
			break LOOP
		case <-time.After(time.Second * 10):
			if w.done {
				break LOOP
			}
			w.checkPubPath()
		}
	}
}

//checkPubPath 检查已发布的节点，不存在则创建
func (w *ApiResponsiveServer) checkPubPath() {
	w.pubLock.Lock()
	defer w.pubLock.Unlock()
	for path, data := range w.pubs {
		if w.done {
			break
		}
		ok, err := w.engine.GetRegistry().Exists(path)
		if err != nil {
			break
		}
		if !ok {
			err := w.engine.GetRegistry().CreateTempNode(path, data)
			if err != nil {
				break
			}
			w.Logger.Infof("节点(%s)已恢复", path)
		}
	}
}

//unpublish 删除已发布的节点
func (w *ApiResponsiveServer) unpublish() {
	w.pubLock.Lock()
	defer w.pubLock.Unlock()
	for path := range w.pubs {
		w.engine.GetRegistry().Delete(path)
	}
	w.pubs = make(map[string]string)
}
