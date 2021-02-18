module github.com/micro-plat/hydra

go 1.16

require (
	github.com/BurntSushi/toml v0.3.1
	github.com/SkyAPM/go2sky v0.6.0
	github.com/asaskevich/govalidator v0.0.0-20200907205600-7a23bdc65eef
	github.com/bradfitz/gomemcache v0.0.0-20190913173617-a41fca850d0b
	github.com/clbanning/mxj v1.8.4
	github.com/d5/tengo/v2 v2.6.2
	github.com/dsnet/compress v0.0.1 // indirect
	github.com/eclipse/paho.mqtt.golang v1.2.0
	github.com/frankban/quicktest v1.11.3 // indirect
	github.com/fsnotify/fsnotify v1.4.9
	github.com/gin-gonic/gin v1.6.3
	github.com/gmallard/stompngo v1.0.13 // indirect
	github.com/go-redis/redis v6.15.9+incompatible
	github.com/golang/protobuf v1.4.3
	github.com/golang/snappy v0.0.2
	github.com/gorilla/websocket v1.4.2
	github.com/lib4dev/cli v1.2.8
	github.com/manifoldco/promptui v0.8.0
	github.com/mholt/archiver v3.1.1+incompatible
	github.com/micro-plat/gmq v1.0.1
	github.com/micro-plat/lib4go v1.0.10
	github.com/nwaples/rardecode v1.1.0 // indirect
	github.com/pierrec/lz4 v2.6.0+incompatible // indirect
	github.com/pkg/profile v1.2.1
	github.com/pkg/sftp v1.12.0
	github.com/robfig/cron/v3 v3.0.1
	github.com/stretchr/testify v1.6.1
	github.com/ugorji/go/codec v1.2.2
	github.com/ulikunitz/xz v0.5.9 // indirect
	github.com/urfave/cli v1.22.5
	github.com/xi2/xz v0.0.0-20171230120015-48954b6210f8 // indirect
	github.com/zkfy/go-cache v2.1.0+incompatible
	github.com/zkfy/go-metrics v0.0.0-20161128210544-1f30fe9094a5
	github.com/zkfy/log v0.0.0-20180312054228-b2704c3ef896
	github.com/zkfy/stompngo v0.0.0-20170803022748-9378e70ca481
	golang.org/x/crypto v0.0.0-20201221181555-eec23a3978ad
	golang.org/x/net v0.0.0-20201216054612-986b41b23924
	golang.org/x/sys v0.0.0-20201221093633-bc327ba9c2f0
	golang.org/x/time v0.0.0-20201208040808-7e3f01d25324
	google.golang.org/grpc v1.34.0
	gopkg.in/yaml.v2 v2.3.0
	gopkg.in/yaml.v3 v3.0.0-20200615113413-eeeca48fe776
)

replace github.com/micro-plat/lib4go => ../../../github.com/micro-plat/lib4go
