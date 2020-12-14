package pkgs

import (
	"github.com/micro-plat/hydra/global"
	"github.com/micro-plat/hydra/hydra/cmds/pkgs/service"
	"github.com/micro-plat/hydra/services"
)

//Stop Stop
func (p *ServiceApp) Stop(s service.Service) (err error) {
	globalLogger := global.Def.Log()
	//8. 关闭服务器释放所有资源
	globalLogger.Info(global.AppName, "正在退出...")

	if p.server != nil {
		//if !reflect.ValueOf(p.server).IsNil() {
		//关闭服务器
		p.server.Shutdown()
	}

	if p.trace != nil {
		p.trace.Stop()
	}
	//关闭各服务
	if err := services.Def.Close(); err != nil {
		globalLogger.Error("关闭服务失败:", err)
	}

	//通知关闭各组件
	global.Def.Close()

	globalLogger.Info(global.AppName, "已安全退出")

	return nil
}
