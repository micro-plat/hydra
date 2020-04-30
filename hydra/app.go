package hydra

import (
	"github.com/micro-plat/cli"
	"github.com/micro-plat/hydra/application"
	"github.com/micro-plat/lib4go/logger"
)

//MicroApp  微服务应用
type MicroApp struct {
	app *cli.App
}

//NewApp 创建微服务应用
func NewApp(opts ...Option) (m *MicroApp) {
	m = &MicroApp{}
	for _, opt := range opts {
		opt()
	}
	return m
}

//Start 启动服务器
func (m *MicroApp) Start() {
	defer logger.Close()
	m.app = cli.New(cli.WithVersion(application.Version))
	m.app.Start()

}
