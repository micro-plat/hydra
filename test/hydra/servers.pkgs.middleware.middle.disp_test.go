package hydra

import (
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"testing"

	"github.com/micro-plat/hydra/components/queues/mq/redis"
	"github.com/micro-plat/hydra/conf/server/queue"
	"github.com/micro-plat/hydra/conf/server/task"
	"github.com/micro-plat/hydra/hydra/servers/cron"
	"github.com/micro-plat/hydra/hydra/servers/mqc"
	"github.com/micro-plat/hydra/hydra/servers/pkg/dispatcher"
	"github.com/micro-plat/hydra/hydra/servers/pkg/middleware"
	"github.com/micro-plat/hydra/test/assert"
	"github.com/micro-plat/hydra/test/mocks"
)

// It accepts
//   - Full crontab specs, e.g. "* * * * * ?"
//   - Descriptors, e.g. "@midnight", "@every 1h30m"
func getTestCronTask(cronName string, service string, opts ...task.Option) *cron.CronTask {
	c, _ := cron.NewCronTask(task.NewTask(cronName, service, opts...))
	return c
}

//message 为json
func getTestMqcQueue(queueName, service, message string, hasData bool) *mqc.Request {
	c, _ := mqc.NewRequest(queue.NewQueue(queueName, service), &redis.RedisMessage{Message: message, HasData: hasData})
	return c
}

func Test_dispCtx_GetRouterPath(t *testing.T) {
	tests := []struct {
		name    string
		request dispatcher.IRequest
		want    string
	}{
		{name: "cron", request: getTestCronTask("@every 1h30m", "cron_service"), want: "cron_service"},
		{name: "mqc", request: getTestMqcQueue("queue_name", "queue_service", `{"data":"message"}`, true), want: "queue_service"},
	}
	for _, tt := range tests {
		g := middleware.NewDispCtx()
		g.Request = tt.request
		got := g.GetRouterPath()
		assert.Equal(t, tt.want, got, tt.name)
	}
}

func Test_dispCtx_GetBody(t *testing.T) {
	tests := []struct {
		name    string
		request dispatcher.IRequest
		want    string
	}{
		{name: "cron获取不到", request: getTestCronTask("@every 1h30m", "cron_service"), want: ""},
		{name: "mqc", request: getTestMqcQueue("queue_name", "queue_service", `{"data":"message"}`, true), want: `{"data": "message"}`},
	}
	for _, tt := range tests {
		g := middleware.NewDispCtx()
		g.Request = tt.request
		got := g.GetBody()
		s, err := ioutil.ReadAll(got)
		assert.Equal(t, false, err != nil, tt.name)
		fmt.Println(string(s), err)
		assert.Equal(t, tt.want, got, tt.name)
	}
}

func Test_dispCtx_GetMethod(t *testing.T) {
	tests := []struct {
		name    string
		request dispatcher.IRequest
		want    string
	}{
		{name: "cron", request: getTestCronTask("@every 1h30m", "cron_service"), want: "GET"},
		{name: "mqc", request: getTestMqcQueue("queue_name", "queue_service", `{"data":"message"}`, true), want: `GET`},
	}
	for _, tt := range tests {
		g := middleware.NewDispCtx()
		g.Request = tt.request
		got := g.GetMethod()
		assert.Equal(t, tt.want, got, tt.name)
	}
}

func Test_dispCtx_GetURL(t *testing.T) {
	tests := []struct {
		name    string
		request dispatcher.IRequest
		want    *url.URL
	}{
		{name: "cron", request: getTestCronTask("@every 1h30m", "/cron/service"), want: &url.URL{Path: "/cron/service"}},
		{name: "mqc", request: getTestMqcQueue("queue_name", "/mqc/service", `{"data":"message"}`, true), want: &url.URL{Path: "/mqc/service"}},
	}
	for _, tt := range tests {
		g := middleware.NewDispCtx()
		g.Request = tt.request
		got := g.GetURL()
		assert.Equal(t, tt.want, got, tt.name)
	}
}
func Test_dispCtx_GetHeaders(t *testing.T) {
	tests := []struct {
		name    string
		request dispatcher.IRequest
		want    http.Header
	}{
		{name: "cron", request: getTestCronTask("@every 1h30m", "/cron/service"), want: http.Header{
			"Client-IP": []string{"127.0.0.1"},
		}},
		{name: "mqc", request: getTestMqcQueue("queue_name", "/mqc/service", `{"data":"message"}`, true), want: http.Header{
			"Client-IP": []string{"127.0.0.1"},
		}},
	}
	for _, tt := range tests {
		g := middleware.NewDispCtx()
		g.Request = tt.request
		got := g.GetHeaders()
		assert.Equal(t, tt.want, got, tt.name)
	}
}

func Test_dispCtx_GetQuery(t *testing.T) {
	tests := []struct {
		name    string
		request dispatcher.IRequest
		k       string
		want    string
		wantOK  bool
	}{
		{name: "cron获取不到", request: getTestCronTask("@every 1h30m", "cron_service")},
		{name: "mqc获取错误的key", request: getTestMqcQueue("queue_name", "queue_service", `{"data":"message"}`, true)},
		{name: "mqc获取正确的", request: getTestMqcQueue("queue_name", "queue_service", `{"data":"message"}`, true), want: `{"data":"message"}`, k: "__body_", wantOK: true},
	}
	for _, tt := range tests {
		g := middleware.NewDispCtx()
		g.Request = tt.request
		got, gotOK := g.GetQuery(tt.k)
		assert.Equal(t, tt.wantOK, gotOK, tt.name)
		if gotOK {
			assert.Equal(t, tt.want, got, tt.name)
		}
	}
}

func Test_dispCtx_GetFormValue(t *testing.T) {
	tests := []struct {
		name    string
		request dispatcher.IRequest
		k       string
		want    string
		wantOK  bool
	}{
		{name: "cron获取不到", request: getTestCronTask("@every 1h30m", "cron_service")},
		{name: "mqc获取错误的key", request: getTestMqcQueue("queue_name", "queue_service", `{"data":"message"}`, true)},
		{name: "mqc获取正确的", request: getTestMqcQueue("queue_name", "queue_service", `{"data":"message"}`, true), want: `{"data":"message"}`, k: "__body_", wantOK: true},
	}
	for _, tt := range tests {
		g := middleware.NewDispCtx()
		g.Request = tt.request
		got, gotOK := g.GetFormValue(tt.k)
		assert.Equal(t, tt.wantOK, gotOK, tt.name)
		if gotOK {
			assert.Equal(t, tt.want, got, tt.name)
		}
	}
}

func Test_dispCtx_GetForm(t *testing.T) {
	tests := []struct {
		name    string
		request dispatcher.IRequest
		want    url.Values
	}{
		{name: "cron", request: getTestCronTask("@every 1h30m", "cron_service"), want: url.Values{}},
		{name: "mqc", request: getTestMqcQueue("queue_name", "queue_service", `{"data":"message"}`, true), want: url.Values{
			"__body_": []string{`{"data":"message"}`},
			"data":    []string{"message"},
		}},
	}
	for _, tt := range tests {
		g := middleware.NewDispCtx()
		g.Request = tt.request
		got := g.GetForm()
		fmt.Println(got)
		assert.Equal(t, tt.want, got, tt.name)
	}
}

func Test_dispCtx_Status(t *testing.T) {
	tests := []struct {
		name   string
		status int
	}{
		{name: "设置200", status: 200},
		{name: "设置302", status: 302},
	}
	for _, tt := range tests {
		g := middleware.NewDispCtx()
		g.Writer = &mocks.MockResponseWriter2{}
		g.WStatus(tt.status)
		got := g.Status()
		assert.Equal(t, tt.status, got, tt.name)
	}
}

func Test_dispCtx_File(t *testing.T) {
	//构建
	g := middleware.NewDispCtx()
	writer := &mocks.MockResponseWriter2{}
	writer.Reset()
	g.Writer = writer
	//判断未写入状态
	w := g.Written()
	assert.Equal(t, false, w, "是否写入数据判断")
	fileName := "servers.pkgs.middleware.middle.test.txt"

	//写入文件
	g.File(fileName)

	//判断Header
	f := g.WHeader("file")
	assert.Equal(t, fileName, f, "写入数据Header判断")

	//判断content-type
	c := g.WHeader("Content-Type")
	assert.Equal(t, "application/json; charset=utf-8", c, "写入数据Header ctp判断")

	//判断写入状态
	w = g.Written()
	assert.Equal(t, true, w, "写入数据判断")
	s := g.Status()
	assert.Equal(t, 200, s, "写入数据状态判断")

	//写入数据判断
	content, _ := ioutil.ReadFile(fileName)
	jsonBytes, _ := json.Marshal(map[string]string{"__body_": base64.StdEncoding.EncodeToString(content)})
	assert.Equal(t, len(jsonBytes), g.Writer.Size(), "写入数据长度判断")
	assert.Equal(t, jsonBytes, g.Writer.Data(), "写入数据判断")

}

func Test_dispCtx_ShouldBind(t *testing.T) {
	tests := []struct {
		name    string
		request dispatcher.IRequest
		bind    interface{}
		wantErr string
		want    interface{}
	}{
		{name: "GetForm为空", request: getTestCronTask("@every 1h30m", "cron_service"), bind: map[string]interface{}{}, want: map[string]interface{}{}},
		{name: "mqc,__body_为空", request: getTestMqcQueue("queue_name", "queue_service", `{}`, true), bind: map[string]interface{}{}, want: map[string]interface{}{}},
		{name: "mqc,__body_不为空", request: getTestMqcQueue("queue_name", "queue_service", `{"data":"message"}`, true),
			bind: map[string]interface{}{}, want: map[string]interface{}{
				"__body_": `{"data":"message"}`,
				"data":    "message",
			}},
	}
	for _, tt := range tests {
		g := middleware.NewDispCtx()
		g.Request = tt.request
		gotErr := g.ShouldBind(&tt.bind)
		if tt.wantErr != "" {
			assert.Equal(t, tt.wantErr, gotErr.Error(), tt.name)
			continue
		}
		assert.Equal(t, tt.want, tt.bind, tt.name)
	}
}
