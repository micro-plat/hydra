module github.com/micro-plat/hydra

go 1.15

require (
	github.com/BurntSushi/toml v0.3.1
	github.com/asaskevich/govalidator v0.0.0-20190424111038-f61b66f89f4a
	github.com/bradfitz/gomemcache v0.0.0-20190913173617-a41fca850d0b
	github.com/clbanning/mxj v1.8.4
	github.com/d5/tengo/v2 v2.6.2
	github.com/dsnet/compress v0.0.1 // indirect
	github.com/eclipse/paho.mqtt.golang v1.2.0
	github.com/frankban/quicktest v1.10.2 // indirect
	github.com/fsnotify/fsnotify v1.4.9
	github.com/gin-gonic/gin v1.6.2
	github.com/go-redis/redis v6.15.6+incompatible
	github.com/golang/protobuf v1.4.2
	github.com/golang/snappy v0.0.1
	github.com/gorilla/websocket v1.4.2
	github.com/json-iterator/go v1.1.10 // indirect
	github.com/kr/pretty v0.2.1
	github.com/mattn/go-oci8 v0.0.8 // indirect
	github.com/mholt/archiver v3.1.1+incompatible
	github.com/micro-plat/cli v1.1.0
	github.com/micro-plat/gmq v1.0.1
	github.com/micro-plat/lib4go v1.0.2
	github.com/modern-go/concurrent v0.0.0-20180306012644-bacd9c7ef1dd // indirect
	github.com/modern-go/reflect2 v1.0.1 // indirect
	github.com/nwaples/rardecode v1.1.0 // indirect
	github.com/pierrec/lz4 v2.5.2+incompatible // indirect
	github.com/pkg/profile v1.4.0
	github.com/pkg/sftp v1.12.0
	github.com/shopspring/decimal v1.2.0
	github.com/ugorji/go/codec v1.1.7
	github.com/ulikunitz/xz v0.5.7 // indirect
	github.com/urfave/cli v1.22.4
	github.com/xi2/xz v0.0.0-20171230120015-48954b6210f8 // indirect
	github.com/zkfy/cron v0.0.0-20170309132418-df38d32658d8
	github.com/zkfy/go-cache v2.1.0+incompatible
	github.com/zkfy/go-metrics v0.0.0-20161128210544-1f30fe9094a5
	github.com/zkfy/log v0.0.0-20180312054228-b2704c3ef896
	github.com/zkfy/stompngo v0.0.0-20170803022748-9378e70ca481
	golang.org/x/crypto v0.0.0-20200820211705-5c72a883971a
	golang.org/x/net v0.0.0-20200707034311-ab3426394381
	golang.org/x/sys v0.0.0-20200803210538-64077c9b5642
	golang.org/x/text v0.3.3
	golang.org/x/time v0.0.0-20200630173020-3af7569d3a1e
	golang.org/x/xerrors v0.0.0-20200804184101-5ec99f83aff1 // indirect
	google.golang.org/genproto v0.0.0-20200806141610-86f49bd18e98 // indirect
	google.golang.org/grpc v1.32.0
	google.golang.org/protobuf v1.25.0 // indirect
	gopkg.in/check.v1 v1.0.0-20190902080502-41f04d3bba15 // indirect
	gopkg.in/yaml.v2 v2.3.0
)

replace github.com/micro-plat/lib4go => ../../../github.com/micro-plat/lib4go
