package main

import "github.com/micro-plat/hydra/context"

var funcAPI1 func(ctx context.IContext) (r interface{}) = func(ctx context.IContext) (r interface{}) {
	ctx.Log().Info("api-all 接口服务测试")
	res, err := ctx.Request().GetBody()
	if err != nil {
		ctx.Log().Errorf("ctx.Request().err:%v", err)
		return err
	}
	ctx.Log().Info("ctx.Request().GetBody():", res)
	resmap, err := ctx.Request().GetMap()
	if err != nil {
		ctx.Log().Errorf("ctx.Request()GetMap.err:%v", err)
		return err
	}
	ctx.Log().Info("ctx.Request().GetMap():", resmap)
	return "success"
}

var funcAPI2 func(ctx context.IContext) (r interface{}) = func(ctx context.IContext) (r interface{}) {
	ctx.Log().Info("api-all 接口服务测试")

	ctx.Log().Info("ctx.Request().Param(xxx):", ctx.Request().Param("xxx"))
	res, err := ctx.Request().GetBody()
	if err != nil {
		ctx.Log().Errorf("ctx.Request().err:%v", err)
		return err
	}
	ctx.Log().Info("ctx.Request().GetBody():", res)
	resmap, err := ctx.Request().GetMap()
	if err != nil {
		ctx.Log().Errorf("ctx.Request()GetMap.err:%v", err)
		return err
	}
	ctx.Log().Info("ctx.Request().GetMap():", resmap)
	return "success"
}

type apiGet struct{}

//如果该方法改名为GetHandle   则无法正常访问
func (s *apiGet) Handle(ctx context.IContext) (r interface{}) {
	ctx.Log().Info("api-get 接口服务测试")
	res, err := ctx.Request().GetBody()
	if err != nil {
		ctx.Log().Errorf("ctx.Request().err:%v", err)
		return err
	}
	ctx.Log().Info("ctx.Request().GetBody():", res)
	resmap, err := ctx.Request().GetMap()
	if err != nil {
		ctx.Log().Errorf("ctx.Request()GetMap.err:%v", err)
		return err
	}
	ctx.Log().Info("ctx.Request().GetMap():", resmap)
	return "api-get-success"
}

func (s *apiGet) GetHandle(ctx context.IContext) (r interface{}) {
	ctx.Log().Info("api-get 接口服务测试1")
	res, err := ctx.Request().GetBody()
	if err != nil {
		ctx.Log().Errorf("ctx.Request().err:%v", err)
		return err
	}
	ctx.Log().Info("ctx.Request().GetBody():", res)
	resmap, err := ctx.Request().GetMap()
	if err != nil {
		ctx.Log().Errorf("ctx.Request()GetMap.err:%v", err)
		return err
	}
	ctx.Log().Info("ctx.Request().GetMap():", resmap)
	return "api-get-success1"
}

func (s *apiGet) QueryHandle(ctx context.IContext) (r interface{}) {
	ctx.Log().Info("api-get 接口服务测试2")
	res, err := ctx.Request().GetBody()
	if err != nil {
		ctx.Log().Errorf("ctx.Request().err:%v", err)
		return err
	}
	ctx.Log().Info("ctx.Request().GetBody():", res)
	resmap, err := ctx.Request().GetMap()
	if err != nil {
		ctx.Log().Errorf("ctx.Request()GetMap.err:%v", err)
		return err
	}
	ctx.Log().Info("ctx.Request().GetMap():", resmap)
	return "api-get-success2"
}

type apiGetgbk struct{}

func (s *apiGetgbk) Handle(ctx context.IContext) (r interface{}) {
	ctx.Log().Info("api-GBK 接口服务测试")
	res, err := ctx.Request().GetBody()
	if err != nil {
		ctx.Log().Errorf("ctx.Request().err:%v", err)
		return err
	}
	ctx.Log().Info("ctx.Request().GetBody():", res)
	resmap, err := ctx.Request().GetMap()
	if err != nil {
		ctx.Log().Errorf("ctx.Request()GetMap.err:%v", err)
		return err
	}
	ctx.Log().Info("ctx.Request().GetMap():", resmap)
	ctx.Log().Info("ctx.Request().GetHeaders():", ctx.Request().GetHeaders())
	return "api-GBK-success"
}

type apiGetgb2312 struct{}

func (s *apiGetgb2312) Handle(ctx context.IContext) (r interface{}) {
	ctx.Log().Info("api-GB2312 接口服务测试")
	res, err := ctx.Request().GetBody()
	if err != nil {
		ctx.Log().Errorf("ctx.Request().err:%v", err)
		return err
	}
	ctx.Log().Info("ctx.Request().GetBody():", res)
	resmap, err := ctx.Request().GetMap()
	if err != nil {
		ctx.Log().Errorf("ctx.Request()GetMap.err:%v", err)
		return err
	}
	ctx.Log().Info("ctx.Request().GetMap():", resmap)
	ctx.Log().Info("ctx.Request().GetHeaders():", ctx.Request().GetHeaders())
	return "api-GB2312-success"
}
