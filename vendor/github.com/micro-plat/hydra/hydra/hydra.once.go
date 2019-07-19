package hydra

type once struct {
	app *MicroApp
}

func (p *once) run() {
	p.app.service.Run(p)
}
func (p *once) Run() {
	s, err := p.app.hydra.Start()
	if err != nil {
		p.app.logger.Error(err)
		return
	}
	p.app.logger.Info(s)
}
func (p *once) Start() {
	go func() {
		p.Run()
	}()
	return
}

func (p *once) Stop() {
	msg, err := p.app.service.Stop()
	if err != nil {
		p.app.logger.Error(err)
		return
	}
	p.app.logger.Info(msg)
}
