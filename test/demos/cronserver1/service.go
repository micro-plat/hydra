package main

import (
	"fmt"

	"github.com/micro-plat/hydra/context"
)

//----------
type mqcService struct {
}

func (s *mqcService) Handle(ctx context.IContext) interface{} {
	return "object.success"
}

func FuncService(ctx context.IContext) interface{} {
	return "func.success"
}

func NewObjNoneError() (*mqcService, error) {
	return &mqcService{}, nil
}

func NewObjWithError() *mqcService {
	return &mqcService{}
}

type cronService struct{}

func (s *cronService) Handle(ctx context.IContext) interface{} {
	fmt.Println("--cron")
	//	queueObj := hydra.C.Queue().GetRegularQueue("mqcqueue")
	//	return queueObj.Push("mqcservice:proc1", fmt.Sprintf(`{"xxx":"%d"}`, time.Now().Unix()))
	return nil
}
