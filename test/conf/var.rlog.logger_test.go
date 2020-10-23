package conf

import (
	"fmt"
	"testing"

	"github.com/micro-plat/hydra/conf"
	"github.com/micro-plat/hydra/conf/vars/rlog"
	"github.com/micro-plat/hydra/test/assert"
	"github.com/micro-plat/hydra/test/mocks"
)

func TestRLogNew(t *testing.T) {

	tests := []struct {
		name    string
		service string
		opts    []rlog.Option
		want    *rlog.Layout
	}{
		{
			name:    "新增-无option",
			service: "/rlog",
			want: &rlog.Layout{
				Level:   "Info",
				Service: "/rlog",
				Layout:  rlog.DefaultLayout,
				Disable: false,
			},
		},
		{
			name:    "新增-WithLayout",
			service: "/rlog",
			opts: []rlog.Option{
				rlog.WithLayout("customerlayout"),
			},
			want: &rlog.Layout{
				Level:   "Info",
				Service: "/rlog",
				Layout:  "customerlayout",
				Disable: false,
			},
		},
		{
			name:    "新增-WithDisable",
			service: "/rlog",
			opts: []rlog.Option{
				rlog.WithLayout("customerlayout"),
				rlog.WithDisable(),
			},
			want: &rlog.Layout{
				Level:   "Info",
				Service: "/rlog",
				Layout:  "customerlayout",
				Disable: true,
			},
		},
		{
			name:    "新增-WithEnable",
			service: "/rlog",
			opts: []rlog.Option{
				rlog.WithLayout("customerlayout"),
				rlog.WithEnable(),
			},
			want: &rlog.Layout{
				Level:   "Info",
				Service: "/rlog",
				Layout:  "customerlayout",
				Disable: false,
			},
		},
		{
			name:    "新增-WithInfo",
			service: "/rlog",
			opts: []rlog.Option{
				rlog.WithLayout("customerlayout"),
				rlog.WithEnable(),
				rlog.WithInfo(),
			},
			want: &rlog.Layout{
				Level:   "Info",
				Service: "/rlog",
				Layout:  "customerlayout",
				Disable: false,
			},
		},
		{
			name:    "新增-WithOff",
			service: "/rlog",
			opts: []rlog.Option{
				rlog.WithLayout("customerlayout"),
				rlog.WithEnable(),
				rlog.WithOff(),
			},
			want: &rlog.Layout{
				Level:   "Off",
				Service: "/rlog",
				Layout:  "customerlayout",
				Disable: false,
			},
		},
		{
			name:    "新增-WithWarn",
			service: "/rlog",
			opts: []rlog.Option{
				rlog.WithLayout("customerlayout"),
				rlog.WithEnable(),
				rlog.WithWarn(),
			},
			want: &rlog.Layout{
				Level:   "Warn",
				Service: "/rlog",
				Layout:  "customerlayout",
				Disable: false,
			},
		},
		{
			name:    "新增-WithError",
			service: "/rlog",
			opts: []rlog.Option{
				rlog.WithLayout("customerlayout"),
				rlog.WithEnable(),
				rlog.WithError(),
			},
			want: &rlog.Layout{
				Level:   "Error",
				Service: "/rlog",
				Layout:  "customerlayout",
				Disable: false,
			},
		},
		{
			name:    "新增-WithFatal",
			service: "/rlog",
			opts: []rlog.Option{
				rlog.WithLayout("customerlayout"),
				rlog.WithEnable(),
				rlog.WithFatal(),
			},
			want: &rlog.Layout{
				Level:   "Fatal",
				Service: "/rlog",
				Layout:  "customerlayout",
				Disable: false,
			},
		},
		{
			name:    "新增-WithAll",
			service: "/rlog",
			opts: []rlog.Option{
				rlog.WithLayout("customerlayout"),
				rlog.WithEnable(),
				rlog.WithAll(),
			},
			want: &rlog.Layout{
				Level:   "All",
				Service: "/rlog",
				Layout:  "customerlayout",
				Disable: false,
			},
		},
	}
	for _, tt := range tests {
		got := rlog.New(tt.service, tt.opts...)
		assert.Equal(t, tt.want, got, tt.name)
	}
}

func TestRLogGetConf(t *testing.T) {

	type args struct {
		cnfData []byte
		logName string
		version int32
	}
	tests := []struct {
		name     string
		args     args
		want     *rlog.Layout
		IsNilErr bool
	}{
		{
			name: "测试-var中无该配置",
			args: args{
				logName: "xxrlog",
				cnfData: []byte(`{"level":"Info","service":"/rlog"}`),
			},
			want: &rlog.Layout{
				Layout:  rlog.DefaultLayout,
				Disable: true,
			},
			IsNilErr: true,
		},

		{
			name: "测试-var中有配置-不满足规则",
			args: args{
				logName: "rlog",
				cnfData: []byte(`{"level":"xInfo","service":"/rlog"}`),
			},
			want:     nil,
			IsNilErr: false,
		},
		{
			name: "测试-var中有配置-正确",
			args: args{
				logName: "rlog",
				cnfData: []byte(`{"level":"Info","service":"/rlog"}`),
			},
			want: &rlog.Layout{
				Level:   "Info",
				Service: "/rlog",
				Layout:  rlog.DefaultLayout,
				Disable: false,
			},
			IsNilErr: true,
		},

		// TODO: Add test cases.
	}
	for _, tt := range tests {

		rawCnf, err := conf.NewRawConfByJson(tt.args.cnfData, tt.args.version)
		fmt.Println(tt.name)
		cnf := &mocks.MockVarConf{
			Version: tt.args.version,
			ConfData: map[string]map[string]*conf.RawConf{
				"app": map[string]*conf.RawConf{
					tt.args.logName: rawCnf,
				},
			},
		}

		got, err := rlog.GetConf(cnf)
		//fmt.Println("rlog.GetConf:", got, err)
		assert.IsNil(t, tt.IsNilErr, err, tt.name)

		assert.Equal(t, tt.want, got, tt.name)
	}
}
