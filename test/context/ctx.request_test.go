package context

import (
	"net/url"
	"reflect"
	"testing"

	"github.com/micro-plat/hydra/conf"
	"github.com/micro-plat/hydra/context"
	"github.com/micro-plat/hydra/context/ctx"
	"github.com/micro-plat/hydra/test/mocks"
)

func Test_request_Bind(t *testing.T) {
	confObj := mocks.NewConf()         //构建对象
	confObj.API(":8080")               //初始化参数
	serverConf := confObj.GetAPIConf() //获取配置

	t.Run("非指针传递,无法进行数据绑定", func(t *testing.T) {
		out := 1
		r := ctx.NewRequest(&mocks.TestContxt{
			Form: url.Values{"__body_": []string{`9`}},
		}, serverConf, conf.NewMeta())
		if err := r.Bind(out); (err != nil) != false {
			t.Errorf("request.Bind() error = %v, wantErr %v", err, false)
		}
		if !reflect.DeepEqual(out, 1) {
			t.Errorf("request.Bind() out= %v, want %v", out, 1)
		}
	})

	t.Run("参数类型与数据的类型不一致", func(t *testing.T) {
		out := 1
		r := ctx.NewRequest(&mocks.TestContxt{
			Form: url.Values{"__body_": []string{`{"key":"value"}`}},
		}, serverConf, conf.NewMeta())
		if err := r.Bind(&out); (err != nil) != true {
			t.Errorf("request.Bind() error = %v, wantErr %v", err, false)
		}
	})

	t.Run("参数数据绑定", func(t *testing.T) {
		out := 1
		r := ctx.NewRequest(&mocks.TestContxt{
			Form: url.Values{"__body_": []string{`9`}},
		}, serverConf, conf.NewMeta())
		if err := r.Bind(&out); (err != nil) != false {
			t.Errorf("request.Bind() error = %v, wantErr %v", err, false)
		}
		if !reflect.DeepEqual(out, 9) {
			t.Errorf("request.Bind() out= %v, want %v", out, 9)
		}
	})

	t.Run("参数类型为interface{}", func(t *testing.T) {
		var out interface{}
		out = 1
		r := ctx.NewRequest(&mocks.TestContxt{
			Form: url.Values{"__body_": []string{`9`}},
		}, serverConf, conf.NewMeta())
		if err := r.Bind(&out); (err != nil) != false {
			t.Errorf("request.Bind() error = %v, wantErr %v", err, false)
		}
		var want float64
		want = 9
		if !reflect.DeepEqual(out, want) {
			t.Errorf("request.Bind() out= %v, want %v", out, 9)
		}
	})

	type result struct {
		Key   string `json:"key" valid:"required"`
		Value string `json:"value" valid:"required"`
	}

	t.Run("参数类型为结构体,绑定数据非空验证", func(t *testing.T) {
		out := result{}
		r := ctx.NewRequest(&mocks.TestContxt{
			Form: url.Values{"__body_": []string{`{"key":"","value":"2"}`}},
		}, serverConf, conf.NewMeta())
		if err := r.Bind(&out); err.Error() != "输入参数有误 key: non zero value required" {
			t.Errorf("request.Bind() error = %v, wantErr 输入参数有误 key: non zero value required", err)
		}
	})
}

func Test_request_Check(t *testing.T) {
	confObj := mocks.NewConf()         //构建对象
	confObj.API(":8080")               //初始化参数
	serverConf := confObj.GetAPIConf() //获取配置

	type fields struct {
		ctx context.IInnerContext
	}
	type args struct {
		field []string
	}
	tests := []struct {
		name       string
		fields     fields
		args       args
		wantErr    bool
		wantErrStr string
	}{
		{name: "验证非空数据判断", fields: fields{ctx: &mocks.TestContxt{
			Body:            `{"key1":"value1","key2":"value2"}`,
			HttpContentType: context.JSONF,
		}}, args: args{field: []string{"key1", "key2"}}, wantErr: false},
		{name: "验证空数据判断", fields: fields{ctx: &mocks.TestContxt{
			Body:            `{"key1":"","key2":"value2"}`,
			HttpContentType: context.JSONF,
		}}, args: args{field: []string{"key1", "key2"}}, wantErrStr: "输入参数:key1值不能为空", wantErr: true},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			r := ctx.NewRequest(tt.fields.ctx, serverConf, conf.NewMeta())
			if err := r.Check(tt.args.field...); err != nil {
				if (err != nil) != tt.wantErr {
					t.Errorf("request.Check() error = %v, wantErr %v", err, tt.wantErr)
				}
				if err.Error() != tt.wantErrStr {
					t.Errorf("request.Check() error = %v, wantErrStr %v", err, tt.wantErrStr)
				}
			}
		})
	}
}

func Test_request_GetKeys(t *testing.T) {
	confObj := mocks.NewConf()         //构建对象
	confObj.API(":8080")               //初始化参数
	serverConf := confObj.GetAPIConf() //获取配置

	type fields struct {
		ctx context.IInnerContext
	}
	tests := []struct {
		name   string
		fields fields
		want   []string
	}{
		{name: "获取所有key", fields: fields{ctx: &mocks.TestContxt{
			Form:            url.Values{"key3": []string{}},
			Body:            `{"key1":"value1","key2":"value2"}`,
			HttpContentType: context.JSONF,
		}}, want: []string{"key3", "key1", "key2"}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			r := ctx.NewRequest(tt.fields.ctx, serverConf, conf.NewMeta())
			if got := r.GetKeys(); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("request.GetKeys() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_request_GetMap(t *testing.T) {
	confObj := mocks.NewConf()         //构建对象
	confObj.API(":8080")               //初始化参数
	serverConf := confObj.GetAPIConf() //获取配置
	type fields struct {
		ctx context.IInnerContext
	}
	tests := []struct {
		name    string
		fields  fields
		want    map[string]interface{}
		wantErr bool
	}{
		{name: "获取所有map", fields: fields{ctx: &mocks.TestContxt{
			Form:            url.Values{"key3": []string{"value3"}},
			Body:            `{"key1":"value1","key2":"value2"}`,
			HttpContentType: context.JSONF,
		}}, want: map[string]interface{}{"key1": "value1", "key2": "value2", "key3": "value3"}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			r := ctx.NewRequest(tt.fields.ctx, serverConf, conf.NewMeta())
			got, err := r.GetMap()
			if (err != nil) != tt.wantErr {
				t.Errorf("request.GetMap() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("request.GetMap() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_request_Get(t *testing.T) {
	confObj := mocks.NewConf()         //构建对象
	confObj.API(":8080")               //初始化参数
	serverConf := confObj.GetAPIConf() //获取配置

	r := ctx.NewRequest(&mocks.TestContxt{
		Body:            `{"key1":"value1"}`,
		Form:            url.Values{"key2": []string{"%20+value2"}},
		HttpContentType: context.JSONF,
	}, serverConf, conf.NewMeta())

	type args struct {
		name string
	}
	tests := []struct {
		name       string
		args       args
		wantResult string
		wantOk     bool
	}{
		{name: "通过BodyMap获取key对应的value", args: args{name: "key1"}, wantResult: "value1", wantOk: true},
		{name: "通过FormValue获取key对应的value", args: args{name: "key2"}, wantResult: "  value2", wantOk: true},
		{name: "获取不存在key的值", args: args{name: "key3"}, wantResult: "", wantOk: false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotResult, gotOk := r.Get(tt.args.name)
			if gotResult != tt.wantResult {
				t.Errorf("request.Get() gotResult = %v, want %v", gotResult, tt.wantResult)
			}
			if gotOk != tt.wantOk {
				t.Errorf("request.Get() gotOk = %v, want %v", gotOk, tt.wantOk)
			}
		})
	}
}
