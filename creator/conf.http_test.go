package creator

import (
	"reflect"
	"testing"

	"github.com/micro-plat/hydra/conf/server/acl/blacklist"
	"github.com/micro-plat/hydra/conf/server/acl/whitelist"
	"github.com/micro-plat/hydra/conf/server/api"
	"github.com/micro-plat/hydra/conf/server/auth/apikey"
	"github.com/micro-plat/hydra/conf/server/auth/basic"
	"github.com/micro-plat/hydra/conf/server/auth/jwt"
	"github.com/micro-plat/hydra/conf/server/auth/ras"
	"github.com/micro-plat/hydra/conf/server/gray"
	"github.com/micro-plat/hydra/conf/server/header"
	"github.com/micro-plat/hydra/conf/server/limiter"
	"github.com/micro-plat/hydra/conf/server/metric"
	"github.com/micro-plat/hydra/conf/server/render"
	"github.com/micro-plat/hydra/conf/server/static"
	"github.com/micro-plat/hydra/services"
)

func Test_newHTTP(t *testing.T) {
	type args struct {
		tp      string
		address string
		f       func(string) *services.ORouter
		opts    []api.Option
	}
	tests := []struct {
		name string
		args args
		want *httpBuilder
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := newHTTP(tt.args.tp, tt.args.address, tt.args.f, tt.args.opts...); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("newHTTP() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_httpBuilder_Load(t *testing.T) {
	type fields struct {
		CustomerBuilder CustomerBuilder
		tp              string
		fnGetRouter     func(string) *services.ORouter
	}
	tests := []struct {
		name   string
		fields fields
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			b := &httpBuilder{
				CustomerBuilder: tt.fields.CustomerBuilder,
				tp:              tt.fields.tp,
				fnGetRouter:     tt.fields.fnGetRouter,
			}
			b.Load()
		})
	}
}

func Test_httpBuilder_Jwt(t *testing.T) {
	type fields struct {
		CustomerBuilder CustomerBuilder
		tp              string
		fnGetRouter     func(string) *services.ORouter
	}
	type args struct {
		opts []jwt.Option
	}
	tests := []struct {
		name   string
		fields fields
		args   args
		want   *httpBuilder
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			b := &httpBuilder{
				CustomerBuilder: tt.fields.CustomerBuilder,
				tp:              tt.fields.tp,
				fnGetRouter:     tt.fields.fnGetRouter,
			}
			if got := b.Jwt(tt.args.opts...); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("httpBuilder.Jwt() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_httpBuilder_APIKEY(t *testing.T) {
	type fields struct {
		CustomerBuilder CustomerBuilder
		tp              string
		fnGetRouter     func(string) *services.ORouter
	}
	type args struct {
		secret string
		opts   []apikey.Option
	}
	tests := []struct {
		name   string
		fields fields
		args   args
		want   *httpBuilder
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			b := &httpBuilder{
				CustomerBuilder: tt.fields.CustomerBuilder,
				tp:              tt.fields.tp,
				fnGetRouter:     tt.fields.fnGetRouter,
			}
			if got := b.APIKEY(tt.args.secret, tt.args.opts...); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("httpBuilder.APIKEY() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_httpBuilder_Basic(t *testing.T) {
	type fields struct {
		CustomerBuilder CustomerBuilder
		tp              string
		fnGetRouter     func(string) *services.ORouter
	}
	type args struct {
		opts []basic.Option
	}
	tests := []struct {
		name   string
		fields fields
		args   args
		want   *httpBuilder
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			b := &httpBuilder{
				CustomerBuilder: tt.fields.CustomerBuilder,
				tp:              tt.fields.tp,
				fnGetRouter:     tt.fields.fnGetRouter,
			}
			if got := b.Basic(tt.args.opts...); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("httpBuilder.Basic() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_httpBuilder_WhiteList(t *testing.T) {
	type fields struct {
		CustomerBuilder CustomerBuilder
		tp              string
		fnGetRouter     func(string) *services.ORouter
	}
	type args struct {
		opts []whitelist.Option
	}
	tests := []struct {
		name   string
		fields fields
		args   args
		want   *httpBuilder
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			b := &httpBuilder{
				CustomerBuilder: tt.fields.CustomerBuilder,
				tp:              tt.fields.tp,
				fnGetRouter:     tt.fields.fnGetRouter,
			}
			if got := b.WhiteList(tt.args.opts...); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("httpBuilder.WhiteList() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_httpBuilder_BlackList(t *testing.T) {
	type fields struct {
		CustomerBuilder CustomerBuilder
		tp              string
		fnGetRouter     func(string) *services.ORouter
	}
	type args struct {
		opts []blacklist.Option
	}
	tests := []struct {
		name   string
		fields fields
		args   args
		want   *httpBuilder
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			b := &httpBuilder{
				CustomerBuilder: tt.fields.CustomerBuilder,
				tp:              tt.fields.tp,
				fnGetRouter:     tt.fields.fnGetRouter,
			}
			if got := b.BlackList(tt.args.opts...); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("httpBuilder.BlackList() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_httpBuilder_Ras(t *testing.T) {
	type fields struct {
		CustomerBuilder CustomerBuilder
		tp              string
		fnGetRouter     func(string) *services.ORouter
	}
	type args struct {
		opts []ras.Option
	}
	tests := []struct {
		name   string
		fields fields
		args   args
		want   *httpBuilder
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			b := &httpBuilder{
				CustomerBuilder: tt.fields.CustomerBuilder,
				tp:              tt.fields.tp,
				fnGetRouter:     tt.fields.fnGetRouter,
			}
			if got := b.Ras(tt.args.opts...); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("httpBuilder.Ras() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_httpBuilder_Header(t *testing.T) {
	type fields struct {
		CustomerBuilder CustomerBuilder
		tp              string
		fnGetRouter     func(string) *services.ORouter
	}
	type args struct {
		opts []header.Option
	}
	tests := []struct {
		name   string
		fields fields
		args   args
		want   *httpBuilder
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			b := &httpBuilder{
				CustomerBuilder: tt.fields.CustomerBuilder,
				tp:              tt.fields.tp,
				fnGetRouter:     tt.fields.fnGetRouter,
			}
			if got := b.Header(tt.args.opts...); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("httpBuilder.Header() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_httpBuilder_Metric(t *testing.T) {
	type fields struct {
		CustomerBuilder CustomerBuilder
		tp              string
		fnGetRouter     func(string) *services.ORouter
	}
	type args struct {
		host string
		db   string
		cron string
		opts []metric.Option
	}
	tests := []struct {
		name   string
		fields fields
		args   args
		want   *httpBuilder
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			b := &httpBuilder{
				CustomerBuilder: tt.fields.CustomerBuilder,
				tp:              tt.fields.tp,
				fnGetRouter:     tt.fields.fnGetRouter,
			}
			if got := b.Metric(tt.args.host, tt.args.db, tt.args.cron, tt.args.opts...); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("httpBuilder.Metric() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_httpBuilder_Static(t *testing.T) {
	type fields struct {
		CustomerBuilder CustomerBuilder
		tp              string
		fnGetRouter     func(string) *services.ORouter
	}
	type args struct {
		opts []static.Option
	}
	tests := []struct {
		name   string
		fields fields
		args   args
		want   *httpBuilder
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			b := &httpBuilder{
				CustomerBuilder: tt.fields.CustomerBuilder,
				tp:              tt.fields.tp,
				fnGetRouter:     tt.fields.fnGetRouter,
			}
			if got := b.Static(tt.args.opts...); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("httpBuilder.Static() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_httpBuilder_Limit(t *testing.T) {
	type fields struct {
		CustomerBuilder CustomerBuilder
		tp              string
		fnGetRouter     func(string) *services.ORouter
	}
	type args struct {
		opts []limiter.Option
	}
	tests := []struct {
		name   string
		fields fields
		args   args
		want   *httpBuilder
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			b := &httpBuilder{
				CustomerBuilder: tt.fields.CustomerBuilder,
				tp:              tt.fields.tp,
				fnGetRouter:     tt.fields.fnGetRouter,
			}
			if got := b.Limit(tt.args.opts...); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("httpBuilder.Limit() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_httpBuilder_Gray(t *testing.T) {
	type fields struct {
		CustomerBuilder CustomerBuilder
		tp              string
		fnGetRouter     func(string) *services.ORouter
	}
	type args struct {
		opts []gray.Option
	}
	tests := []struct {
		name   string
		fields fields
		args   args
		want   *httpBuilder
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			b := &httpBuilder{
				CustomerBuilder: tt.fields.CustomerBuilder,
				tp:              tt.fields.tp,
				fnGetRouter:     tt.fields.fnGetRouter,
			}
			if got := b.Gray(tt.args.opts...); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("httpBuilder.Gray() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_httpBuilder_Render(t *testing.T) {
	type fields struct {
		CustomerBuilder CustomerBuilder
		tp              string
		fnGetRouter     func(string) *services.ORouter
	}
	type args struct {
		opts []render.Option
	}
	tests := []struct {
		name   string
		fields fields
		args   args
		want   *httpBuilder
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			b := &httpBuilder{
				CustomerBuilder: tt.fields.CustomerBuilder,
				tp:              tt.fields.tp,
				fnGetRouter:     tt.fields.fnGetRouter,
			}
			if got := b.Render(tt.args.opts...); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("httpBuilder.Render() = %v, want %v", got, tt.want)
			}
		})
	}
}
