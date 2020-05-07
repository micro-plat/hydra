package middleware

//RASAuth 远程服务验证
func RASAuth() Handler {
	return func(ctx IMiddleContext) {

		//获取FSA配置
		auth := ctx.ServerConf().GetRASConf()
		if len(auth) == 0 {
			ctx.Next()
			return
		}

		//检查必须参数

		// b, auth := auth.Contains(ctx.Request().Path().GetService())
		// if !b {
		// 	ctx.Next()
		// 	return
		// }

		ctx.Response().AddSpecial("ras")

		// header := ctx.Request().Path().GetHeaders()
		// cookie := ctx.Request().Path().GetCookies()
		// header["method"] = ctx.Request().Path().GetMethod()

		// input := types.NewXMapByMap(ctx.Request().GetData())

		// input["__auth_"], err = auth.AuthString()
		// if err != nil {
		// 	return fmt.Errorf("将service.auth转换为__auth_失败:%v", err)
		// }

		// status, result, _, err := ctx.RPC.Request(auth.Service, header, input.ToMap(), true)
		// if err != nil || status != 200 {
		// 	return context.NewErrorf(types.GetMax(status, 403), "远程认证失败:%s,err:%v(%d)", err, result, status)
		// }
		// tmp := types.XMap{}
		// if err := json.Unmarshal([]byte(result), &tmp); err != nil {
		// 	return err
		// }
		// ctx.Request.Metadata.SetStrings(tmp.ToSMap())
		return
	}
}
