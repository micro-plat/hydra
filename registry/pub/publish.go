package pub

import (
	"fmt"
	"net"
	"sync"
	"time"

	"github.com/micro-plat/hydra/registry/conf"
	"github.com/micro-plat/hydra/registry/conf/server/api"
	"github.com/micro-plat/hydra/registry/watcher"
	"github.com/micro-plat/lib4go/jsons"
	"github.com/micro-plat/lib4go/logger"
)

//IPublisher 服务发布程序
type IPublisher interface {
	WatchClusterChange(notify func(isMaster bool, sharding int, total int)) error
	Publish(serverName string, data map[string]interface{}, service ...string) error
	Clear()
	Close()
}

//Publisher 服务发布程序
type Publisher struct {
	c          conf.IMainConf
	log        logger.ILogging
	serverNode string
	serverName string
	lock       sync.Mutex
	closeChan  chan struct{}
	watchChan  chan struct{}
	pubs       map[string]string
	done       bool
}

//New 构建服务发布程序
func New(c conf.IMainConf) *Publisher {
	p := &Publisher{
		c:         c,
		watchChan: make(chan struct{}, 1),
		closeChan: make(chan struct{}),
		pubs:      make(map[string]string),
		log:       logger.New("publisher"),
	}
	go p.loopCheck()
	return p
}

//WatchClusterChange 监控集群服务节点变化
func (p *Publisher) WatchClusterChange(notify func(isMaster bool, sharding int, total int)) error {
	watcher, err := watcher.NewChildWatcherByRegistry(p.c.GetRegistry(), []string{p.c.GetServerPubPath()}, p.log)
	if err != nil {
		return err
	}

	//启动监控
	ch, err := watcher.Start()
	if err != nil {
		return err
	}

	//异步检查变化
	go func() {
	LOOP:
		for {
			select {
			case <-p.closeChan:
				watcher.Close()
				break LOOP
			case <-p.watchChan:
				watcher.Close()
				break LOOP
			case c := <-ch:
				total := p.c.GetMainConf().GetInt("sharding", 0)
				sharding, isMaster := GetSharding(true, total, p.serverNode, c.Children)
				notify(isMaster, sharding, total)
			}
		}
	}()
	return nil
}

//Publish 发布所有服务（集群节点，服务节点，DNS节点）
func (p *Publisher) Publish(serverName string, input map[string]interface{}, service ...string) error {

	buff, _ := jsons.Marshal(input)
	data := string(buff)
	if err := p.pubServerNode(serverName, data); err != nil {
		return err
	}
	switch p.c.GetServerType() {
	case "API", "WEB":
		if err := p.pubDNSNode(serverName); err != nil {
			return err
		}
		return p.pubAPIServiceNode(serverName, data)
	case "RPC":
		for _, srv := range service {
			if err := p.pubRPCServiceNode(serverName, srv, data); err != nil {
				return err
			}
		}
	}
	return nil
}

//pubRPCServiceNode 发布RPC服务节点
func (p *Publisher) pubRPCServiceNode(serverName string, service string, data string) error {
	path := fmt.Sprintf("%s/%s_", p.c.GetServicePubPathByService(service), serverName)
	npath, err := p.c.GetRegistry().CreateSeqNode(path, data)
	if err != nil {
		return fmt.Errorf("服务发布失败:(%s)[%v]", path, err)
	}
	p.appendPub(npath, data)
	return nil
}

//pubAPIServiceNode 发布API服务节点
func (p *Publisher) pubAPIServiceNode(serverName string, data string) error {
	path := fmt.Sprintf("%s/%s_", p.c.GetServicePubPath(), serverName)
	npath, err := p.c.GetRegistry().CreateSeqNode(path, data)
	if err != nil {
		return fmt.Errorf("服务发布失败:(%s)[%v]", path, err)
	}
	p.appendPub(npath, data)
	return nil
}

//pubDNSNode 发布DNS服务节点
func (p *Publisher) pubDNSNode(serverName string) error {
	//获取服务嚣配置
	server, err := api.GetConf(p.c)
	if err != nil {
		return err
	}
	if server.Domain == "" {
		return nil
	}

	//创建DNS节点
	ip, _, err := net.SplitHostPort(serverName)
	if err != nil {
		return err
	}
	path := fmt.Sprintf("/%s/%s", p.c.GetDNSPubPath(server.Domain), ip)
	err = p.c.GetRegistry().CreateTempNode(path, "")
	if err != nil {
		err = fmt.Errorf("DNS服务发布失败:(%s)[%v]", path, err)
		return err
	}
	p.appendPub(path, "")
	return nil
}

//pubServerNode 发布集群节点，用于服务监控
func (p *Publisher) pubServerNode(serverName string, data string) error {
	path := fmt.Sprintf("/%s/%s_%s_", p.c.GetServerPubPath(), p.c.GetClusterID(), serverName)

	npath, err := p.c.GetRegistry().CreateSeqNode(path, data)
	if err != nil {
		return fmt.Errorf("服务发布失败:(%s)[%v]", path, err)
	}
	p.serverNode = npath
	p.appendPub(npath, data)
	return nil
}
func (p *Publisher) appendPub(path string, data string) {
	p.lock.Lock()
	defer p.lock.Unlock()
	p.pubs[path] = data
}

//publishCheck 定时检查节点数据是否存在
func (p *Publisher) loopCheck() {
LOOP:
	for {
		select {
		case <-p.closeChan:
			break LOOP
		case <-time.After(time.Second * 30):
			if p.done {
				break LOOP
			}
			p.check()
		}
	}
}

//checkPubPath 检查已发布的节点，不存在则创建
func (p *Publisher) check() {
	p.lock.Lock()
	defer p.lock.Unlock()
	for path, data := range p.pubs {
		if p.done {
			break
		}
		ok, err := p.c.GetRegistry().Exists(path)
		if err != nil {
			break
		}
		if !ok {
			err := p.c.GetRegistry().CreateTempNode(path, data)
			if err != nil {
				break
			}
			p.log.Infof("节点(%s)已恢复", path)
		}
	}
}

//Close 关闭当前发布删除所有节点
func (p *Publisher) Close() {
	p.done = true
	close(p.closeChan)
	p.Clear()
}

//Clear 清除所有发布节点
func (p *Publisher) Clear() {
	p.lock.Lock()
	defer p.lock.Unlock()
	p.watchChan <- struct{}{}
	for path := range p.pubs {
		p.c.GetRegistry().Delete(path)
	}
	p.pubs = make(map[string]string)
}
