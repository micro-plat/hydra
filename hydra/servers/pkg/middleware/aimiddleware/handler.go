package aimiddleware

import (
	"context"
	"fmt"
	"io"
	"net/http"
	"time"

	"github.com/micro-plat/hydra/components"
	aigwconf "github.com/micro-plat/hydra/conf/server/aigw"
	hctx "github.com/micro-plat/hydra/context"
	"github.com/micro-plat/hydra/global"
	"github.com/micro-plat/hydra/hydra/servers/pkg/middleware"
	"github.com/micro-plat/hydra/pkgs"
	"github.com/micro-plat/hydra/services"
)

// ExecuteHandler AIGW业务处理Handler。
func ExecuteHandler() middleware.Handler {
	return func(ctx middleware.IMiddleContext) {
		service := ctx.Request().Path().GetService()
		if ctx.Request().Path().IsLimited() {
			fallback(ctx, service)
			return
		}
		if addr, ok := global.IsProto(service, global.ProtoRPC); ok {
			response, err := components.Def.RPC().GetRegularRPC().Swap(addr, ctx)
			if err != nil {
				writeOpenAIError(ctx, rpcErrorStatus(response), err)
				return
			}
			if response == nil {
				writeOpenAIError(ctx, http.StatusBadGateway, fmt.Errorf("rpc response is nil"))
				return
			}
			headers := response.GetHeaders()
			for k := range headers {
				ctx.Response().Header(k, headers.GetString(k))
			}
			ctx.Response().Write(response.GetStatus(), response.GetResult())
			return
		}

		serverType := ctx.APPConf().GetServerConf().GetServerType()
		method := ctx.Request().Path().GetMethod()
		if !services.Def.Has(serverType, service, method) {
			writeOpenAIError(ctx, http.StatusNotFound, fmt.Errorf("not found path %s", ctx.Request().Path().GetRequestPath()))
			return
		}

		result := services.Def.Call(ctx, service)
		switch v := result.(type) {
		case RawJSON:
			writeRawJSON(ctx, v)
		case *RawJSON:
			if v == nil {
				return
			}
			writeRawJSON(ctx, *v)
		case SSEProxy:
			writeSSEProxy(ctx, v)
		case *SSEProxy:
			if v == nil {
				return
			}
			writeSSEProxy(ctx, *v)
		default:
			if ok, r := hctx.IsSSEData(result); ok {
				if stream, ok := r.(interface {
					LoopWriteWithContext(context.Context, http.ResponseWriter)
				}); ok {
					stream.LoopWriteWithContext(ctx.Request().GetHTTPRequest().Context(), ctx.Response().GetHTTPReponse())
					return
				}
				r.LoopWrite(ctx.Response().GetHTTPReponse())
				return
			}
			ctx.Response().WriteAny(result)
		}
	}
}

func writeRawJSON(ctx middleware.IMiddleContext, data RawJSON) {
	code := data.StatusCode
	if code <= 0 {
		code = http.StatusOK
	}
	for k, values := range data.Header {
		for _, v := range values {
			ctx.Response().GetHTTPReponse().Header().Add(k, v)
		}
	}
	if ctx.Response().GetHTTPReponse().Header().Get("Content-Type") == "" {
		ctx.Response().GetHTTPReponse().Header().Set("Content-Type", "application/json; charset=utf-8")
	}
	ctx.Response().GetHTTPReponse().WriteHeader(code)
	_, _ = ctx.Response().GetHTTPReponse().Write(data.Body)
	ctx.Response().NoNeedWrite(code)
}

func writeSSEProxy(ctx middleware.IMiddleContext, proxy SSEProxy) {
	w := ctx.Response().GetHTTPReponse()
	if _, ok := w.(http.Flusher); !ok {
		writeOpenAIError(ctx, http.StatusInternalServerError, fmt.Errorf("streaming unsupported"))
		return
	}
	code := proxy.StatusCode
	if code <= 0 {
		code = http.StatusOK
	}
	for k, values := range proxy.Header {
		for _, v := range values {
			w.Header().Add(k, v)
		}
	}
	w.Header().Set("Content-Type", "text/event-stream; charset=utf-8")
	w.Header().Set("Cache-Control", "no-cache")
	w.Header().Set("Connection", "keep-alive")
	w.WriteHeader(code)
	ctx.Response().NoNeedWrite(code)

	req := ctx.Request().GetHTTPRequest()
	if err := copySSEProxy(req.Context(), w, proxy, streamTimeout(ctx)); err != nil {
		ctx.Log().Warn("aigw.stream.closed", err)
	}
}

func rpcErrorStatus(response *pkgs.Rspns) int {
	if response == nil {
		return http.StatusBadGateway
	}
	status := response.GetStatus()
	if status <= 0 {
		return http.StatusBadGateway
	}
	return status
}

func streamTimeout(ctx middleware.IMiddleContext) time.Duration {
	conf, err := aigwconf.GetConf(ctx.APPConf().GetServerConf())
	if err != nil {
		return time.Duration(aigwconf.DefaultStreamTimeOut) * time.Second
	}
	return time.Duration(conf.GetStreamTimeout()) * time.Second
}

func copySSEProxy(ctx context.Context, w http.ResponseWriter, proxy SSEProxy, idleTimeout time.Duration) error {
	if proxy.Body == nil {
		return nil
	}
	flusher, _ := w.(http.Flusher)
	closer, _ := proxy.Body.(io.Closer)
	if closer != nil {
		defer closer.Close()
	}
	closeBody := func() {
		if closer != nil {
			_ = closer.Close()
		}
	}
	if ctx != nil && ctx.Done() != nil {
		go func() {
			<-ctx.Done()
			closeBody()
		}()
	}
	var timer *time.Timer
	if idleTimeout > 0 {
		timer = time.AfterFunc(idleTimeout, closeBody)
		defer timer.Stop()
	}
	buf := make([]byte, 32*1024)
	for {
		if ctx != nil {
			select {
			case <-ctx.Done():
				return ctx.Err()
			default:
			}
		}
		n, err := proxy.Body.Read(buf)
		if timer != nil {
			timer.Reset(idleTimeout)
		}
		if n > 0 {
			if _, writeErr := w.Write(buf[:n]); writeErr != nil {
				return writeErr
			}
			if flusher != nil {
				flusher.Flush()
			}
		}
		if err == io.EOF {
			return nil
		}
		if err != nil {
			if ctx != nil && ctx.Err() != nil {
				return ctx.Err()
			}
			return err
		}
	}
}

func writeOpenAIError(ctx middleware.IMiddleContext, status int, err error) {
	if status <= 0 {
		status = http.StatusBadRequest
	}
	ctx.Response().ContentType("application/json; charset=utf-8")
	ctx.Response().Write(status, map[string]interface{}{
		"error": map[string]interface{}{
			"message": err.Error(),
			"type":    "gateway_error",
			"code":    http.StatusText(status),
		},
	})
}

func fallback(ctx middleware.IMiddleContext, service string) {
	if ctx.Request().Path().AllowFallback() {
		if h, ok := services.Def.GetFallback(ctx.APPConf().GetServerConf().GetServerType(), service); ok {
			ctx.Response().WriteAny(h.Handle(ctx))
			return
		}
	}
	writeOpenAIError(ctx, http.StatusTooManyRequests, fmt.Errorf("too many requests"))
}
