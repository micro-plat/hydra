package middleware

import (
	"github.com/micro-plat/hydra/servers/pkg/circuit"

	"github.com/gin-gonic/gin"
	"github.com/micro-plat/hydra/conf"
)

//CircuitBreak 熔断处理
func CircuitBreak(conf *conf.MetadataConf) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		circuitBreaker, ok := conf.GetMetadata("__circuit-breaker_").(*circuit.NamedCircuitBreakers)
		if !ok {
			ctx.Next()
			return
		}
		url := ctx.Request.URL.Path
		breaker := circuitBreaker.GetBreaker(url)
		isOpen, allowRequest := breaker.GetCircuitStatus()
		setIsCircuitBreaker(ctx, isOpen)
		if !allowRequest {
			breaker.ReportEvent(circuit.EventReject, 1)
			ctx.AbortWithStatus(503)
			return
		}
		ctx.Next()
		nctx := getCTX(ctx)
		if nctx == nil {
			return
		}
		success := nctx.Response.GetStatus() < 400
		if !isOpen {
			if success {
				breaker.ReportEvent(circuit.EventSuccess, 1)
				return
			}
			breaker.ReportEvent(circuit.EventFailure, 1)
			return
		}
		setExt(ctx, "fb")
		if success {
			breaker.ReportEvent(circuit.EventFallbackSuccess, 1)
			return
		}
		breaker.ReportEvent(circuit.EventFallbackFailure, 1)
		return
	}

}
