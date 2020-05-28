package rcontent

import (
	"fmt"

	"github.com/micro-plat/hydra/conf"
	"github.com/micro-plat/hydra/conf/server/render"
)

//GetConf 设置GetRender配置
func GetConf(cnf conf.IMainConf) (rsp *render.Render) {
	rsp = &render.Render{}
	_, err := cnf.GetSubObject("render/content", rsp)
	if err != nil && err != conf.ErrNoSetting {
		panic(fmt.Errorf("render/content配置有误:%v", err))
	}
	if err == conf.ErrNoSetting {
		rsp.Disable = true
		return rsp
	}
	return rsp
}
