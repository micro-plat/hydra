package main

import (
	"fmt"

	"github.com/micro-plat/hydra/context"
)

func (api *apiserver) handling() {
	api.MicroApp.Handling(func(ctx *context.Context) (rt interface{}) {
		if err := ctx.Request.Check("merchant_id"); err != nil {
			return err
		}
		key, err := merchant.GetKey(ctx.Request.GetInt(merchant_id))
		if err != nil {
			return err
		}
		if !ctx.Request.CheckSign(key) {
			return fmt.Errorf(908, "商户签名错误")
		}
		return nil
	})
}
