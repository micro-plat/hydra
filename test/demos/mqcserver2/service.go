package main

import "github.com/micro-plat/hydra/context"

//----------
type objService struct {
}

func (s *objService) Handle(ctx context.IContext) interface{} {
	return "object.success"
}

func FuncService(ctx context.IContext) interface{} {
	return "func.success"
}

func NewObjNoneError() (*objService, error) {
	return &objService{}, nil
}

func NewObjWithError() *objService {
	return &objService{}
}
