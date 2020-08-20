package apm

import "github.com/micro-plat/hydra/components/pkgs/apm"

//IAPM 缓存接口
type IAPM = apm.IAPM

//IComponentAPM Component APM
type IComponentAPM interface {
	GetRegularAPM(instance string, names ...string) (c IAPM)
	GetAPM(instance string, names ...string) (c IAPM, err error)
}
