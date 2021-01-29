package pkgs

import (
	"regexp"

	"github.com/asaskevich/govalidator"
)

func init() {
	//注册服务名验证方式
	govalidator.TagMap["spath"] = govalidator.Validator(func(str string) bool {
		return regexp.MustCompile(`^(/[^/\x00]*)+/?$`).Match([]byte(str))
	})
}
