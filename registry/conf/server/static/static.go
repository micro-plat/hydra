package static

import (
	"fmt"

	"github.com/asaskevich/govalidator"
	"github.com/micro-plat/hydra/conf"
)

//Static 设置静态文件配置
type Static struct {
	*option
}

//New 构建静态文件配置信息
func New(opts ...Option) *Static {
	s := &Static{option: newOption()}
	for _, opt := range opts {
		opt(s.option)
	}
	return s
}

//GetStatic 设置static
func GetStatic(cnf conf.IMainConf) (static *Static, err error) {
	//设置静态文件路由
	_, err = cnf.GetSubObject("static", &static)
	if err != nil && err != conf.ErrNoSetting {
		return nil, err
	}
	if err == conf.ErrNoSetting {
		static = New()
	}
	if b, err := govalidator.ValidateStruct(&static); !b {
		return nil, fmt.Errorf("static配置有误:%v", err)
	}
	return
}
