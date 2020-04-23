package engines

import (
	"encoding/json"
	"fmt"
	"strings"

	"github.com/micro-plat/hydra/conf"
	"github.com/micro-plat/hydra/context"
	"github.com/micro-plat/lib4go/types"
)

//checkSignByFixedSecret 根据固定secret检查签名
func checkSignByFixedSecret(ctx *context.Context) error {
	fsConf, err := ctx.Request.GetFixedSecretConfig()
	if err == conf.ErrNoSetting || !ctx.IsMicroServer() {
		return nil
	}
	if !fsConf.Contains(ctx.Service) {
		return nil
	}
	ctx.Response.SetHeader("__auth_tag_", "FAUTH")
	if err := ctx.Request.Check("sign", "timestamp"); err != nil {
		return context.NewError(context.ERR_NOT_ACCEPTABLE, err)
	}
	_, err = ctx.Request.CheckSign(fsConf.Secret, fsConf.Mode)
	if err == nil {
		return nil
	}
	return context.NewErrorf(401, "签名认证失败%v", err)
}

//checkSignByRemoteSecret 根据固定secret检查签名
func checkSignByRemoteSecret(ctx *context.Context) error {
	fsConf, err := ctx.Request.GetRemoteAuthConfig()
	if err == conf.ErrNoSetting || !ctx.IsMicroServer() {
		return nil
	}
	b, auth := fsConf.Contains(ctx.Service)
	if !b {
		return nil
	}
	header, _ := ctx.Request.Http.GetHeader()
	cookie, _ := ctx.Request.Http.GetCookies()
	for k, v := range cookie {
		header[k] = v
	}
	header["method"] = strings.ToUpper(ctx.Request.GetMethod())
	input := types.NewXMapByMap(ctx.Request.GetRequestMap())
	iparam := types.XMap(auth.Params)
	input.Merge(iparam)
	input["__auth_"], err = auth.AuthString()
	if err != nil {
		return fmt.Errorf("将service.auth转换为__auth_失败:%v", err)
	}
	ctx.Response.SetHeader("__auth_tag_", "RAUTH")
	status, result, _, err := ctx.RPC.Request(auth.Service, header, input.ToMap(), true)
	if err != nil || status != 200 {
		return context.NewErrorf(types.GetMax(status, 403), "远程认证失败:%s,err:%v(%d)", err, result, status)
	}
	tmp := types.XMap{}
	if err := json.Unmarshal([]byte(result), &tmp); err != nil {
		return err
	}
	ctx.Request.Metadata.SetStrings(tmp.ToSMap())
	return nil

}
