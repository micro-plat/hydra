package skywalking

import (
	"encoding/json"
	"fmt"
	"time"

	"github.com/SkyAPM/go2sky"
	grpcreporter "github.com/SkyAPM/go2sky/reporter"

	"github.com/micro-plat/hydra/components/pkgs/apm"
	"github.com/micro-plat/lib4go/types"
	"google.golang.org/grpc/credentials"
)

func NewReporter(serverAddr string, config string) (reporter apm.Reporter, err error) {
	//fmt.Println("NewReporter:", serverAddr, config)
	opts, err := buildOptions(config)
	if err != nil {
		//fmt.Println("NewReporter.buildOptions.err:", err)
		return
	}

	skyreporter, err := grpcreporter.NewGRPCReporter(serverAddr, opts...)
	if err != nil {
		err = fmt.Errorf("构建skywalking Reporter 失败：%+v", err)
		//fmt.Println("NewReporter.NewGRPCReporter.err:", serverAddr, config)
		return
	}
	reporter = &innerreporter{
		reporter: skyreporter,
	}
	//fmt.Println("NewReporter.success")
	return
}

func buildOptions(config string) (opts []grpcreporter.GRPCReporterOption, err error) {
	/*	```
		{
			"check_interval":1,
			"max_send_queue_size:500000,
			"instance_props":{"":""},
			"transport_credentials":{"":""},
			"authentication_key":""
		}
		```
	*/
	cfgMap := types.XMap{}
	opts = make([]grpcreporter.GRPCReporterOption, 0)
	json.Unmarshal([]byte(config), &cfgMap)
	if v, ok := cfgMap["check_interval"]; ok {
		opts = append(opts, grpcreporter.WithCheckInterval(time.Second*time.Duration(types.GetInt64(v))))
	}

	if v, ok := cfgMap["max_send_queue_size"]; ok {
		opts = append(opts, grpcreporter.WithMaxSendQueueSize(types.GetInt(v)))
	}

	if v, ok := cfgMap["instance_props"]; ok {
		props := map[string]string{}
		tmpProps := v.(map[string]interface{})
		for k, v := range tmpProps {
			props[k] = types.GetString(v)
		}

		opts = append(opts, grpcreporter.WithInstanceProps(props))
	}
	if v, ok := cfgMap["transport_credentials"]; ok {
		var tmp types.XMap = v.(map[string]interface{})

		cred, err1 := credentials.NewClientTLSFromFile(tmp.GetString("cert_file"), tmp.GetString("server_name"))
		if err1 != nil {
			err = fmt.Errorf("从文件获取证书失败：%+v;cert_file:%s,server_name:%s", err1, tmp.GetString("cert_file"), tmp.GetString("server_name"))
			return
		}
		opts = append(opts, grpcreporter.WithTransportCredentials(cred))
	}

	if v, ok := cfgMap["authentication_key"]; ok {
		opts = append(opts, grpcreporter.WithAuthentication(types.GetString(v)))
	}

	return
}

type innerreporter struct {
	reporter go2sky.Reporter
}

func (r innerreporter) Boot(service string, serviceInstance string) {
	r.reporter.Boot(service, serviceInstance)
}

func (r innerreporter) Send(spans []apm.Span) {
	var reportspans = make([]go2sky.ReportedSpan, len(spans))
	for i, s := range spans {
		reportspans[i] = s.GetRealSpan().(go2sky.ReportedSpan)
	}
	r.reporter.Send(reportspans)
}
func (r innerreporter) Close() {
	r.reporter.Close()
}

func (r innerreporter) GetRealReporter() interface{} {
	return r.reporter
}
