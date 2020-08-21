// Licensed to SkyAPM org under one or more contributor
// license agreements. See the NOTICE file distributed with
// this work for additional information regarding copyright
// ownership. SkyAPM org licenses this file to you under
// the Apache License, Version 2.0 (the "License"); you may
// not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing,
// software distributed under the License is distributed on an
// "AS IS" BASIS, WITHOUT WARRANTIES OR CONDITIONS OF ANY
// KIND, either express or implied.  See the License for the
// specific language governing permissions and limitations
// under the License.

package http

import (
	"fmt"
	"net/http"
	"strconv"
	"time"

	"github.com/micro-plat/hydra/components/pkgs/apm"
	"github.com/micro-plat/hydra/global"
)

var errInvalidTracer = fmt.Errorf("invalid tracer")

const componentIDGOHttpClient = 5005

type ClientConfig struct {
	name      string
	client    *http.Client
	extraTags map[string]string
	apmInfo   *apm.APMInfo
}

// ClientOption allows optional configuration of Client.
type ClientOption func(*ClientConfig)

// WithOperationName override default operation name.
func WithClientOperationName(name string) ClientOption {
	return func(c *ClientConfig) {
		c.name = name
	}
}

// WithClientTag adds extra tag to client spans.
func WithClientTag(key string, value string) ClientOption {
	return func(c *ClientConfig) {
		if c.extraTags == nil {
			c.extraTags = make(map[string]string)
		}
		c.extraTags[key] = value
	}
}

// WithClient set customer http client.
func WithClient(client *http.Client) ClientOption {
	return func(c *ClientConfig) {
		c.client = client
	}
}

// newTracerClient returns an HTTP Client with tracer
func newTracerClient(apmInfo *apm.APMInfo, options ...ClientOption) (*http.Client, error) {
	co := &ClientConfig{
		apmInfo: apmInfo,
	}
	for _, option := range options {
		option(co)
	}
	if co.client == nil {
		co.client = &http.Client{}
	}
	if !global.Def.IsUseAPM() {
		return co.client, nil
	}

	tp := &transport{
		ClientConfig: co,
		delegated:    http.DefaultTransport,
	}
	if co.client.Transport != nil {
		tp.delegated = co.client.Transport
	}
	co.client.Transport = tp
	return co.client, nil
}

type transport struct {
	*ClientConfig
	delegated http.RoundTripper
}

func (t *transport) RoundTrip(req *http.Request) (res *http.Response, err error) {
	apmInfo := t.apmInfo
	rootCtx := apmInfo.RootCtx
	tracer := apmInfo.Tracer

	span, err := tracer.CreateExitSpan(rootCtx, getOperationName(t.name, req), req.Host, func(header string) error {
		fmt.Println("CreateExitSpan:", req.URL.Host, req.URL.Port(), getOperationName("", req), header)
		req.Header.Set(apm.Header, header)
		return nil
	})
	if err != nil {
		return t.delegated.RoundTrip(req)
	}
	defer span.End()
	span.SetComponent(componentIDGOHttpClient)
	for k, v := range t.extraTags {
		span.Tag(k, v)
	}
	span.Tag(apm.TagHTTPMethod, req.Method)
	span.Tag(apm.TagURL, req.URL.String())
	span.SetSpanLayer(apm.SpanLayer_Http)
	res, err = t.delegated.RoundTrip(req)
	if err != nil {
		span.Error(time.Now(), err.Error())
		return
	}
	span.Tag(apm.TagStatusCode, strconv.Itoa(res.StatusCode))
	if res.StatusCode >= http.StatusBadRequest {
		span.Error(time.Now(), "Errors on handling client")
	}
	return res, nil
}

func getOperationName(name string, r *http.Request) string {
	if name == "" {
		return fmt.Sprintf("/%s%s", r.Method, r.URL.Path)
	}
	return name
}
