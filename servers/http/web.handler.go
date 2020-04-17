package http

import (
	x "net/http"
	"path"
	"path/filepath"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/micro-plat/hydra/conf"
	"github.com/micro-plat/hydra/servers"
	"github.com/micro-plat/hydra/servers/http/middleware"
)

func (s *WebServer) getHandler(routers []*conf.Router) (h x.Handler, err error) {
	if !servers.IsDebug {
		gin.SetMode(gin.ReleaseMode)
	}
	s.gin = gin.New()
	s.gin.Use(gin.Recovery())
	s.gin.Use(middleware.Logging(s.conf)) //记录请求日志
	s.gin.Use(middleware.Recovery())

	s.gin.Use(s.option.metric.Handle())       //生成metric报表
	s.gin.Use(middleware.Host(s.conf))        // 检查主机头是否合法
	s.gin.Use(middleware.Static(s.conf))      //处理静态文件
	s.gin.Use(middleware.JwtAuth(s.conf))     //jwt安全认证
	s.gin.Use(middleware.Body())              //处理请求form
	s.gin.Use(middleware.WebResponse(s.conf)) //处理返回值
	s.gin.Use(middleware.Header(s.conf))      //设置请求头
	s.gin.Use(middleware.JwtWriter(s.conf))   //jwt回写
	if err = setRouters(s.gin, routers); err != nil {
		return nil, err
	}
	return s.gin, nil
}
func (s *WebServer) loadHTMLGlob() (viewFiles []string, err error) {
	defer func() {
		if err1 := recover(); err1 != nil {
			err = err1.(error)
		}
	}()
	viewFiles = make([]string, 0, 8)
	viewRoot := "../views"
	if view, ok := s.conf.GetMetadata("view").(*conf.View); ok && !view.Disable {
		viewRoot = view.Path
	} else {
		s.conf.SetMetadata("view", &conf.View{Path: viewRoot})
	}

	dirs := []string{
		path.Join(viewRoot, "/*.html"),
	}
	for _, name := range dirs {
		filenames, err := filepath.Glob(name)
		if err != nil {
			return nil, err
		}
		viewFiles = append(viewFiles, filenames...)
	}
	if len(viewFiles) > 0 {
		s.gin.LoadHTMLFiles(viewFiles...)
	}
	s.conf.SetMetadata("viewFiles", viewFiles)
	servers.TraceIf(len(viewFiles) > 0, s.Logger.Infof, s.Logger.Debugf,
		getEnableName(len(viewFiles) > 0), "view模板", strings.Join(viewFiles, "\n"))
	return nil, nil
}
