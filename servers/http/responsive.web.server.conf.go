package http

import (
	"fmt"

	"github.com/micro-plat/hydra/conf"
	"github.com/micro-plat/hydra/servers"
)

//Notify 服务器配置变更通知
func (w *WebResponsiveServer) Notify(conf conf.IServerConf) error {
	w.mu.Lock()
	defer w.mu.Unlock()
	w.restarted = false
	//检查是否需要重启服务器
	restart, err := w.NeedRestart(conf)
	if err != nil {
		return err
	}
	if restart { //服务器地址已变化，则重新启动新的server,并停止当前server
		servers.Trace(w.Infof, "关键配置发生变化，准备重启服务器")
		return w.Restart(conf)
	}
	//服务器地址未变化，更新服务器当前配置，并立即生效
	servers.Trace(w.Infof, "配置发生变化，准备更新")
	if err = w.SetConf(false, conf); err != nil {
		return err
	}
	w.engine.UpdateVarConf(conf)
	w.currentConf = conf
	return nil
}

//NeedRestart 检查配置判断是否需要重启服务器
func (w *WebResponsiveServer) NeedRestart(cnf conf.IServerConf) (bool, error) {
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
	if comparer.IsValueChanged("status", "address", "host", "rTimeout", "wTimeout", "rhTimeout") {
		return true, nil
	}
	ok, err := comparer.IsRequiredSubConfChanged("router")
	if ok {
		return ok, nil
	}
	if err != nil {
		return ok, fmt.Errorf("路由未配置或配置有误:%v", err)
	}
	if ok := comparer.IsSubConfChanged("app"); ok {
		return ok, nil
	}
	if ok := comparer.IsSubConfChanged("view"); ok {
		return ok, nil
	}
	if ok := comparer.IsSubConfChanged("circuit"); ok {
		return ok, nil
	}
	return false, nil
}

//SetConf 设置配置参数
func (w *WebResponsiveServer) SetConf(restart bool, conf conf.IServerConf) (err error) {
	if err = w.ApiResponsiveServer.SetConf(restart, conf); err != nil {
		return err
	}
	//设置metric
	var ok bool
	if ok, err = SetView(w.webServer, conf); err != nil {
		return err
	}
	servers.TraceIf(ok, w.Infof, w.Debugf, getEnableName(ok), "view设置")
	return nil
}
