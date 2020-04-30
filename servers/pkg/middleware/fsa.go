package middleware



//FixedSecretAuth 静态密钥验证
func FixedSecretAuth() Handler {
	return func(ctx IMiddleContext) {


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
		ctx.Next()
	}
}