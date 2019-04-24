package hydra

type executable struct {
	app *MicroApp
}

func (p *executable) run() {
	p.app.service.Run(p)
}
func (p *executable) Run() {
	s, err := p.app.hydra.Start()
	if err != nil {
		p.app.logger.Error(err)
		return
	}
	p.app.logger.Info(s)
}
func (p *executable) Start() {
	go func() {
		p.Run()
	}()
	return
}

func (p *executable) Stop() {
	msg, err := p.app.service.Stop()
	if err != nil {
		p.app.logger.Error(err)
		return
	}
	p.app.logger.Info(msg)
}
