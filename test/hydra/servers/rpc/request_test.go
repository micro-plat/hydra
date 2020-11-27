package rpc

import (
	"encoding/json"
	"strings"
	"testing"

	"github.com/micro-plat/hydra/components/rpcs/rpc/pb"
	"github.com/micro-plat/hydra/hydra/servers/rpc"
	"github.com/micro-plat/hydra/test/assert"
)

func TestNewRequest(t *testing.T) {
	type args struct {
		Service string
		Method  string
		Header  string
		Input   string
	}
	tests := []struct {
		name    string
		args    args
		wantErr string
	}{
		{name: "1. 初始化入参是空对象", args: args{Service: "", Method: "", Header: "", Input: ""}, wantErr: "rpc请求头转换失败"},
		{name: "2. 初始化入参是错误的header", args: args{Service: "service1", Method: "get", Header: "错误的header", Input: `{"taosy":"test"}`}, wantErr: "rpc请求头转换失败"},
		{name: "3. 初始化入参是错误的input", args: args{Service: "service1", Method: "get", Header: `{"Host":"www.baidu.com"}`, Input: "错误的input"}, wantErr: "rpc请求参数转换失败"},
		{name: "4. 初始化正确的参数对象", args: args{Service: "service1", Method: "get", Header: `{"Host":"www.baidu.com"}`, Input: `{"taosy":"test"}`}, wantErr: ""},
	}
	for _, tt := range tests {
		request := &pb.RequestContext{
			Service: tt.args.Service,
			Method:  tt.args.Method,
			Header:  tt.args.Header,
			Input:   tt.args.Input,
		}

		gotR, err := rpc.NewRequest(request)
		if tt.wantErr != "" {
			assert.Equalf(t, true, strings.Contains(err.Error(), tt.wantErr), tt.name, tt.wantErr, err)
			break
		}

		name := gotR.GetName()
		assert.Equalf(t, tt.args.Service, name, tt.name, tt.args.Service, name)

		host := gotR.GetHost()
		if strings.Contains(tt.args.Header, "Host") {
			assert.Equalf(t, true, strings.Contains(tt.args.Header, host), tt.name, tt.args.Service, name)
		}

		method := gotR.GetMethod()
		assert.Equalf(t, tt.args.Method, method, tt.name, tt.args.Method, method)

		service := gotR.GetService()
		assert.Equalf(t, tt.args.Service, service, tt.name, tt.args.Service, service)

		headMap := map[string]string{}
		err = json.Unmarshal([]byte(request.Header), &headMap)
		assert.Equalf(t, true, err == nil, tt.name, err)
		hm := gotR.GetHeader()
		assert.Equalf(t, headMap, hm, tt.name, headMap, hm)

		fromMap := map[string]interface{}{}
		err = json.Unmarshal([]byte(request.Input), &fromMap)
		assert.Equalf(t, true, err == nil, tt.name, err)
		fm := gotR.GetForm()
		assert.Equalf(t, fromMap, fm, tt.name, fromMap, fm)
	}
}
