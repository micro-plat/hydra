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
	engine := gin.New()
	if s.views, err = s.loadHTMLGlob(engine); err != nil {
		s.Logger.Debugf("%s未找到模板:%v", s.conf.Name, err)
		return nil, err
	}
	engine.Use(gin.Recovery())
	engine.Use(middleware.Logging(s.conf)) //记录请求日志
	engine.Use(middleware.Recovery())

	engine.Use(s.option.metric.Handle())       //生成metric报表
	engine.Use(middleware.Host(s.conf))        // 检查主机头是否合法
	engine.Use(middleware.Static(s.conf))      //处理静态文件
	engine.Use(middleware.JwtAuth(s.conf))     //jwt安全认证
	engine.Use(middleware.Body())              //处理请求form
	engine.Use(middleware.WebResponse(s.conf)) //处理返回值
	engine.Use(middleware.Header(s.conf))      //设置请求头
	engine.Use(middleware.JwtWriter(s.conf))   //jwt回写
	if err = setRouters(engine, routers); err != nil {
		return nil, err
	}
	return engine, nil
}
func (s *WebServer) loadHTMLGlob(engine *gin.Engine) (viewFiles []string, err error) {
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
		path.Join(viewRoot, "/**/*.html"),
	}
	for _, name := range dirs {
		filenames, err := filepath.Glob(name)
		if err != nil {
			return nil, err
		}
		viewFiles = append(viewFiles, filenames...)
	}
	if len(viewFiles) > 0 {
		engine.LoadHTMLFiles(viewFiles...)
	}
	s.conf.SetMetadata("viewFiles", viewFiles)
	servers.TraceIf(len(viewFiles) > 0, s.Logger.Infof, s.Logger.Debugf,
		getEnableName(len(viewFiles) > 0), "view模板", strings.Join(viewFiles, "\n"))
	return nil, nil
}
