package engines

import (
	"github.com/micro-plat/hydra/component"
)

//ServiceLoader 服务加载器
type ServiceLoader func(r *component.StandardComponent, i component.IContainer) error

var serviceLoaders = make([]ServiceLoader, 0, 2)

//AddServiceLoader 添加引擎加载器
func AddServiceLoader(f ServiceLoader) {
	serviceLoaders = append(serviceLoaders, f)
}

func (r *ServiceEngine) loadEngineServices() error {
	for _, loader := range serviceLoaders {
		if err := loader(r.StandardComponent, r); err != nil {
			return err
		}
	}
	return nil
}
