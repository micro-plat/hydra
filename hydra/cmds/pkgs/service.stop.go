package pkgs

import (
	"time"

	"github.com/micro-plat/hydra/global"
	"github.com/micro-plat/hydra/hydra/cmds/pkgs/service"
	"github.com/micro-plat/hydra/services"
	"github.com/micro-plat/lib4go/logger"
)

//Stop Stop
func (p *ServiceApp) Stop(s service.Service) (err error) {

	//8. 关闭服务器释放所有资源
	global.Def.Log().Info(global.AppName, "正在退出...")

	if p.server != nil {
		//if !reflect.ValueOf(p.server).IsNil() {
		//关闭服务器
		p.server.Shutdown()
	}

	//关闭各服务
	if err := services.Def.Close(); err != nil {
		global.Def.Log().Error("关闭服务失败:", err)
	}

	if p.trace != nil {
		p.trace.Stop()
	}

	//通知关闭各组件
	global.Def.Close()

	global.Def.Log().Info(global.AppName, "已安全退出")

	//关闭日志
	logger.Close()
	time.Sleep(time.Millisecond * 1000)
	return nil
}
