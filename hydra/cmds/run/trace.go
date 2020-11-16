package run

import (
	"fmt"
	"net/http"
	_ "net/http/pprof"
	"os"
	"runtime/trace"
	"strconv"
	"time"

	"github.com/micro-plat/lib4go/logger"
	"github.com/pkg/profile"
)

var supportTraces = []string{"cpu", "mem", "block", "mutex", "web"}

//startTrace 启用项目性能跟踪
func startTrace(trace, tracePort string) error {
	switch trace {
	case "cpu":
		defer profile.Start(profile.CPUProfile, profile.ProfilePath("."), profile.NoShutdownHook).Stop()
	case "mem":
		defer profile.Start(profile.MemProfile, profile.ProfilePath("."), profile.NoShutdownHook).Stop()
	case "block":
		defer profile.Start(profile.BlockProfile, profile.ProfilePath("."), profile.NoShutdownHook).Stop()
	case "mutex":
		defer profile.Start(profile.MutexProfile, profile.ProfilePath("."), profile.NoShutdownHook).Stop()
	case "web":
		return startWebTrace(tracePort)
	case "":
		return nil
	default:
		return fmt.Errorf("不支持trace命令:%v", trace)
	}
	return nil
}
func startWebTrace(tracePort string) error {
	errChan := make(chan error, 1)
	go startTraceServer(tracePort, errChan)
	select {
	case err := <-errChan:
		return err
	case <-time.After(time.Millisecond * 200):
		return nil
	}
}

func startTraceServer(tracePort string, errChan chan error) {
	f, err := os.Create("trace.out")
	if err != nil {
		errChan <- fmt.Errorf("启动pprof，创建监控输出文件错误：%w", err)
		return
	}
	defer f.Close()
	err = trace.Start(f)
	if err != nil {
		errChan <- fmt.Errorf("启动pprof，trace.Start错误：%w", err)
		return
	}
	defer trace.Stop()
	if tracePort == "" {
		tracePort = "19999"
	}

	_, err = strconv.ParseInt(tracePort, 10, 32)
	if err != nil {
		errChan <- fmt.Errorf("参数：traceport/tp错误：%w", err)
		return
	}

	addr := fmt.Sprintf("0.0.0.0:%s", tracePort)
	logger.New("trace").Infof("启动成功:pprof.web(addr:http://0.0.0.0:%s/debug/pprof/)", tracePort)
	err = http.ListenAndServe(addr, nil)
	if err != nil {
		errChan <- fmt.Errorf("启动pprof监控服务错误：%w", err)
	}
}
