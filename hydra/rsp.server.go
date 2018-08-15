package hydra

import (
	"sync"
	"time"

	"github.com/micro-plat/hydra/component"

	"github.com/micro-plat/hydra/servers"

	"github.com/micro-plat/hydra/conf"
	"github.com/micro-plat/hydra/registry"
	"github.com/micro-plat/hydra/registry/watcher"
	"github.com/micro-plat/lib4go/logger"
)

type rspServer struct {
	servers      map[string]*server
	mu           sync.Mutex
	registry     registry.IRegistry
	registryAddr string
	logger       *logger.Logger
	handler      component.IComponentHandler
	done         bool
}

func newRspServer(registryAddr string, registry registry.IRegistry, handler component.IComponentHandler, logger *logger.Logger) *rspServer {
	return &rspServer{
		registry:     registry,
		registryAddr: registryAddr,
		servers:      make(map[string]*server),
		handler:      handler,
		logger:       logger,
	}
}

//Change 服务器发生变更
func (s *rspServer) Change(u *watcher.ContentChangeArgs) {
	if s.done {
		return
	}
	switch u.OP {
	case watcher.ADD, watcher.CHANGE:
		func() {
			s.mu.Lock()
			defer s.mu.Unlock()
			//初始化服务器配置
			conf, err := conf.NewServerConf(u.Path, u.Content, u.Version, s.registry)
			if err != nil {
				s.logger.Error(err)
				return
			}
			conf.Set("__component_handler_", s.handler)
			if _, ok := s.servers[u.Path]; !ok {
				//添加新服务器
				if conf.IsStop() {
					s.logger.Warnf("服务器(%s)配置为:stop", u.Path)
					return
				}
				server := newServer(conf, s.registryAddr, s.registry)
				server.logger.Infof("开始启动...")
				if err = server.Start(); err != nil {
					server.logger.Errorf("启动失败 %v", err)
					return
				}
				s.servers[u.Path] = server
				server.logger.Infof("启动成功(%s,%d)", server.GetAddress(), len(server.GetServices()))
			} else {
				//修改服务器
				server := s.servers[u.Path]
				if !conf.IsStop() {
					if err = server.Notify(conf); err != nil {
						server.logger.Errorf("未完成更新 %v", err)
					} else {
						if server.Restarted() {
							server.logger.Infof("配置更新成功(%s,%d)", server.GetAddress(), len(server.GetServices()))
						} else {
							server.logger.Info("配置更新成功")
						}
					}
				} else {
					server.logger.Warnf("服务器配置为:stop")
				}
				if conf.IsStop() || server.GetStatus() != servers.ST_RUNNING {
					server.logger.Debug("关闭服务器")
					server.Shutdown()
					delete(s.servers, u.Path)
					return
				}
			}
		}()

	case watcher.DEL:
		func() {
			s.mu.Lock()
			defer s.mu.Unlock()
			if server, ok := s.servers[u.Path]; ok {
				server.logger.Errorf("%s配置已删除", u.Path)
				server.Shutdown()
				server.logger.Info("服务器已关闭")
				delete(s.servers, u.Path)
				return
			}
		}()
	}
}

//Change 服务器发生变更
func (s *rspServer) Shutdown() {
	s.done = true
	s.mu.Lock()
	defer s.mu.Unlock()
	cl := make(chan struct{})
	go func() {
		for _, server := range s.servers {
			server.Shutdown()
		}
		close(cl)
	}()
	select {
	case <-time.After(time.Second * 30):
		return
	case <-cl:
		return
	}

}
