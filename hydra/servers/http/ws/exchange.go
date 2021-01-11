package ws

import (
	"fmt"
	"net/http"
	"sync"

	"github.com/micro-plat/hydra"
	"github.com/micro-plat/hydra/context"
	"github.com/micro-plat/lib4go/concurrent/cmap"
	"github.com/micro-plat/lib4go/errs"
)

//WSExchange web socket message exchange
var exchange = NewExchange()

//DataExchange 数据交换处理
var DataExchange IDataExchange = exchange

var confName = "queue"

//Conf 配置管理
func Conf(queueConfName string) {
	confName = queueConfName
}

//IDataExchange 数据交换接口
type IDataExchange interface {
	Notify(uuid string, data interface{}) error
}

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
		queueFormatName: "ws:exchange:%s",
		service:         "/ws/handle",
	}
}

//Subscribe 订阅消息通知
func (e *Exchange) Subscribe(uuid string, f func(...interface{}) error) error {
	e.once.Do(func() {
		hydra.S.MQC(e.service, e.handle) //注册MQC服务
	})

	if ok, _ := e.uuid.SetIfAbsent(uuid, f); ok {
		hydra.MQC.Add(e.getQueueName(uuid), e.service) //为每个用户添加处理队列
	}
	return nil
}

//Unsubscribe 取消订阅
func (e *Exchange) Unsubscribe(uuid string) {
	e.uuid.Remove(uuid)
	hydra.MQC.Remove(e.getQueueName(uuid), e.service) //关闭队列
}

//Notify 发送通知消息
func (e *Exchange) Notify(uuid string, msg interface{}) error {
	queue, err := hydra.C.Queue().GetQueue(confName)
	if err != nil {
		return err
	}
	return queue.Send(e.getQueueName(uuid), msg, uuid)
}

//handle 业务回调处理
func (e *Exchange) handle(ctx context.IContext) interface{} {
	uuid := ctx.User().GetTraceID()
	v, ok := e.uuid.Get(uuid)
	if !ok {
		return errs.NewError(http.StatusNoContent, nil)
	}
	callback := v.(func(...interface{}) error)
	body, err := ctx.Request().GetBody()
	if err != nil {
		return err
	}
	return callback(string(body))
}

func (e *Exchange) getQueueName(id string) string {
	return fmt.Sprintf(e.queueFormatName, id)
}
