package rpc

import (
	"fmt"

	"github.com/micro-plat/hydra/conf"
	"github.com/micro-plat/hydra/servers"
)

//Notify 服务器配置变更通知
func (w *RpcResponsiveServer) Notify(nConf conf.IServerConf) error {
	w.mu.Lock()
	defer w.mu.Unlock()
	w.restarted = false
	//检查是否需要重启服务器
	restart, err := w.NeedRestart(nConf)
	if err != nil {
		return err
	}
	if restart { //服务器地址已变化，则重新启动新的server,并停止当前server
		servers.Trace(w.Infof, "关键配置发生变化，准备重启服务器")
		return w.Restart(nConf)
	}

	servers.Trace(w.Infof, "配置发生变化，准备更新")
	//服务器地址未变化，更新服务器当前配置，并立即生效
	if err = w.SetConf(false, nConf); err != nil {
		return err
	}
	w.engine.UpdateVarConf(nConf)
	w.currentConf = nConf
	return nil
}

//NeedRestart 检查配置判断是否需要重启服务器
func (w *RpcResponsiveServer) NeedRestart(cnf conf.IServerConf) (bool, error) {
	if cnf.ForceRestart() {
		return true, nil
	}
	comparer := conf.NewComparer(w.currentConf, cnf)
	if !comparer.IsChanged() {
		return false, nil
	}
	if comparer.IsVarChanged() {
		return true, nil
	}
	if comparer.IsValueChanged("status", "address", "host") {
		return true, nil
	}
	ok, err := comparer.IsRequiredSubConfChanged("router")
	if ok {
		return ok, nil
	}
	if err != nil {
		return ok, fmt.Errorf("路由未配置或配置有误:%s(%+v)", cnf.GetServerName(), err)
	}
	if ok := comparer.IsSubConfChanged("app"); ok {
		return ok, nil
	}
	return false, nil

}

//SetConf 设置配置参数
func (w *RpcResponsiveServer) SetConf(restart bool, conf conf.IServerConf) (err error) {

	var ok bool
	//设置路由
	if restart {
		if _, err = SetRouters(w.engine, conf, w.server, nil); err != nil {
			return err
		}
	}
	if err != nil {
		return
	}
	//设置jwt安全认证
	if ok, err = SetJWT(w.server, conf); err != nil {
		return err
	}
	servers.TraceIf(ok, w.Infof, w.Debugf, getEnableName(ok), "jwt设置")

	//设置请求头
	if ok, err = SetHeaders(w.server, conf); err != nil {
		return err
	}
	servers.TraceIf(ok, w.Infof, w.Debugf, getEnableName(ok), "header设置")

	//设置metric
	if ok, err = SetMetric(w.server, conf); err != nil {
		return err
	}
	servers.TraceIf(ok, w.Infof, w.Debugf, getEnableName(ok), "metric设置")

	//设置host
	if ok, err = SetHosts(w.server, conf); err != nil {
		return err
	}
	servers.TraceIf(ok, w.Infof, w.Debugf, getEnableName(ok), "host设置")
	return nil
}
func getEnableName(b bool) string {
	if b {
		return "启用"
	}
	return "未启用"
}
