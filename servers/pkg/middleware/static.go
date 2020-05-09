package middleware

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/micro-plat/hydra/registry/conf/server/static"
)

//Static 静态文件处理插件
func Static() Handler {
	return func(ctx IMiddleContext) {
		opt := ctx.ServerConf().GetStaticConf()
		if opt.Disable || opt.AllowRequest(ctx.Request().Path().GetMethod()) {
			ctx.Next()
			return
		}
		var rPath = ctx.Request().Path().GetPath()
		s, xname := MustStatic(opt, getDefPath(opt, rPath))
		if !s {
			ctx.Next()
			return
		}

		ctx.Response().AddSpecial("static")
		fPath := filepath.Join(opt.Dir, xname)
		finfo, err := os.Stat(fPath)
		if err != nil {
			if os.IsNotExist(err) {
				ctx.Log().Error(fmt.Errorf("找不到文件:%s", fPath))
				ctx.Response().Abort(404)
				return
			}
			ctx.Log().Error(fmt.Errorf("%s,err:%v", fPath, err))
			ctx.Response().Abort(500)
			return
		}
		if finfo.IsDir() {
			ctx.Log().Error(fmt.Errorf("找不到文件:%s", fPath))
			ctx.Response().Abort(404)
			return
		}
		//文件已存在，则返回文件
		ctx.Response().File(fPath)
		return
	}
}

func checkPrefix(s *static.Static, rPath string) bool {
	if strings.HasPrefix(rPath, s.Prefix) {
		return true
	}
	return false
}

func checkExclude(all []string, rPath string) bool {
	if len(all) == 0 {
		return true
	}
	for _, v := range all {
		if strings.Contains(rPath, v) {
			return false
		}
	}
	return true
}

func checkExt(s *static.Static, rPath string) bool {
	name := filepath.Base(rPath)
	hasExt := strings.Contains(filepath.Ext(name), ".")
	if len(s.Exts) > 0 {
		if s.Exts[0] == "*" && (rPath == "/" || hasExt) {
			return true
		}
		pExt := filepath.Ext(name)
		for _, ext := range s.Exts {
			if pExt == ext {
				return true
			}
		}
		return false
	}
	return true
}

//MustStatic 判断当前文件是否一定是静态文件 0:非静态文件  1：是静态文件  2：未知
func MustStatic(s *static.Static, rPath string) (b bool, xname string) {
	if rPath == "/favicon.ico" || rPath == "robots.txt" {
		return true, rPath
	}
	if checkExclude(s.Exclude, rPath) && checkPrefix(s, rPath) && checkExt(s, rPath) {
		return true, strings.TrimPrefix(rPath, s.Prefix)
	}
	return false, ""
}
func getDefPath(s *static.Static, p string) string {
	pExt := filepath.Ext(p)
	for _, c := range s.Rewriters {
		if c == p || (c == "*" && pExt == "") {
			if s.FirstPage != "" {
				return filepath.Join("/", s.FirstPage)
			}
			return "/index.html"
		}
	}
	if p == "" || p == "/" {
		if s.FirstPage != "" {
			return filepath.Join("/", s.FirstPage)
		}
		return "/index.html"
	}
	return p
}
