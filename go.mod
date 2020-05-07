module github.com/micro-plat/hydra

go 1.13

require (
	github.com/asaskevich/govalidator v0.0.0-20190424111038-f61b66f89f4a
	github.com/fsnotify/fsnotify v1.4.7
	github.com/gin-gonic/gin v1.6.2
	github.com/golang/protobuf v1.3.3
	github.com/micro-plat/cli v1.0.0
	github.com/micro-plat/lib4go v0.3.1
	github.com/pkg/profile v1.4.0
	github.com/ugorji/go/codec v1.1.7
	github.com/urfave/cli v1.22.4
	github.com/zkfy/log v0.0.0-20180312054228-b2704c3ef896
	golang.org/x/net v0.0.0-20190311183353-d8887717615a
	golang.org/x/sys v0.0.0-20200116001909-b77594299b42
	google.golang.org/grpc v1.27.1
	gopkg.in/yaml.v2 v2.2.8
)

replace github.com/micro-plat/lib4go => ../../../github.com/micro-plat/lib4go
