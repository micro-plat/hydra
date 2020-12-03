package pkgs

import "github.com/micro-plat/hydra/hydra/cmds/pkgs/service"

//Start Start
func (p *ServiceApp) Start(s service.Service) (err error) {
	err = p.run()
	return err
}
