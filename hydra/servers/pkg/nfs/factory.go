package nfs

import (
	"fmt"
	"strings"

	"github.com/micro-plat/hydra/conf/app"
	"github.com/micro-plat/hydra/conf/server/nfs"
	"github.com/micro-plat/hydra/hydra/servers/pkg/nfs/infs"
	"github.com/micro-plat/hydra/hydra/servers/pkg/nfs/lnfs"
	"github.com/micro-plat/hydra/hydra/servers/pkg/nfs/obs"
)

func getNFS(app app.IAPPConf, c *nfs.NFS) infs.Infs {
	fmt.Println("obs:", c.CloudNFS)
	switch strings.ToUpper(c.CloudNFS) {
	case "OBS":
		return obs.NewOBS(c.AccessKey, c.SecretKey, c.Local, c.Endpoint)
	case "":
		return lnfs.NewNFS(app, c)
	default:
		panic(fmt.Sprintf("不支持的NFS类型:%s", c.CloudNFS))
	}

}
