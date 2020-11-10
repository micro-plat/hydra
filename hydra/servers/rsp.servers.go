package servers

import (
	"fmt"
	"runtime/debug"
	"strings"
	"sync"
	"time"

	"github.com/micro-plat/hydra/conf/app"
	"github.com/micro-plat/hydra/global"
	"github.com/micro-plat/hydra/registry"
	"github.com/micro-plat/hydra/registry/watcher"
	"github.com/micro-plat/lib4go/logger"
)

//RspServers 响应式服务管理器,监控配置变更自动创建、停止服务器
type RspServers struct {
	registryAddr string
	registry     registry.IRegistry
	delayChan    chan string
	path         []string
	mpath        string
	notify       chan *watcher.ValueChangeArgs
	done         bool
	closeChan    chan struct{}
	log          logger.ILogger
	servers      map[string]IResponsiveServer
	lock         sync.Mutex
}

//NewRspServers 构建响应式服务器
func NewRspServers(registryAddr string, platName, sysName string, serverTypes []string, clusterName string) *RspServers {
	server := &RspServers{
		registryAddr: registryAddr,
		delayChan:    make(chan string, 10),
		closeChan:    make(chan struct{}),
		servers:      make(map[string]IResponsiveServer),
		log:          logger.New("hydra"),
		mpath:        registry.Join(platName, sysName, strings.Join(serverTypes, "-"), clusterName, "conf"),
	}
	for _, t := range serverTypes {
		server.path = append(server.path, registry.Join(platName, sysName, t, clusterName, "conf"))
	}
	return server
}

//Start 启动服务器
func (r *RspServers) Start() (err error) {

	r.log.Info("初始化:", r.mpath)

	//初始化注册中心
	r.registry, err = registry.NewRegistry(r.registryAddr, r.log)
	if err != nil {
		err = fmt.Errorf("注册中心初始化失败 %w", err)
		return
	}

	//监听配置变化
	watcher, err := watcher.NewValueWatcherByRegistry(r.registry, r.path, r.log)
	if err != nil {
		return fmt.Errorf("服务器watcher初始化失败 %s,%w", r.path, err)
	}

	//处理配置更变通知消息
	r.notify, err = watcher.Start()
	if err != nil {
		return err
	}
	go r.loopRecvNotify()
	return nil
}

//loopRecvNotify 接收注册中心配置变更消息
func (r *RspServers) loopRecvNotify() {
	go func() {
		tk := time.NewTicker(time.Second * 120)
	LOOP:
		for {
			select {
			case <-tk.C:
				debug.FreeOSMemory()
			case <-r.closeChan:
				break LOOP
			}
		}

	}()
	waitTimeout := time.After(time.Second)
LOOP:
	for {
		select {
		case <-r.closeChan:
			break LOOP
		case <-waitTimeout:
			if len(r.servers) == 0 {
				r.log.Debug("监听服务器配置...")
			}
		case p := <-r.delayChan:
			if r.done {
				break LOOP
			}
			if err := r.checkServer(p); err != nil {
				r.log.Error(err)
			}
		case u := <-r.notify:
			if r.done {
				break LOOP
			}
			if err := r.checkServer(u.Path); err != nil {
				r.log.Error(err)
			}
		}
	}
}

//checkServer 通知server配置变更或创建新server
func (r *RspServers) checkServer(path string) error {
	defer func() {
		if err := recover(); err != nil {
			r.log.Errorf("[Recovery] panic recovered:\n%s\n%s", err, global.GetStack())
		}
	}()
	//拉取配置信息
	conf, err := app.NewAPPConf(path, r.registry)
	if err != nil {
		r.log.Errorf("获取%s配置发生错误:%v", path, err)
	}

	//同一时间只允许一个流程处理配置变更
	r.lock.Lock()
	defer r.lock.Unlock()
	if r.done {
		return nil
	}

	//检查服务器是否已创建
	srvr, ok := r.servers[conf.GetServerConf().GetServerType()]
	if ok {
		//通知已创建的服务器
		r.log.Debugf("配置发生变化%s", conf.GetServerConf().GetServerPath())
		change, err := srvr.Notify(conf)
		if err != nil {
			return err
		}
		if !change {
			r.log.Debug("服务配置未发生变化")
		} else {
			r.log.Info("配置更新完成")
		}

	} else {
		//创建新服务器
		serverType := conf.GetServerConf().GetServerType()
		if creator, ok := creators[serverType]; ok {
			srvr, err := creator.Create(conf)
			if err != nil {
				return fmt.Errorf("[%s]服务器构建失败:%w", serverType, err)
			}
			r.log.Infof("启动[%s]服务...", serverType)
			if err := srvr.Start(); err != nil {
				r.delayPub(path)
				return fmt.Errorf("[%s]服务器启动失败:%w", serverType, err)
			}
			r.servers[serverType] = srvr
		} else {
			r.log.Errorf("服务器类型[%s]不支持或未注册", conf.GetServerConf().GetServerPath())
			return nil
		}
	}

	return nil

}

//delayPub 延迟启动，当依赖的服务没有正确启动时通过延迟重试进行启动
func (r *RspServers) delayPub(p string) {
	go func() {
		if r.done {
			return
		}
		time.Sleep(time.Second * 300)
		if r.done {
			return
		}
		r.delayChan <- p
	}()
}

//Shutdown 关闭所有服务器
func (r *RspServers) Shutdown() {
	r.done = true
	r.lock.Lock()
	defer r.lock.Unlock()
	cl := make(chan struct{})

	//新协程关闭服务器
	go func() {
		for _, server := range r.servers {
			server.Shutdown()
		}
		close(cl)
	}()

	//最长等待30秒
	select {
	case <-time.After(time.Second * 30):
		return
	case <-cl:
		return
	}
}
