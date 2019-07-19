package rqs

import (
	"errors"
	"fmt"
	"strings"
	"time"

	"github.com/micro-plat/hydra/context"
	"github.com/micro-plat/hydra/registry"
	"github.com/micro-plat/hydra/servers"

	"github.com/micro-plat/hydra/servers/http"
	"github.com/micro-plat/lib4go/jsons"
	"github.com/micro-plat/lib4go/logger"
	"github.com/micro-plat/lib4go/net"
	"github.com/micro-plat/lib4go/sysinfo/cpu"
	"github.com/micro-plat/lib4go/sysinfo/disk"
	"github.com/micro-plat/lib4go/sysinfo/memory"
)

var statusLocalPort = []int{10160, 10162, 10166, 10168}

//RemoteQueryService 远程查询服务
type RemoteQueryService struct {
	server        *http.ApiServer
	platName      string
	systemName    string
	serverTypes   []string
	clusterName   string
	closeChan     chan struct{}
	registry      registry.IRegistry
	logger        *logger.Logger
	version       string
	pubs          []string
	HydraShutdown func()
	done          bool
}

//NewHRemoteQueryService 创建HRemoteQueryService
func NewHRemoteQueryService(platName string, systemName string, serverTypes []string, clusterName string, registry registry.IRegistry, version string) (h *RemoteQueryService, err error) {
	port := net.GetAvailablePort(statusLocalPort)
	h = &RemoteQueryService{
		platName:    platName,
		systemName:  systemName,
		serverTypes: serverTypes,
		clusterName: clusterName,
		logger:      logger.GetSession("RQS", logger.CreateSession()),
		version:     version,
		registry:    registry,
		closeChan:   make(chan struct{}),
	}
	routers := http.GetRouters()
	routers.Route("get", "/server/query", h.queryHandler())
	routers.Route("get", "/update/:v", h.updateHandler())
	server, err := http.NewApiServer("RQS", fmt.Sprintf(":%d", port), routers.Get())
	if err != nil {
		return nil, err
	}
	h.server = server
	return h, nil

}
func (h *RemoteQueryService) Start() error {
	h.logger.Info("开始启动...")
	if err := h.server.Run(); err != nil {
		return err
	}
	if err := h.publish(); err != nil {
		h.logger.Errorf("启动失败 %v", err)
		return err
	}
	h.logger.Infof("启动成功(%s)", h.server.GetAddress())
	return nil
}

func (h *RemoteQueryService) Shutdown() {
	close(h.closeChan)
	h.done = true
	h.unpublish()
	h.server.Shutdown(time.Second)
}

//publish 将当前服务器的节点信息发布到注册中心
func (h *RemoteQueryService) publish() (err error) {
	addr := h.server.GetAddress()
	ipPort := strings.Split(addr, "://")[1]
	data := map[string]interface{}{
		"address":      h.server.GetAddress(),
		"time":         time.Now().Unix(),
		"plat-name":    h.platName,
		"system-name":  h.systemName,
		"server-type":  h.serverTypes,
		"cluster-name": h.clusterName,
		"version":      h.version,
	}
	jsonData, _ := jsons.Marshal(data)
	nodeData := string(jsonData)
	h.pubs = []string{}
	pubPath := registry.Join("/", h.platName, "services", "rcs", ipPort)
	r, err := h.registry.CreateSeqNode(pubPath, nodeData)
	if err != nil {
		err = fmt.Errorf("服务发布失败:(%s)[%v]", pubPath, err)
		return err
	}
	h.pubs = append(h.pubs, r)

	go h.publishCheck(nodeData)
	return
}

//publishCheck 定时检查节点数据是否存在
func (h *RemoteQueryService) publishCheck(data string) {
LOOP:
	for {
		select {
		case <-h.closeChan:
			break LOOP
		case <-time.After(time.Second * 10):
			if h.done {
				break LOOP
			}
			h.checkPubPath(data)
		}
	}
}

//checkPubPath 检查已发布的节点，不存在则创建
func (h *RemoteQueryService) checkPubPath(data string) {
	for _, path := range h.pubs {
		if h.done {
			break
		}
		ok, err := h.registry.Exists(path)
		if err != nil {
			break
		}
		if !ok {
			err := h.registry.CreateTempNode(path, data)
			if err != nil {
				break
			}
			h.logger.Infof("节点(%s)已恢复", path)
		}
	}
}

//unpublish 删除已发布的节点
func (h *RemoteQueryService) unpublish() {
	for _, path := range h.pubs {
		h.registry.Delete(path)
	}
	h.pubs = make([]string, 0, 0)
}

func (h *RemoteQueryService) queryHandler() servers.IExecuteHandler {
	return func(ctx *context.Context) (rs interface{}) {
		ctx.Response.SetJSON()
		data := make(map[string]interface{})
		data["cpu_used_precent"] = fmt.Sprintf("%.2f", cpu.GetInfo(time.Millisecond*200).UsedPercent)
		data["mem_used_precent"] = fmt.Sprintf("%.2f", memory.GetInfo().UsedPercent)
		data["disk_used_precent"] = fmt.Sprintf("%.2f", disk.GetInfo().UsedPercent)
		data["app_memory"] = memory.GetAPPMemory()
		data["plat-name"] = h.platName
		data["system-name"] = h.systemName
		data["server-type"] = h.serverTypes
		data["cluster-name"] = h.clusterName
		data["version"] = h.version
		return data
	}
}
func (h *RemoteQueryService) updateHandler() servers.IExecuteHandler {
	return func(ctx *context.Context) (rs interface{}) {
		ctx.Response.SetJSON()
		v := ctx.Request.Param.GetString("v")
		if v == "" {
			return errors.New("未指定版本号")
		}
		b, pkg, err := NeedUpdate(h.registry, h.platName, h.systemName, v)
		if err != nil {
			return err
		}
		if b {
			if err = UpdateNow(pkg, h.logger, func() {
				if h.HydraShutdown != nil {
					h.HydraShutdown()
				}
				//关闭服务器
			}); err != nil {
				return err
			}
			return nil
		}
		ctx.Response.SetStatus(202)
		return nil
	}
}
