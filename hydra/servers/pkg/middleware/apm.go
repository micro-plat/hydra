package middleware

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

import (
	"fmt"
	"net/http"
	"strconv"
	"time"

	"github.com/micro-plat/hydra/components/pkgs/apm"
	"github.com/micro-plat/hydra/global"
)

//APM 调用链
func APM() Handler {

	return func(ctx IMiddleContext) {
		//fmt.Println("middleware.apm")
		//获取apm配置
		apmconf := ctx.ServerConf().GetAPMConf()
		if apmconf.Disable {
			ctx.Next()
			return
		}
		ctx.Response().AddSpecial("apm")

		octx := ctx.Context()
		//fmt.Println("middleware.apm-1", octx)
		oreq, _ := ctx.GetHttpReqResp()

		tracer, err := global.Def.APM.CreateTracer(global.Def.GetAPMService())
		if err != nil {
			ctx.Log().Warnf("APM.CreateTracer:%+v", err)
			ctx.Next()
			return
		}
		//fmt.Println("middleware.apm-2", tracer, err)
		span, _, err := tracer.CreateEntrySpan(octx, getOperationName("", oreq), func() (string, error) {
			return oreq.Header.Get(apm.Header), nil
		})
		if err != nil {
			ctx.Log().Warnf("APM.CreateEntrySpan:%+v", err)
			ctx.Next()
			return
		}
		//fmt.Println("middleware.apm-3", oreq.Header.Get("X-Request-Id"))
		span.SetComponent(componentIDGOHttpServer)
		span.Tag("X-Request-Id", oreq.Header.Get("X-Request-Id"))
		// for k, v := range h.extraTags {
		//
		// }
		span.Tag(apm.TagHTTPMethod, oreq.Method)
		span.Tag(apm.TagURL, fmt.Sprintf("%s%s", oreq.Host, oreq.URL.Path))
		span.SetSpanLayer(apm.SpanLayer_Http)

		defer func() {
			statusCode, _ := ctx.Response().GetRawResponse()
			code := statusCode
			if code >= 400 {
				span.Error(time.Now(), "Error on handling request")
			}
			//fmt.Println("middleware.apm-4", statusCode)
			span.Tag(apm.TagStatusCode, strconv.Itoa(code))
			span.End()
		}()

		ctx.Next()
	}
}

const componentIDGOHttpServer = 5004

func getOperationName(name string, r *http.Request) string {
	if name == "" {
		return fmt.Sprintf("/%s%s", r.Method, r.URL.Path)
	}
	return name
}
