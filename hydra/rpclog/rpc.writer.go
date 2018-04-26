package rpclog

import (
	"bytes"
	"encoding/json"
	"fmt"

	"github.com/golang/snappy"
	"github.com/micro-plat/hydra/rpc"
)

type rpcWriter struct {
	rpcInvoker *rpc.Invoker
	service    string
}

func newRPCWriter(service string, invoker *rpc.Invoker) (r *rpcWriter) {
	return &rpcWriter{service: service, rpcInvoker: invoker}
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

	dst := snappy.Encode(nil, buff.Bytes())
	_, _, _, err = r.rpcInvoker.Request(r.service, "GET", map[string]string{
		"__encode_snappy_": "true",
	}, map[string]string{
		"__body_": string(dst),
	}, true)
	if err != nil {
		return 0, err
	}
	return len(p) - 1, nil
}
func (r *rpcWriter) Close() error {
	if r.rpcInvoker != nil {
		r.rpcInvoker.Close()
	}
	return nil
}
