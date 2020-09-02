module github.com/micro-plat/hydra

go 1.14

require (
	github.com/BurntSushi/toml v0.3.1
	github.com/SkyAPM/go2sky v0.5.0
	github.com/asaskevich/govalidator v0.0.0-20190424111038-f61b66f89f4a
	github.com/bradfitz/gomemcache v0.0.0-20190913173617-a41fca850d0b
	github.com/dsnet/compress v0.0.1 // indirect
	github.com/eclipse/paho.mqtt.golang v1.2.0
	github.com/frankban/quicktest v1.10.2 // indirect
	github.com/gin-gonic/gin v1.6.2
	github.com/go-redis/redis v6.15.6+incompatible
	github.com/golang/protobuf v1.3.3
	github.com/golang/snappy v0.0.1
	github.com/gorilla/websocket v1.4.2
	github.com/keybase/go-ps v0.0.0-20190827175125-91aafc93ba19 // indirect
	github.com/mattn/go-oci8 v0.0.8
	github.com/mholt/archiver v3.1.1+incompatible
	github.com/micro-plat/cli v1.1.0
	github.com/micro-plat/gmq v1.0.1
	github.com/micro-plat/lib4go v0.4.0
	github.com/nwaples/rardecode v1.1.0 // indirect
	github.com/pierrec/lz4 v2.5.2+incompatible // indirect
	github.com/pkg/profile v1.4.0
	github.com/pkg/sftp v1.12.0
	github.com/qxnw/lib4go v0.0.0-20180426074627-c80c7e84b925
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
	golang.org/x/net v0.0.0-20200506145744-7e3656a0809f
	golang.org/x/sys v0.0.0-20200508214444-3aab700007d7
	golang.org/x/time v0.0.0-20200630173020-3af7569d3a1e
	google.golang.org/grpc v1.27.1
	gopkg.in/yaml.v2 v2.2.8
)

replace github.com/micro-plat/lib4go => ../../../github.com/micro-plat/lib4go
