package pkgs

import (
	"fmt"
	"net/http"
	_ "net/http/pprof"
	"strconv"
	"strings"
	"time"

	"github.com/micro-plat/lib4go/logger"
	"github.com/pkg/profile"
)

type itrace interface {
	Stop()
}

var supportTraces = []string{"cpu", "mem", "block", "mutex", "web"}

//startTrace 启用项目性能跟踪
func startTrace(trace, tracePort string) (itrace, error) {
	switch strings.ToLower(trace) {
	case "cpu":
		return profile.Start(profile.CPUProfile, profile.ProfilePath("."), profile.NoShutdownHook), nil
	case "mem":
		return profile.Start(profile.MemProfile, profile.ProfilePath("."), profile.NoShutdownHook), nil
	case "block":
		return profile.Start(profile.BlockProfile, profile.ProfilePath("."), profile.NoShutdownHook), nil
	case "mutex":
		return profile.Start(profile.MutexProfile, profile.ProfilePath("."), profile.NoShutdownHook), nil
	case "web":
		web := &webTrace{port: tracePort}
		return web, web.Start()
	case "":
		return &emptyTrace{}, nil
	default:
		return nil, fmt.Errorf("不支持trace命令:%v", trace)
	}
	return &emptyTrace{}, nil
}

type webTrace struct {
	port string
}

func (w *webTrace) Start() error {
	errChan := make(chan error, 1)
	go startServer(w.port, errChan)
	select {
	case err := <-errChan:
		return err
	case <-time.After(time.Millisecond * 200):
		return nil
	}
	return nil
}

func startServer(tracePort string, errChan chan error) {

	if tracePort == "" {
		tracePort = "19999"
	}

	_, err := strconv.ParseInt(tracePort, 10, 32)
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

func (w *webTrace) Stop() {
	return
}

type emptyTrace struct {
}

func (e *emptyTrace) Stop() {
	return
}
