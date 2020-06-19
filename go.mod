module github.com/micro-plat/hydra

go 1.13

require (
	github.com/BurntSushi/toml v0.3.1
	github.com/asaskevich/govalidator v0.0.0-20190424111038-f61b66f89f4a
	github.com/bradfitz/gomemcache v0.0.0-20190913173617-a41fca850d0b
	github.com/eclipse/paho.mqtt.golang v1.2.0
	github.com/gin-gonic/gin v1.6.2
	github.com/go-redis/redis v6.15.6+incompatible
	github.com/golang/protobuf v1.3.3
	github.com/golang/snappy v0.0.1
	github.com/micro-plat/cli v1.1.0
	github.com/micro-plat/gmq v1.0.1
	github.com/micro-plat/lib4go v0.4.0
	github.com/pkg/profile v1.4.0
	github.com/qxnw/lib4go v0.0.0-20180426074627-c80c7e84b925
	github.com/ugorji/go/codec v1.1.7
	github.com/urfave/cli v1.22.4
	github.com/zkfy/cron v0.0.0-20170309132418-df38d32658d8
	github.com/zkfy/go-cache v2.1.0+incompatible
	github.com/zkfy/go-metrics v0.0.0-20161128210544-1f30fe9094a5
	github.com/zkfy/log v0.0.0-20180312054228-b2704c3ef896
	github.com/zkfy/stompngo v0.0.0-20170803022748-9378e70ca481
	golang.org/x/net v0.0.0-20200506145744-7e3656a0809f
	golang.org/x/sys v0.0.0-20200508214444-3aab700007d7
	google.golang.org/grpc v1.27.1
	gopkg.in/yaml.v2 v2.2.8
)

replace github.com/micro-plat/lib4go => ../../../github.com/micro-plat/lib4go
