package main

import (
	"embed"

	"github.com/micro-plat/hydra"
	"github.com/micro-plat/hydra/conf/server/static"
	"github.com/micro-plat/hydra/hydra/servers/http"
)

//go:embed web/*
var mgrweb embed.FS

var app = hydra.NewApp(
	hydra.WithPlatName("hydra", "hydra"),
	hydra.WithSystemName("docs", "文档中心"),
	hydra.WithServerTypes(http.Web),
)

func main() {
	//设置配置参数
	hydra.Conf.Web("8080").
		Static(static.WithAutoRewrite(), static.WithEmbed("web", mgrweb), static.WithHomePage("index.html"))
	app.Start()
}
