package middleware

//RASAuth RAS远程认证
func RASAuth() Handler {
	return func(ctx IMiddleContext) {

		//获取FSA配置
		auth := ctx.ServerConf().GetRASConf()
		if len(auth) == 0 {
			ctx.Next()
			return
		}

		//检查必须参数
		// b, auth := auth.Contains(ctx.Request().Path().GetRouter().Path)
		// if !b {
		// 	ctx.Next()
		// 	return
		// }

		// ctx.Response().AddSpecial("ras")

		// input, err := ctx.Request().GetData()
		// if err != nil {
		// 	ctx.Response().AbortWithError(500, err)
		// 	return
		// }

		// input["__auth_"], err = auth.AuthString()
		// if err != nil {
		// 	return fmt.Errorf("将service.auth转换为__auth_失败:%v", err)
		// }

		// respones, err := components.Def.RPC().GetRegularRPC().Request(ctx.Context(), auth.Service, input)
		// if err != nil || !respones.Success() {
		// 	return context.NewErrorf(types.GetMax(respones.Status, 403), "远程认证失败:%s,err:%v(%d)", err, result, respones.Status)
		// }
		// tmp := types.XMap{}
		// if err := json.Unmarshal([]byte(result), &tmp); err != nil {
		// 	return err
		// }

		// ctx.Request.Metadata.SetStrings(tmp.ToSMap())
		return
	}
}
