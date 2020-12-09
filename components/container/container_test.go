package container

import (
	"reflect"
	"testing"

	"github.com/micro-plat/hydra/conf"
	"github.com/micro-plat/lib4go/assert"
	"github.com/micro-plat/lib4go/concurrent/cmap"
)

func TestNewContainer(t *testing.T) {
	l := NewContainer()
	if !reflect.DeepEqual(l, &Container{
		cache:     cmap.New(8),
		histories: newHistories(),
	}) {
		t.Error("NewContainer() didn't return *Container")
	}
}

func TestContainer_GetOrCreate(t *testing.T) {

	c := NewContainer()
	type args struct {
		typ     string
		name    string
		creator func(conf *conf.RawConf) (interface{}, error)
	}
	creator := func(conf *conf.RawConf) (interface{}, error) {
		return nil, nil
	}
	tests := []struct {
		name    string
		args    args
		want    interface{}
		wantErr bool
	}{
		{name: "1. 初始化容器内容db对象", args: args{typ: "db", name: "db", creator: creator}, want: nil, wantErr: true},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := c.GetOrCreate(tt.args.typ, tt.args.name, tt.args.creator)
			assert.Nil(t, tt.wantErr, err, tt.name)
			assert.Equal(t, tt.wantErr, err == nil, tt.name)
			assert.Equal(t, tt.want, got, tt.name)
		})
	}
}
