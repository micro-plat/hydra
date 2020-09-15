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

	"github.com/micro-plat/hydra/context"
	"github.com/micro-plat/hydra/context/apm"
)

var errInvalidTracer = fmt.Errorf("invalid tracer")

type ClientConfig struct {
	apmCtx    context.IAPMContext
	name      string
	client    *http.Client
	extraTags map[string]string
}

// ClientOption allows optional configuration of Client.
type ClientOption func(*ClientConfig)

// WithClientOperationName override default operation name.
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

// newAPMClient returns an HTTP Client with tracer
func newAPMClient(apmctx context.IAPMContext, options ...ClientOption) (*http.Client, error) {
	co := &ClientConfig{
		apmCtx: apmctx,
	}
	for _, option := range options {
		option(co)
	}
	if co.client == nil {
		co.client = &http.Client{}
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

	if t.apmCtx == nil {
		return t.delegated.RoundTrip(req)
	}

	tracer := t.apmCtx.GetTracer()
	span, err := tracer.CreateExitSpan(t.apmCtx.GetRootCtx(), getOperationName(t.name, req), req.Host, func(header string) error {
		fmt.Println("CreateExitSpan:", req.URL.Host, req.URL.Port(), getOperationName("", req), header)
		req.Header.Set(apm.Header, header)
		return nil
	})
	if err != nil {
		return t.delegated.RoundTrip(req)
	}
	defer span.End()
	span.SetComponent(apm.ComponentIDGOHttpClient)
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
		return fmt.Sprintf("%s", r.URL.Path)
	}
	return name
}
