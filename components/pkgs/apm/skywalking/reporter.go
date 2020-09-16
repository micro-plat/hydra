package skywalking

import (
	"fmt"
	"time"

	"github.com/SkyAPM/go2sky"

	grpcreporter "github.com/SkyAPM/go2sky/reporter"
	conf "github.com/micro-plat/hydra/conf/vars/apm"
	"github.com/micro-plat/hydra/context/apm"
	"google.golang.org/grpc/credentials"
)

func NewReporter(serverAddr string, config *conf.APM) (reporter apm.Reporter, err error) {
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

func buildOptions(config *conf.APM) (opts []grpcreporter.GRPCReporterOption, err error) {
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
	opts = make([]grpcreporter.GRPCReporterOption, 0)
	if config.ReportCheckInterval > 0 {
		opts = append(opts, grpcreporter.WithCheckInterval(time.Second*time.Duration(config.ReportCheckInterval)))
	}

	if config.MaxSendQueueSize > 0 {
		opts = append(opts, grpcreporter.WithMaxSendQueueSize(config.MaxSendQueueSize))
	}

	if len(config.InstanceProps) > 0 {
		opts = append(opts, grpcreporter.WithInstanceProps(config.InstanceProps))
	}
	if config.Credentials != nil {
		tmp := config.Credentials
		if tmp.CertFile == "" || tmp.ServerName == "" {
			err = fmt.Errorf("构建skywalking Reporter 失败：Credentials设置有值，但CertFile与ServerName存在为空数据")
			return
		}
		cred, err1 := credentials.NewClientTLSFromFile(tmp.CertFile, tmp.CertFile)
		if err1 != nil {
			err = fmt.Errorf("从文件获取证书失败：%+v;cert_file:%s,server_name:%s", err1, tmp.CertFile, tmp.ServerName)
			return
		}

		opts = append(opts, grpcreporter.WithTransportCredentials(cred))
	}

	if len(config.AuthenticationKey) > 0 {
		opts = append(opts, grpcreporter.WithAuthentication(config.AuthenticationKey))
	}

	return
}

type innerreporter struct {
	reporter go2sky.Reporter
}

func (r *innerreporter) Boot(service string, serviceInstance string) {
	r.reporter.Boot(service, serviceInstance)
}

func (r *innerreporter) Send(spans []apm.Span) {
	var reportspans = make([]go2sky.ReportedSpan, len(spans))
	for i, s := range spans {
		reportspans[i] = s.GetRealSpan().(go2sky.ReportedSpan)
	}
	r.reporter.Send(reportspans)
}
func (r *innerreporter) Close() {
	r.reporter.Close()
}

func (r innerreporter) GetRealReporter() interface{} {
	return r.reporter
}
