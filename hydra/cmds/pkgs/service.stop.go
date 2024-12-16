/*
 * @Description:
 * @Autor: taoshouyin
 * @Date: 2020-12-23 15:43:54
 * @LastEditors: taoshouyin
 * @LastEditTime: 2021-09-26 17:57:50
 */
package pkgs

import (
	"github.com/micro-plat/hydra/components"
	"github.com/micro-plat/hydra/global"
	"github.com/micro-plat/hydra/hydra/cmds/pkgs/service"
	"github.com/micro-plat/hydra/registry"
	"github.com/micro-plat/hydra/services"
)

//Stop Stop
func (p *ServiceApp) Stop(s service.Service) (err error) {
	globalLogger := global.Def.Log()
	//8. 关闭服务器释放所有资源
	globalLogger.Info(global.AppName, "正在退出...")

	//关闭所有组件
	if err := components.Def.Container().Close(); err != nil {
		globalLogger.Error("关闭容器中组件失败:", err)
	}

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

	//关闭注册中心
	registry.Close()

	globalLogger.Info(global.AppName, "已安全退出")

	return nil
}
