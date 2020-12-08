package main

import (
	"fmt"

	"github.com/micro-plat/hydra"
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
	obj := hydra.C.Queue().GetRegularQueue("mqcqueue")
	err := obj.Send("service:queue1", `{}`)
	fmt.Println("Send1:", err)
	err = obj.Send("service:queue6", `{}`)
	fmt.Println("Send2:", err)
	//	queueObj := hydra.C.Queue().GetRegularQueue("mqcqueue")
	//	return queueObj.Push("mqcservice:proc1", fmt.Sprintf(`{"xxx":"%d"}`, time.Now().Unix()))
	return nil
}
