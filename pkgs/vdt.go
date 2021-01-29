package pkgs

import (
	"regexp"

	"github.com/asaskevich/govalidator"
)

func init() {
	//注册服务名验证方式
	govalidator.ParamTagRegexMap["spath"] = regexp.MustCompile(`^(/[^/\x00]*)+/?$`)
}
