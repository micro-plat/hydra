package lnfs

import (
	"github.com/micro-plat/hydra/conf/app"
	"github.com/micro-plat/hydra/global"
	"github.com/micro-plat/hydra/hydra/servers/pkg/nfs/infs"
	"github.com/micro-plat/hydra/services"
	"github.com/micro-plat/lib4go/concurrent/cmap"
)

var servicesCache = cmap.New(2)

func init() {
	global.OnReady(func() {
		//处理服务初始化
		services.Def.OnSetup(func(c app.IAPPConf) error {
			unregistry(c.GetServerConf().GetServerType())
			return nil
		}, global.API, global.Web)

	})

}

func unregistry(tp string) error {
	srv, ok := servicesCache.Get(tp)
	if !ok {
		return nil
	}
	service := srv.([]string)
	for _, s := range service {
		services.Def.Remove(s)
		servicesCache.Remove(s)
	}
	return nil

}

func (c *LNFS) Registry(tp string) {
	if tp == global.API {
		//内部服务
		services.Def.API(infs.RMT_FP_GET, c.GetFP)
		services.Def.API(infs.RMT_FP_NOTIFY, c.RecvNotify)
		services.Def.API(infs.RMT_FP_QUERY, c.Query)
		services.Def.API(infs.RMT_FILE_DOWNLOAD, c.GetFile)
		service := make([]string, 0, 4)
		service = append(service, infs.RMT_FP_GET, infs.RMT_FP_NOTIFY, infs.RMT_FP_QUERY, infs.RMT_FILE_DOWNLOAD)
		servicesCache.Set(tp, service)
	}
	if tp == global.Web {

		//内部服务
		services.Def.Web(infs.RMT_FP_GET, c.GetFP)
		services.Def.Web(infs.RMT_FP_NOTIFY, c.RecvNotify)
		services.Def.Web(infs.RMT_FP_QUERY, c.Query)
		services.Def.Web(infs.RMT_FILE_DOWNLOAD, c.GetFile)
		service := make([]string, 0, 4)
		service = append(service, infs.RMT_FP_GET, infs.RMT_FP_NOTIFY, infs.RMT_FP_QUERY, infs.RMT_FILE_DOWNLOAD)
		servicesCache.Set(tp, service)
	}

}
