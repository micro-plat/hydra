package nfs

import (
	"github.com/micro-plat/hydra/conf/app"
	"github.com/micro-plat/hydra/conf/server/nfs"
	"github.com/micro-plat/hydra/global"
	"github.com/micro-plat/hydra/hydra/servers/pkg/nfs/infs"
	"github.com/micro-plat/hydra/services"
	"github.com/micro-plat/lib4go/concurrent/cmap"
	"github.com/micro-plat/lib4go/types"
)

//currentModule 当前Module
var currentModule infs.Infs

type cnfs struct {
	c        *nfs.NFS
	app      app.IAPPConf
	services []string
	infs     infs.Infs
}

func newNFS(app app.IAPPConf, c *nfs.NFS) *cnfs {
	currentModule = getNFS(app, c)
	return &cnfs{c: c,
		app:      app,
		infs:     currentModule,
		services: make([]string, 0, 3)}
}
func (c *cnfs) Start() error {
	return c.infs.Start()
}

func (r *cnfs) Close() error {
	return r.infs.Close()
}

func init() {
	global.OnReady(func() {
		//处理服务初始化
		services.Def.OnSetup(func(c app.IAPPConf) error {
			//取消服务注册
			closeNFS(c.GetServerConf().GetServerType())
			n, err := c.GetNFSConf()
			if err != nil {
				return err
			}
			if n.Disable {
				return nil
			}

			//构建并缓存nfs
			cnfs := newNFS(c, n)
			nfsCaches.Set(c.GetServerConf().GetServerType(), cnfs)

			//注册服务
			registry(c.GetServerConf().GetServerType(), cnfs, n)
			cnfs.infs.Registry(c.GetServerConf().GetServerType())
			return nil
		})

		//处理服务启动完成
		services.Def.OnStarted(func(c app.IAPPConf) error {
			return startNFS(c.GetServerConf().GetServerType())
		})

	})

}

var nfsCaches cmap.ConcurrentMap = cmap.New(2)

func startNFS(tp string) error {
	v, ok := nfsCaches.Get(tp)
	if !ok {
		return nil
	}
	m := v.(*cnfs)
	return m.Start()
}

func closeNFS(tp string) error {
	nfsCaches.RemoveIterCb(func(k string, v interface{}) bool {
		if k == tp {
			m := v.(*cnfs)
			for _, s := range m.services {
				services.Def.Remove(s)
			}

			m.Close()
			return true
		}
		return false
	})
	return nil

}

func registry(tp string, cnfs *cnfs, cnf *nfs.NFS) {
	if tp == global.API {
		//注册服务
		if !cnf.DiableUpload {
			s := types.GetString(cnf.UploadService, infs.SVSUpload)
			services.Def.API(s, cnfs.Upload)
			cnfs.services = append(cnfs.services, s)
		}

		if cnf.AllowDownload {
			s := types.GetString(cnf.DownloadService, infs.SVSDonwload)
			services.Def.API(s, cnfs.Download)
			cnfs.services = append(cnfs.services, s)
		}
		if cnf.AllowListFile {
			s := types.GetString(cnf.ListFileService, infs.SVSList)
			services.Def.API(s, cnfs.GetFileList)
			cnfs.services = append(cnfs.services, s)
		}

		if cnf.AllowListDir {
			s := types.GetString(cnf.ListDirService, infs.SVSDir)
			services.Def.API(s, cnfs.GetDirList)
			cnfs.services = append(cnfs.services, s)
		}

		if cnf.AllowPreview {
			s := types.GetString(cnf.PreviewService, infs.SVSPreview)
			services.Def.API(s, cnfs.GetPDF4Preview)
			cnfs.services = append(cnfs.services, s)
		}

		if cnf.AllowScaleImage {
			s := types.GetString(cnf.ScaleImageService, infs.SVSScalrImage)
			services.Def.API(s, cnfs.ImgScale)
			cnfs.services = append(cnfs.services, s)
		}
		if cnf.AllowCreateDir {
			s := types.GetString(cnf.CreateDirService, infs.SVSCreateDir)
			services.Def.API(s, cnfs.CreateDir)
			cnfs.services = append(cnfs.services, s)
		}

		if cnf.AllowRenameDir {
			s := types.GetString(cnf.RenameDirService, infs.SVSRenameDir)
			services.Def.API(s, cnfs.RenameDir)
			cnfs.services = append(cnfs.services, s)
		}
	}

	if tp == global.Web {

		if !cnf.DiableUpload {
			s := types.GetString(cnf.UploadService, infs.SVSUpload)
			services.Def.Web(s, cnfs.Upload)
			cnfs.services = append(cnfs.services, s)
		}

		if cnf.AllowDownload {
			s := types.GetString(cnf.DownloadService, infs.SVSDonwload)
			services.Def.Web(s, cnfs.Download)
			cnfs.services = append(cnfs.services, s)
		}

		if cnf.AllowListFile {
			s := types.GetString(cnf.ListFileService, infs.SVSList)
			services.Def.Web(s, cnfs.GetFileList)
			cnfs.services = append(cnfs.services, s)
		}

		if cnf.AllowListDir {
			s := types.GetString(cnf.ListDirService, infs.SVSDir)
			services.Def.Web(s, cnfs.GetDirList)
			cnfs.services = append(cnfs.services, s)
		}

		if cnf.AllowPreview {
			s := types.GetString(cnf.PreviewService, infs.SVSPreview)
			services.Def.Web(s, cnfs.GetPDF4Preview)
			cnfs.services = append(cnfs.services, s)
		}

		if cnf.AllowScaleImage {
			s := types.GetString(cnf.ScaleImageService, infs.SVSScalrImage)
			services.Def.Web(s, cnfs.ImgScale)
			cnfs.services = append(cnfs.services, s)
		}

		if cnf.AllowCreateDir {
			s := types.GetString(cnf.CreateDirService, infs.SVSCreateDir)
			services.Def.Web(s, cnfs.CreateDir)
			cnfs.services = append(cnfs.services, s)
		}

		if cnf.AllowRenameDir {
			s := types.GetString(cnf.RenameDirService, infs.SVSRenameDir)
			services.Def.Web(s, cnfs.RenameDir)
			cnfs.services = append(cnfs.services, s)
		}
	}
}
