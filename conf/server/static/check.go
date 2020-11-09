package static

import (
	"path/filepath"
	"strings"
)

//IsStatic 判断当前文件是否一定是静态文件,及文件的完整路径
func (s *Static) IsStatic(rPath string, method string) (b bool, xname string) {

	if b = s.IsFavRobot(rPath); b {
		return b, filepath.Join(s.Dir, rPath)
	}

	if !s.AllowRequest(method) {
		return false, ""
	}

	if s.IsExclude(rPath) {
		return false, ""
	}
	if s.IsContainExt(rPath) {
		return true, filepath.Join(s.Dir, rPath)
	}
	if s.HasPrefix(rPath) {
		return true, filepath.Join(s.Dir, strings.TrimPrefix(rPath, s.Prefix))
	}
	if s.NeedRewrite(rPath) {
		return true, filepath.Join(s.Dir, s.FirstPage)
	}
	return false, ""
}

//IsFavRobot 是否是favicon.ico 或 robots.txt
func (s *Static) IsFavRobot(rPath string) (b bool) {
	if rPath == "/favicon.ico" || rPath == "/robots.txt" {
		return true
	}
	return false
}

//HasPrefix 是否有指定的前缀
func (s *Static) HasPrefix(rPath string) bool {
	if s.Prefix == "" {
		return false
	}
	return strings.HasPrefix(rPath, s.Prefix)
}

//IsExclude 是否是排除的文件
func (s *Static) IsExclude(rPath string) bool {
	if len(s.Exclude) == 0 {
		return false
	}

	name := filepath.Base(rPath)
	pExt := filepath.Ext(name)
	hasExt := strings.Contains(pExt, ".")
	for _, v := range s.Exclude {
		if hasExt {
			return strings.EqualFold(pExt, v)
		}
		return strings.EqualFold(rPath, v)
	}
	return false
}

//IsContainExt 是否是包含在指定的ext中
func (s *Static) IsContainExt(rPath string) bool {

	name := filepath.Base(rPath)
	pExt := filepath.Ext(name)
	hasExt := strings.Contains(pExt, ".")
	if !hasExt {
		return false
	}
	if len(s.Exts) == 0 {
		return true
	}
	for _, ext := range s.Exts {
		if pExt == ext || ext == "*" {
			return true
		}
	}
	return false

}

//NeedRewrite 是否需要重写请求
func (s *Static) NeedRewrite(p string) bool {
	for _, c := range s.Rewriters {
		if c == p {
			return true
		}
	}
	return false
}
