# 构建静态文件服务

静态文件服务器实际就是 api 或 web 服务器。api/web 服务器启动时已默认启动了静态文件服务。

#### 1. 参数配置

```go
package main

import (
	"github.com/micro-plat/hydra/hydra"
)

type webserver struct {
	*hydra.MicroApp
}

func main() {
	app := &webserver{
		hydra.NewApp(
			hydra.WithPlatName("mall"),
			hydra.WithSystemName("webserver"),
			hydra.WithServerTypes("api"),
		),
	}
    app.Conf.API.SetSubConf("static", `{
            "dir":"./static",
            "first-page":"index.html",
			"rewriters":["*"],
            "exts":[".jpg", ".png", ".gif", ".ico", ".html", ".htm", ".js", ".css", ".map", ".ttf", ".woff", ".woff2"],
            "exclude":["/views/", ".exe", ".so"]
			}`)
	app.Start()
}
```

参数说明:

|   参数名   |                                              默认值                                               | 说明                                                    |
| :--------: | :-----------------------------------------------------------------------------------------------: | ------------------------------------------------------- |
|    dir     |                                             ./static/                                             | 静态文件存放路径                                        |
| first-page |                                            index.html                                             | 首页地址                                                |
| rewriters  |                                 "/", "index.htm", "default.html"                                  | 哪些页面需重写到首页                                    |
|    exts    | ".jpg", ".png", ".gif", ".ico", ".html", ".htm", ".js", ".css", ".map", ".ttf", ".woff", ".woff2" | 支持的文件扩展名                                        |
|  exclude   |                                     "/views/", ".exe", ".so"                                      | 路径中需排除的名称。包括以上名称则请求会返回 404 或 403 |

#### 2. vuejs 项目中使用 hydra 作为 web 服务器

```go
package main

import (
	"github.com/micro-plat/hydra/hydra"
)

type webserver struct {
	*hydra.MicroApp
}

func main() {
	app := &webserver{
		hydra.NewApp(
			hydra.WithPlatName("mall"),
			hydra.WithSystemName("webserver"),
			hydra.WithServerTypes("web"),
		),
	}
    app.Conf.WEB.SetSubConf("static", `{
            "dir":"./static",
			"rewriters":["*"],
			"exts":[".ttf",".woff",".woff2"]
			}`)
	app.Start()
}
```

vuejs 项目中可新建`build.sh`文件，填写以下内容进行项目编译:

```sh
npm run build
go install
cp ${GOPATH}/bin/webserver ./dist

```
