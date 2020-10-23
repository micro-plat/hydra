package conf

import (
	"fmt"
	"testing"

	"github.com/micro-plat/hydra/conf"
	"github.com/micro-plat/hydra/conf/vars/db"
	"github.com/micro-plat/hydra/test/assert"
	"github.com/micro-plat/hydra/test/mocks"
)

func TestDBNew(t *testing.T) {

	tests := []struct {
		name       string
		provider   string
		connString string
		opts       []db.Option
		want       *db.DB
	}{
		{
			name:       "测试新增-无OPTION",
			provider:   "oracle",
			connString: "zhjy/123456@orcl136",
			want: &db.DB{
				Provider:   "oracle",
				ConnString: "zhjy/123456@orcl136",
				MaxOpen:    10,
				MaxIdle:    3,
				LifeTime:   600,
			},
		},
		{
			name:       "测试新增-WithConnect",
			provider:   "oracle",
			connString: "zhjy/123456@orcl136",
			opts: []db.Option{
				db.WithConnect(11, 22, 33),
			},
			want: &db.DB{
				Provider:   "oracle",
				ConnString: "zhjy/123456@orcl136",
				MaxOpen:    11,
				MaxIdle:    22,
				LifeTime:   33,
			},
		},
	}
	for _, tt := range tests {
		got := db.New(tt.provider, tt.connString, tt.opts...)
		assert.Equal(t, tt.want, got, tt.name)
	}
}

func TestDBGetConf(t *testing.T) {
	type args struct {
		cnfData []byte
		version int32
		tp      string
		name    string
	}
	connectString := `{"provider":"oracle","connString":"zhjy/123456@orcl136","maxOpen":11,"maxIdle":22,"lifeTime":33}`
	tests := []struct {
		name     string
		args     args
		want     *db.DB
		IsNilErr bool
	}{
		{
			name: "测试-var中无该配置",
			args: args{
				cnfData: []byte(connectString),
				version: 1,
				tp:      "db",
				name:    "xoracle",
			},
			want:     nil,
			IsNilErr: false,
		},
		{
			name: "测试-var中有配置",
			args: args{
				cnfData: []byte(connectString),
				version: 1,
				tp:      "db",
				name:    "oracle",
			},
			want: &db.DB{
				Provider:   "oracle",
				ConnString: connectString,
				MaxOpen:    11,
				MaxIdle:    22,
				LifeTime:   33,
			},
			IsNilErr: true,
		},

		// TODO: Add test cases.
	}
	for _, tt := range tests {
		fmt.Println(tt.name)

		rawCnf, err := conf.NewRawConfByJson(tt.args.cnfData, tt.args.version)

		assert.IsNil(t, true, err, tt.name)
		cnf := &mocks.MockVarConf{
			Version: tt.args.version,
			ConfData: map[string]map[string]*conf.RawConf{
				"db": map[string]*conf.RawConf{
					"oracle": rawCnf,
				},
			},
		}
		fmt.Println(tt.args.tp, tt.args.name)
		got, err := db.GetConf(cnf, tt.args.tp, tt.args.name)
		//fmt.Println("err:", got, err)
		assert.IsNil(t, tt.IsNilErr, err, tt.name)
		fmt.Printf("want:%+v;got:%+v\r\n", tt.want, got)
		//assert.Equal(t, tt.want, got, tt.name)
		if err != nil && tt.want != nil {
			assert.Equal(t, tt.want.Provider, got.Provider, tt.name)
			assert.Equal(t, tt.want.ConnString, got.ConnString, tt.name)
			assert.Equal(t, tt.want.MaxOpen, got.MaxOpen, tt.name)
			assert.Equal(t, tt.want.MaxIdle, got.MaxIdle, tt.name)
			assert.Equal(t, tt.want.LifeTime, got.LifeTime, tt.name)
		}
	}
}
