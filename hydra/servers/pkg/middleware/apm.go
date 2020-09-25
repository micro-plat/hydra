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
	"strconv"
	"time"

	"github.com/micro-plat/hydra/components"
	"github.com/micro-plat/hydra/context/apm"
	"github.com/micro-plat/hydra/pkgs"
)

//APM 调用链
func APM() Handler {

	return func(ctx IMiddleContext) {
		//fmt.Println("middleware.apm")
		//获取apm配置
		apmconf := ctx.ServerConf().GetAPMConf()
		if apmconf.Disable || apmconf.GetName() == "" {
			ctx.Next()
			return
		}

		mainConf := ctx.ServerConf().GetMainConf()

		instance := fmt.Sprintf("%s_%s", mainConf.GetPlatName(), pkgs.LocalIP())
		apmInstance := components.Def.APM().GetRegularAPM(instance, apmconf.GetName())
		apmSvc := fmt.Sprintf("%s_%s", mainConf.GetPlatName(), mainConf.GetSysName())

		tracer, err := apmInstance.CreateTracer(apmSvc)
		if err != nil {
			ctx.Log().Warnf("APM.CreateTracer:%+v", err)
			ctx.Next()
			return
		}
		callback, err := procMiddle(ctx, tracer)
		if err != nil {
			ctx.Next()
			return
		}
		defer callback()
		ctx.Response().AddSpecial("apm")
		ctx.Next()
	}
}

//ProcMiddle ProcMiddle
func procMiddle(ctx IMiddleContext, tracer apm.Tracer) (callback func(), err error) {

	callback = func() {}

	oreq := ctx.Request()

	//fmt.Println("middleware.apm-2", tracer, err)
	span, rootctx, err := tracer.CreateEntrySpan(ctx.Context(), oreq.Path().GetURL(), func() (string, error) {
		return oreq.Path().GetHeader(apm.Header), nil
	})
	if err != nil {
		ctx.Log().Warnf("APM.CreateEntrySpan:%+v", err)
		return
	}

	ctx.StoreAPMCtx(apm.NewCtx(rootctx, tracer))

	//fmt.Println("middleware.apm-3", oreq.Header.Get("X-Request-Id"))
	span.SetComponent(apm.ComponentIDGOHttpServer)
	span.Tag("X-Request-Id", ctx.User().GetRequestID())
	// for k, v := range h.extraTags {
	//
	// }
	span.Tag(apm.TagHTTPMethod, oreq.Path().GetMethod())
	span.Tag(apm.TagURL, oreq.Path().GetURL())
	span.SetSpanLayer(apm.SpanLayer_Http)

	callback = func() {
		statusCode, v := ctx.Response().GetFinalResponse()
		if statusCode >= 400 {
			if len(v) > 300 {
				v = v[:300]
			}
			span.Error(time.Now(), "Error on handling request,code:"+strconv.Itoa(statusCode), v)
		}
		span.Tag(apm.TagStatusCode, strconv.Itoa(statusCode))
		span.End()
	}
	return
}
