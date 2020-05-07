package application

import "github.com/micro-plat/hydra/services"



//Micro 注册为微服务包括api,web,rpc
func Micro(name string, h interface{}) {
	services.Registry.Micro(name, h)
}

//Flow 注册为流程服务，包括mqc,cron
func Flow(name string, h interface{}) {
	services.Registry.Flow(name, h)
}

//API 注册为API服务
func API(name string, h interface{}) {
	services.Registry.API(name, h)
}

//Web 注册为web服务
func Web(name string, h interface{}) {
	services.Registry.Web(name, h)
}

//WS 注册为websocket服务
func WS(name string, h interface{}) {
	services.Registry.WS(name, h)
}

//MQC 注册为消息队列服务
func MQC(name string, h interface{}) {
	services.Registry.MQC(name, h)
}

//CRON 注册为定时任务服务
func CRON(name string, h interface{}) {
	services.Registry.CRON(name, h)
}
