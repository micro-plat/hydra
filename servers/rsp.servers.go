package servers

import (
	"fmt"
	"sync"
	"time"

	"github.com/micro-plat/hydra/application"
	"github.com/micro-plat/hydra/registry"
	"github.com/micro-plat/hydra/registry/conf/server"
	"github.com/micro-plat/hydra/registry/watcher"
	"github.com/micro-plat/lib4go/logger"
)

//RspServers 响应式服务管理器,监控配置变更自动创建、停止服务器
type RspServers struct {
	registryAddr string
	registry     registry.IRegistry
	path         []string
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
		closeChan:    make(chan struct{}),
		servers:      make(map[string]IResponsiveServer),
		log:          logger.New("hydra"),
	}
	for _, t := range serverTypes {
		server.path = append(server.path, registry.Join(platName, sysName, t, clusterName, "conf"))
	}
	return server
}

//Start 启动服务器
func (r *RspServers) Start() (err error) {

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

	//启动配置监听

	notify := make(chan struct{}, 1)
	go func() {
		f := time.After(time.Second * 10)
		select {
		case <-f:
			for _, p := range r.path {
				r.log.Infof("开始监听:%v", p)
			}
		case <-notify:
			break
		}
	}()
LOOP:
	for {
		select {
		case <-r.closeChan:
			break LOOP
		case u := <-r.notify:
			if r.done {
				break LOOP
			}
			if err := r.checkServer(u.Path); err != nil {
				r.log.Error(err)
			}
			select {
			case notify <- struct{}{}:
			default:
			}
		}
	}
	r.log.Error("exit")
}

//Shutdown 关闭所有服务器
func (r *RspServers) Shutdown() {
	r.done = true
	r.lock.Lock()
	defer r.lock.Unlock()
	cl := make(chan struct{})

	//新协程去关闭服务器
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

//checkServer 通知server配置变更或创建新server
func (r *RspServers) checkServer(path string) error {
	defer func() {
		if err := recover(); err != nil {
			r.log.Errorf("[Recovery] panic recovered:\n%s\n%s", err, application.GetStack())
		}
	}()
	//拉取配置信息
	conf, err := server.NewServerConf(path, r.registry)
	if err != nil {
		r.log.Errorf("加载[%s]配置发生错误:%v", path, err)
	}

	//同一时间只允许一个流程处理配置变更
	r.lock.Lock()
	defer r.lock.Unlock()
	if r.done {
		return nil
	}

	//检查服务器是否已创建
	srvr, ok := r.servers[conf.GetMainConf().GetServerType()]
	if ok {
		//通知已创建的服务器
		if err := srvr.Notify(conf); err != nil {
			return err
		}
	} else {
		//创建新服务器
		if creator, ok := creators[conf.GetMainConf().GetServerType()]; ok {
			srvr, err := creator.Create(conf)
			if err != nil {
				return fmt.Errorf("服务器%s %w", conf.GetMainConf().GetMainPath(), err)
			}
			r.servers[conf.GetMainConf().GetServerType()] = srvr
			if err := srvr.Start(); err != nil {
				return fmt.Errorf("服务器%s %w", conf.GetMainConf().GetMainPath(), err)
			}
		} else {
			r.log.Errorf("服务器类型[%s]不支持或未注册", conf.GetMainConf().GetMainPath())
			return nil
		}
	}

	//缓存服务器配置
	server.Cache.Save(conf)
	return nil

}
