package hydra

import (
	"fmt"
	"os"
	"os/signal"
	"runtime/debug"
	"strings"
	"sync"
	"syscall"
	"time"

	"github.com/micro-plat/hydra/component"
	"github.com/micro-plat/hydra/hydra/rpclog"
	"github.com/micro-plat/hydra/servers"

	"github.com/micro-plat/hydra/registry"
	"github.com/micro-plat/hydra/registry/watcher"
	"github.com/micro-plat/lib4go/logger"
)

//Hydra  hydra app
type Hydra struct {
	logger         *logger.Logger
	closeChan      chan struct{}
	interrupt      chan os.Signal
	isDebug        bool
	platName       string
	systemName     string
	serverTypes    []string
	clusterName    string
	registryAddr   string
	systemRootName string
	rpcLoggerPath  string
	mu             sync.Mutex
	registry       registry.IRegistry
	watcher        *watcher.ConfWatcher
	notify         chan *watcher.ContentChangeArgs
	cHandler       component.IComponentHandler
	rspServer      *rspServer
	trace          string
	done           bool
	remoteLogger   bool
}

//NewHydra 创建hydra服务器
func NewHydra(platName string, systemName string, serverTypes []string, clusterName string, trace string, registryAddr string, isDebug bool, remoteLogger bool, logger *logger.Logger, r component.IComponentHandler) *Hydra {
	servers.IsDebug = isDebug
	return &Hydra{
		cHandler:       r,
		logger:         logger,
		systemRootName: registry.Join("/", platName, systemName, strings.Join(serverTypes, "-"), clusterName),
		rpcLoggerPath:  registry.Join("/", platName, "/var/global/logger"),
		//	binder:         binder,
		closeChan:    make(chan struct{}),
		interrupt:    make(chan os.Signal, 1),
		isDebug:      isDebug,
		platName:     platName,
		systemName:   systemName,
		serverTypes:  serverTypes,
		clusterName:  clusterName,
		registryAddr: registryAddr,
		remoteLogger: remoteLogger,
		trace:        trace,
	}
}

//Start 启动hydra服务器
func (h *Hydra) Start() (s string, err error) {
	//非调试模式时设置日志写协程数为50个
	if !h.isDebug {
		logger.AddWriteThread(49)
	}
	if h.remoteLogger {
		_, err := rpclog.NewRPCLogger(h.rpcLoggerPath, h.registryAddr, h.logger, h.platName, h.systemName, h.clusterName, h.serverTypes)
		if err != nil {
			return "", err
		}
	}
	//创建trace性能跟踪
	if err = startTrace(h.trace, h.logger); err != nil {
		return
	}

	//开始监控服务器配置变更
	if err = h.startWatch(); err != nil {
		return "", err
	}
	//定时释放内存

	go h.freeMemory()

	//堵塞当前进程，等服务结束
	signal.Notify(h.interrupt, os.Interrupt, os.Kill, syscall.SIGTERM) //, syscall.SIGUSR1) //9:kill/SIGKILL,15:SIGTEM,20,SIGTOP 2:interrupt/syscall.SIGINT
LOOP:
	for {
		select {
		case <-h.interrupt:
			h.done = true
			break LOOP
		}
	}
	h.logger.Infof("hydra 正在退出...")
	h.rspServer.Shutdown()
	h.logger.Infof("hydra 已安全退出")
	return "hydra 已安全退出", nil
}

//startWatch 启动服务器配置监控
func (h *Hydra) startWatch() (err error) {
	//创建注册中心
	if h.registry, err = registry.NewRegistryWithAddress(h.registryAddr, h.logger); err != nil {
		err = fmt.Errorf("注册中心初始化失败:%s(%v)", h.registryAddr, err)
		return
	}

	//启动配置监听
	h.watcher, err = watcher.NewConfWatcher(h.platName, h.systemName, h.serverTypes, h.clusterName, h.registry, h.logger)
	if err != nil {
		err = fmt.Errorf("watcher初始化失败 %s,%+v", registry.Join(h.platName, h.systemName), err)
		return
	}
	h.logger.Infof("初始化 %s", h.systemRootName)
	if h.notify, err = h.watcher.Notify(); err != nil {
		return err
	}
	if err != nil {
		err = fmt.Errorf("watcher启动失败 %s,%+v", registry.Join(h.platName, h.systemName), err)
		return
	}

	//创建服务管理器
	h.rspServer = newRspServer(h.registryAddr, h.registry, h.cHandler, h.logger)

	//循环接收服务变更新通知
	go h.loopRecvNotify()
	return nil
}

//freeMemory 每120秒执行1次垃圾回收，清理内存
func (h *Hydra) freeMemory() {
	for {
		select {
		case <-h.closeChan:
			return
		case <-time.After(time.Second * 120):
			debug.FreeOSMemory()
		}
	}
}

//循环接收服务器配置变化通知
func (h *Hydra) loopRecvNotify() {
	notify := make(chan struct{}, 1)
	go func() {
		select {
		case <-time.After(time.Second * 10):
			h.logger.Warnf("%s 未配置", h.systemRootName)
		case <-notify:
			break
		}
	}()
LOOP:
	for {
		select {
		case <-h.closeChan:
			break LOOP
		case u := <-h.notify:
			if h.done {
				break LOOP
			}
			h.rspServer.Change(u)
			select {
			case notify <- struct{}{}:
			default:
			}
		}
	}
}

func (h *Hydra) Shutdown() {
	h.done = true
	close(h.closeChan)
	h.interrupt <- syscall.SIGTERM

}
