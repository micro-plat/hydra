package servers

import (
	"fmt"
	"sync"
	"time"

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
		err = fmt.Errorf("注册中心初始化失败:%s(%v)", r.registryAddr, err)
		return
	}

	//监听配置变化
	watcher, err := watcher.NewValueWatcherByRegistry(r.registry, r.path, r.log)
	if err != nil {
		return fmt.Errorf("服务器watcher初始化失败 %s,%w", r.path, err)
	}

	//启动配置监听
	r.notify, err = watcher.Start()
	if err != nil {
		return err
	}
	go r.loopRecvNotify()
	return nil
}

//loopRecvNotify 接收注册中心配置变更消息
func (r *RspServers) loopRecvNotify() {
	notify := make(chan struct{}, 1)
	go func() {
		select {
		case <-time.After(time.Second * 10):
			r.log.Warnf("%s 未配置", r.path[0])
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
}

//Shutdown 关闭所有服务器
func (r *RspServers) Shutdown() {
	r.done = true
	r.lock.Lock()
	defer r.lock.Unlock()
	cl := make(chan struct{})

	//多个协程去关闭服务器
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

	//拉取配置信息
	conf, err := server.NewServerConf(path, r.registry)
	if err != nil {
		r.log.Error("加载[%s]配置发生错误:%w", path, err)
	}

	//拿到权限再去检查服务器配置
	r.lock.Lock()
	defer r.lock.Unlock()
	if r.done {
		return nil
	}

	//通知已创建的服务器
	if server, ok := r.servers[conf.GetMainConf().GetServerType()]; ok {
		return server.Notify(conf)
	}

	//创建新服务器
	if creator, ok := creators[conf.GetMainConf().GetServerType()]; ok {
		server, err := creator.Create(conf)
		if err != nil {
			return err
		}
		r.servers[conf.GetMainConf().GetServerType()] = server
	}
	return nil

}
