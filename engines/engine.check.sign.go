package engines

import (
	"strings"
	"time"

	"github.com/micro-plat/hydra/conf"
	"github.com/micro-plat/hydra/context"
)

//checkSignByFixedSecret 根据固定secret检查签名
func checkSignByFixedSecret(ctx *context.Context) error {
	fsConf, err := ctx.Request.GetFixedSecretConfig()
	if err == conf.ErrNoSetting {
		return nil
	}
	if err := ctx.Request.Check("sign", "timestamp"); err != nil {
		return context.NewError(context.ERR_NOT_ACCEPTABLE, err)
	}
	_, err = ctx.Request.CheckSign(fsConf.Secret, fsConf.Mode)
	if err == nil {
		return nil
	}
	return context.NewErrorf(403, "签名认证失败%v", err)
}

//checkSignByRemoteSecret 根据固定secret检查签名
func checkSignByRemoteSecret(ctx *context.Context) error {
	fsConf, err := ctx.Request.GetRemoteAuthConfig()
	if err == conf.ErrNoSetting {
		return nil
	}

	header, _ := ctx.Request.Http.GetHeader()
	cookie, _ := ctx.Request.Http.GetCookies()
	for k, v := range cookie {
		header[k] = v
	}
	header["method"] = strings.ToUpper(ctx.Request.GetMethod())
	input := ctx.Request.GetRequestMap()
	timeout := ctx.Request.Setting.GetInt("timeout", 3)
	response := ctx.RPC.AsyncRequest(fsConf.RPCServiceName, header, input, true)
	status, result, params, err := response.Wait(time.Second * time.Duration(timeout))
	if err != nil {
		return context.NewErrorf(403, "调用远程认证服务失败 %v(%d)", err, status)
	}
	if status == 200 {
		return nil
	}
	ctx.Request.Metadata.SetStrings(params)
	return context.NewErrorf(status, "远程认证失败(%d)%s", status, result)

}
