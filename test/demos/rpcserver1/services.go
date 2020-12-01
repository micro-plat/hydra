package main

import "github.com/micro-plat/hydra/context"

var rpcFunc func(ctx context.IContext) (r interface{}) = func(ctx context.IContext) (r interface{}) {
	ctx.Log().Info("RPC func服务测试")
	res, err := ctx.Request().GetMap()
	if err != nil {
		ctx.Log().Errorf("ctx.GetMap().err:%v", err)
		return err
	}
	return res
}

type rpcStruct struct{}

func (s *rpcStruct) Handle(ctx context.IContext) (r interface{}) {
	ctx.Log().Info("RPC struct_handle服务测试")
	res, err := ctx.Request().GetMap()
	if err != nil {
		ctx.Log().Errorf("ctx.GetMap().err:%v", err)
		return err
	}
	return res
}

func (s *rpcStruct) GetHandle(ctx context.IContext) (r interface{}) {
	ctx.Log().Info("RPC struct_get服务测试")
	res, err := ctx.Request().GetMap()
	if err != nil {
		ctx.Log().Errorf("ctx.GetMap().err:%v", err)
		return err
	}
	return res
}

func (s *rpcStruct) QueryHandle(ctx context.IContext) (r interface{}) {
	ctx.Log().Info("RPC struct_post服务测试")
	res, err := ctx.Request().GetMap()
	if err != nil {
		ctx.Log().Errorf("ctx.GetMap().err:%v", err)
		return err
	}
	return res
}

type rpcStruct2 struct{}

func (s rpcStruct2) Handle(ctx context.IContext) (r interface{}) {
	ctx.Log().Info("RPC服务测试")
	res, err := ctx.Request().GetMap()
	if err != nil {
		ctx.Log().Errorf("ctx.GetMap().err:%v", err)
		return err
	}
	return res
}
