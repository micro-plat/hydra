package run

import (
	"fmt"
	"net/http"
	_ "net/http/pprof"
	"os"
	"runtime/trace"

	"github.com/micro-plat/lib4go/logger"
	"github.com/pkg/profile"
)

var supportTraces = []string{"cpu", "mem", "block", "mutex", "web"}

//startTrace 启用项目性能跟踪
func startTrace(trace string) error {
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
		go startTraceServer()
	case "":
		fallthrough
	default:
		return fmt.Errorf("不支持trace命令:%v", trace)
	}
	return nil
}
func startTraceServer() error {
	f, err := os.Create("trace.out")
	if err != nil {
		return err
	}
	defer f.Close()
	err = trace.Start(f)
	if err != nil {
		return err
	}
	defer trace.Stop()
	addr := "0.0.0.0:19999"
	logger.New("trace").Info("启动成功:pprof.web(addr:http://0.0.0.0:19999/debug/pprof/)")
	return http.ListenAndServe(addr, nil)
}
