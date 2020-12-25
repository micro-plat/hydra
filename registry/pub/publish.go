package pub

import (
	"fmt"
	"net"
	"strings"
	"sync"
	"time"

	"github.com/micro-plat/hydra/conf"
	"github.com/micro-plat/hydra/conf/server/api"
	"github.com/micro-plat/hydra/global"
	"github.com/micro-plat/hydra/registry"
	"github.com/micro-plat/hydra/registry/watcher"
	"github.com/micro-plat/lib4go/jsons"
	"github.com/micro-plat/lib4go/logger"
	"github.com/micro-plat/lib4go/types"
)

//IPublisher 服务发布程序
type IPublisher interface {
	//WatchClusterChange(notify func(isMaster bool, sharding int, total int)) error
	Publish(serverName string, serviceAddr string, clusterID string, service ...string) error
	Update(serverName string, serviceAddr string, clusterID string, kv ...string) error
	Clear()
	Close()
}

//Publisher 服务发布程序
type Publisher struct {
	c          conf.IServerConf
	log        logger.ILogging
	serverNode string
	serverName string
	lock       sync.Mutex
	closeChan  chan struct{}
	watchChan  chan struct{}
	pubs       map[string]string
	done       bool
	checkTime  time.Duration
}

//New 构建服务发布程序
func New(c conf.IServerConf, checkTime ...time.Duration) *Publisher {
	p := &Publisher{
		c:         c,
		watchChan: make(chan struct{}),
		closeChan: make(chan struct{}),
		pubs:      make(map[string]string),
		log:       logger.New("publisher"),
	}
	if len(checkTime) > 0 {
		p.checkTime = checkTime[0]
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
				sharding, isMaster := global.IsMaster(true, total, p.serverNode, c.Children)
				notify(isMaster, sharding, total)
			}
		}
	}()
	return nil
}

//Publish 发布所有服务（集群节点，服务节点，DNS节点）
func (p *Publisher) Publish(serverName string, serviceAddr string, clusterID string, service ...string) error {
	input := map[string]interface{}{}
	input["addr"] = serviceAddr
	input["cluster_id"] = clusterID
	input["time"] = time.Now().Unix()
	buff, err := jsons.Marshal(input)
	if err != nil {
		return fmt.Errorf("服务器发布数据转换为json失败:%w", err)
	}
	data := string(buff)
	if _, err := p.PubServerNode(serverName, data); err != nil {
		return err
	}
	switch p.c.GetServerType() {
	case global.API, global.Web:
		if _, err := p.PubDNSNode(serverName, serviceAddr); err != nil {
			return err
		}
		_, err := p.PubAPIServiceNode(serverName, data)
		return err
	case global.RPC:
		for _, srv := range service {
			if _, err := p.PubRPCServiceNode(serverName, srv, data); err != nil {
				return err
			}
		}
	}
	return nil
}

//Update 更新服务器配置
func (p *Publisher) Update(serverName string, serviceAddr string, clusterID string, kv ...string) error {
	input := map[string]interface{}{}
	input["addr"] = serviceAddr
	input["cluster_id"] = clusterID
	input["time"] = time.Now().Unix()

	if len(kv)%2 > 0 {
		return fmt.Errorf("更新服务器发布数据,展参数必须成对出现：%d", len(kv))
	}
	for i := 0; i+1 < len(kv); i = i + 2 {
		input[kv[i]] = kv[i+1]
	}

	buff, err := jsons.Marshal(input)
	if err != nil {
		return fmt.Errorf("更新服务器发布数据失败:%w", err)
	}
	ndata := string(buff)
	p.lock.Lock()
	defer p.lock.Unlock()
	for path := range p.pubs {
		if p.done {
			break
		}
		if strings.Contains(path, serverName) {
			p.pubs[path] = ndata
		}
		err := p.c.GetRegistry().Update(path, ndata)
		if err != nil {
			return err
		}
	}
	return nil
}

//PubRPCServiceNode 发布RPC服务节点
func (p *Publisher) PubRPCServiceNode(serverName string, service string, data string) (map[string]string, error) {
	path := registry.Join(p.c.GetRPCServicePubPath(service), serverName+"_")
	npath, err := p.c.GetRegistry().CreateSeqNode(path, data)
	if err != nil {
		return nil, fmt.Errorf("服务发布失败:(%s)[%v]", path, err)
	}
	p.appendPub(npath, data)
	return p.pubs, nil
}

//PubAPIServiceNode 发布API服务节点
func (p *Publisher) PubAPIServiceNode(serverName string, data string) (map[string]string, error) {
	path := registry.Join(p.c.GetServicePubPath(), serverName+"_")
	npath, err := p.c.GetRegistry().CreateSeqNode(path, data)
	if err != nil {
		return nil, fmt.Errorf("服务发布失败:(%s)[%v]", path, err)
	}
	p.appendPub(npath, data)
	return p.pubs, nil
}

//PubDNSNode 发布DNS服务节点
func (p *Publisher) PubDNSNode(serverName string, serviceAddr string) (map[string]string, error) {
	//获取服务嚣配置
	server, err := api.GetConf(p.c)
	if err != nil {
		return nil, err
	}
	if server.Domain == "" {
		return p.pubs, nil
	}
	proto, addr, err := global.ParseProto(serviceAddr)
	if err != nil {
		return nil, err
	}
	host, port, err := net.SplitHostPort(addr)
	if err != nil {
		return nil, err
	}

	input := map[string]interface{}{
		"plat_name":       p.c.GetPlatName(),
		"plat_cn_name":    global.Def.PlatCNName,
		"system_name":     p.c.GetSysName(),
		"system_cn_name":  types.GetString(server.Name, global.Def.SysCNName),
		"server_type":     p.c.GetServerType(),
		"cluster_name":    p.c.GetClusterName(),
		"server_name":     serverName,
		"service_address": serviceAddr,
		"proto":           proto,
		"host":            host,
		"port":            port,
		"ip":              global.LocalIP(),
	}
	buff, err := jsons.Marshal(input)
	if err != nil {
		return nil, fmt.Errorf("更新dns服务器发布数据失败:%w", err)
	}
	ndata := string(buff)
	domain := strings.TrimPrefix(server.Domain, "www.")
	path := registry.Join(p.c.GetDNSPubPath(domain), fmt.Sprintf("%s:%s", host, port))
	exist, err := p.c.GetRegistry().Exists(path)
	if err != nil {
		err = fmt.Errorf("DNS服务发布失败:(%s)[%v]", path, err)
		return nil, err
	}

	if exist {
		err = p.c.GetRegistry().Update(path, ndata)
	} else {
		err = p.c.GetRegistry().CreateTempNode(path, ndata)
	}

	if err != nil {
		err = fmt.Errorf("DNS服务发布失败:(%s)[%v]", path, err)
		return nil, err
	}

	//加入节点检查
	p.appendPub(path, ndata)
	return p.pubs, nil
}

//PubServerNode 发布集群节点，用于服务监控
func (p *Publisher) PubServerNode(serverName string, data string) (map[string]string, error) {
	path := registry.Join(p.c.GetServerPubPath(), fmt.Sprintf("%s_%s_", serverName, p.c.GetServerID()))
	npath, err := p.c.GetRegistry().CreateSeqNode(path, data)
	if err != nil {
		return nil, fmt.Errorf("服务发布失败:(%s)[%v]", path, err)
	}
	p.serverNode = npath
	p.appendPub(npath, data)
	return p.pubs, nil
}

func (p *Publisher) appendPub(path string, data string) {
	p.lock.Lock()
	defer p.lock.Unlock()
	p.pubs[path] = data
}

//publishCheck 定时检查节点数据是否存在
func (p *Publisher) loopCheck() {
	checkTime := time.Second * 30
	if p.checkTime > 0 {
		checkTime = p.checkTime
	}
LOOP:
	for {
		select {
		case <-p.closeChan:
			break LOOP
		case <-time.After(checkTime):
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
		} else {
			err := p.c.GetRegistry().Update(path, data)
			if err != nil {
				break
			}
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
	close(p.watchChan)
	for path := range p.pubs {
		p.c.GetRegistry().Delete(path)
	}
	p.pubs = make(map[string]string)
	p.watchChan = make(chan struct{})
}
