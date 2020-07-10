package ws

import (
	"fmt"
	"net/http"
	"sync"

	"github.com/micro-plat/hydra"
	"github.com/micro-plat/hydra/context"
	"github.com/micro-plat/hydra/global"
	"github.com/micro-plat/lib4go/concurrent/cmap"
	"github.com/micro-plat/lib4go/errs"
)

//WSExchange web socket message exchange
var WSExchange = NewExchange()

//Exchange 数据交换中心
type Exchange struct {
	uuid            cmap.ConcurrentMap
	queueFormatName string
	service         string
	once            sync.Once
}

//NewExchange 构建数据交换中心
func NewExchange() *Exchange {
	return &Exchange{
		uuid:            cmap.New(8),
		queueFormatName: "%s:ws:%s",
		service:         "/ws/handle",
	}
}

//Subscribe 订阅消息通知
func (e *Exchange) Subscribe(uuid string, f func(...interface{}) error) error {
	e.once.Do(func() {
		hydra.S.MQC(e.service, e.handle) //注册全局处理函数
	})

	if ok, _ := e.uuid.SetIfAbsent(uuid, f); ok {
		hydra.MQC.Add(e.getQueueName(uuid), e.service)
	}
	return nil
}

//Unsubscribe 取消订阅
func (e *Exchange) Unsubscribe(uuid string) {
	e.uuid.Remove(uuid)
	hydra.MQC.Remove(e.getQueueName(uuid), e.service)
}

//Notify 消息通知
func (e *Exchange) Notify(name string, msg string) error {
	hydra.C.Queue().GetRegularQueue().Push(e.getQueueName(name), msg)
	return nil
}

//handle 业务回调处理
func (e *Exchange) handle(ctx context.IContext) interface{} {
	uuid := ctx.User().GetRequestID()
	v, ok := e.uuid.Get(uuid)
	if !ok {
		return errs.NewError(http.StatusNoContent, nil)
	}
	callback := v.(func(...interface{}) error)
	body, err := ctx.Request().GetBody()
	if err != nil {
		return err
	}
	return callback(body)

}
func (e *Exchange) getQueueName(id string) string {
	return fmt.Sprintf(e.queueFormatName, global.Def.PlatName, id)
}
