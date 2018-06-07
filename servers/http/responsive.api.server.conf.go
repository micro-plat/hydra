package http

import (
	"fmt"

	"github.com/micro-plat/hydra/conf"
	"github.com/micro-plat/hydra/servers"
)

//Notify 服务器配置变更通知
func (w *ApiResponsiveServer) Notify(conf conf.IServerConf) error {
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
func (w *ApiResponsiveServer) NeedRestart(cnf conf.IServerConf) (bool, error) {
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
		return true, nil
	}
	if err != nil {
		return false, fmt.Errorf("路由未配置或配置有误:%v", err)
	}
	if ok := comparer.IsSubConfChanged("app"); ok {
		return ok, nil
	}
	if ok := comparer.IsSubConfChanged("circuit"); ok {
		return ok, nil
	}
	return false, nil
}

//SetConf 设置配置参数
func (w *ApiResponsiveServer) SetConf(restart bool, cnf conf.IServerConf) (err error) {

	var ok bool
	//设置路由
	if restart {
		if _, err := SetHttpRouters(w.engine, w.server, cnf); err != nil {
			return err
		}
	}
	//设置静态文件
	if ok, err = SetStatic(w.server, cnf); err != nil {
		return err
	}
	servers.TraceIf(ok, w.Infof, w.Debugf, getEnableName(ok), "静态文件")

	//设置请求头
	if ok, err = SetHeaders(w.server, cnf); err != nil {
		return err
	}
	servers.TraceIf(ok, w.Infof, w.Debugf, getEnableName(ok), "header设置")

	//设置熔断配置
	if ok, err = SetCircuitBreaker(w.server, cnf); err != nil {
		return err
	}
	servers.TraceIf(ok, w.Infof, w.Debugf, getEnableName(ok), "熔断设置")

	//设置jwt安全认证
	if ok, err = SetJWT(w.server, cnf); err != nil {
		return err
	}
	servers.TraceIf(ok, w.Infof, w.Debugf, getEnableName(ok), "jwt设置")

	//设置ajax请求
	if ok, err = SetAjaxRequest(w.server, cnf); err != nil {
		return err
	}
	servers.TraceIf(ok, w.Infof, w.Debugf, getEnableName(ok), "ajax请求限制设置")

	//设置metric
	if ok, err = SetMetric(w.server, cnf); err != nil {
		return err
	}
	servers.TraceIf(ok, w.Infof, w.Debugf, getEnableName(ok), "metric设置")

	//设置host
	if ok, err = SetHosts(w.server, cnf); err != nil {
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
