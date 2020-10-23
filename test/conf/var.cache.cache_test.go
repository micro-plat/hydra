package conf

import (
	"testing"

	"github.com/micro-plat/hydra/conf"
	"github.com/micro-plat/hydra/conf/vars/cache"
	"github.com/micro-plat/hydra/test/assert"
	"github.com/micro-plat/hydra/test/mocks"
)

func TestCacheNew(t *testing.T) {

	tests := []struct {
		name  string
		proto string
		raw   []byte
		want  *cache.Cache
	}{
		{
			name:  "测试新增",
			proto: "redis",
			raw:   []byte(`{"proto":"redis","addrs":["192.168.5.79:6379"]}`),
			want: &cache.Cache{
				Proto: "redis",
				Raw:   []byte(`{"proto":"redis","addrs":["192.168.5.79:6379"]}`),
			},
		},
	}
	for _, tt := range tests {
		got := cache.New(tt.proto, tt.raw)
		assert.Equal(t, tt.want.Proto, got.Proto, tt.name)
		assert.Equal(t, tt.want.Raw, got.Raw, tt.name)

	}
}

func TestCacheGetConf(t *testing.T) {
	type args struct {
		cnfData []byte
		version int32
		tp      string
		name    string
	}
	tests := []struct {
		name     string
		args     args
		want     *cache.Cache
		IsNilErr bool
	}{
		{
			name: "测试-var中无该配置",
			args: args{
				cnfData: []byte(`{"proto":"redis","addrs":["192.168.5.79:6379"]}`),
				tp:      "cache",
				name:    "xredis",
			},
			want:     nil,
			IsNilErr: false,
		},
		{
			name: "测试-var中有配置",
			args: args{
				cnfData: []byte(`{"proto":"redis","addrs":["192.168.5.79:6379"]}`),
				tp:      "cache",
				name:    "redis",
			},
			want: &cache.Cache{
				Proto: "redis",
				Raw:   []byte(`{"proto":"redis","addrs":["192.168.5.79:6379"]}`),
			},
			IsNilErr: true,
		},

		// TODO: Add test cases.
	}
	for _, tt := range tests {
		rawCnf, err := conf.NewRawConfByJson(tt.args.cnfData, tt.args.version)

		cnf := &mocks.MockVarConf{
			Version: tt.args.version,
			ConfData: map[string]map[string]*conf.RawConf{
				"cache": map[string]*conf.RawConf{
					"redis": rawCnf,
				},
			},
		}

		got, err := cache.GetConf(cnf, tt.args.tp, tt.args.name)

		assert.IsNil(t, tt.IsNilErr, err, tt.name)

		assert.Equal(t, tt.want, got, tt.name)
	}
}
