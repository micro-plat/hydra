package rpclog

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"

	"github.com/micro-plat/hydra/components"
)

type rpcWriter struct {
	service     string
	platName    string
	systemName  string
	serverType  string
	clusterName string
}

func newRPCWriter(service string, platName string, systemName string, clusterName string, serverType string) (r *rpcWriter) {
	return &rpcWriter{
		service:     service,
		platName:    platName,
		systemName:  systemName,
		clusterName: clusterName,
		serverType:  serverType,
	}
}
func (r *rpcWriter) Write(p []byte) (n int, err error) {
	if len(p) == 0 {
		return 0, nil
	}
	p[0] = byte('[')
	p = append(p, byte(']'))
	var buff bytes.Buffer
	if err := json.Compact(&buff, []byte(p)); err != nil {
		err = fmt.Errorf("json.compact.err:%v", err)
		return 0, err
	}

	_, err = components.Def.RPC().GetRegularRPC().Request(context.TODO(), r.service,
		map[string]interface{}{
			"__body_":     buff,
			"platName":    r.platName,
			"sysName":     r.systemName,
			"cluster":     r.clusterName,
			"server-type": r.serverType,
		})
	if err != nil {
		return 0, err
	}
	return len(p) - 1, nil
}
func (r *rpcWriter) Close() error {
	return nil
}
