package conf

import (
	"encoding/json"
	"fmt"
	"testing"

	"github.com/micro-plat/hydra/registry"
	"github.com/micro-plat/lib4go/assert"
	"github.com/micro-plat/lib4go/security/md5"
	"github.com/micro-plat/lib4go/types"
)

var _ IServerConf = &TMainConf{}

//TMainConf 服务器主配置
type TMainConf struct {
	rootConf    *RawConf
	rootVersion int32
	subConfs    map[string]RawConf
	IServerPub
	registry     registry.IRegistry
	ClusterNames []string
}

func NewTMainConf(rootVersion int32, data map[string]interface{}) IServerConf {
	raw, _ := json.Marshal(data)
	nRawConf := RawConf{
		XMap:      data,
		version:   rootVersion,
		raw:       raw,
		signature: md5.EncryptBytes(raw),
	}

	data["subc"] = "123456"
	raw1, _ := json.Marshal(data)
	subRawConf := RawConf{
		XMap:      data,
		version:   rootVersion,
		raw:       raw1,
		signature: md5.EncryptBytes(raw1),
	}

	subm := map[string]RawConf{
		"keysub": subRawConf,
	}
	return &TMainConf{rootConf: &nRawConf, rootVersion: rootVersion, subConfs: subm}
}

func NewTMainConf1(rootVersion int32, keySub map[string]string) IServerConf {

	subConf := map[string]RawConf{}
	for k, v := range keySub {
		data := map[string]interface{}{"value": v}
		raw1, _ := json.Marshal(data)
		subConf[k] = RawConf{
			XMap:      data,
			version:   rootVersion,
			raw:       raw1,
			signature: md5.EncryptBytes(raw1),
		}

	}
	return &TMainConf{rootConf: nil, rootVersion: rootVersion, subConfs: subConf}
}

//IsTrace 是否跟踪请求或响应
func (c *TMainConf) IsTrace() bool {
	return true
}

//GetRegistry 获取注册中心
func (c *TMainConf) GetRegistry() registry.IRegistry {
	return nil
}

//IsStarted 当前服务是否已启动
func (c *TMainConf) IsStarted() bool {
	return true
}

//GetVersion 获取版本号
func (c *TMainConf) GetVersion() int32 {
	return c.rootVersion
}

//GetRootConf 获取当前主配置
func (c *TMainConf) GetMainConf() *RawConf {
	return c.rootConf
}

//GetMainObject 获取主配置信息
func (c *TMainConf) GetMainObject(v interface{}) (int32, error) {

	return 0, nil
}

//GetSubConf 指定子配置
func (c *TMainConf) GetSubConf(name string) (*RawConf, error) {
	if v, ok := c.subConfs[name]; ok {
		return &v, nil
	}
	return nil, ErrNoSetting
}

//GetCluster 获取集群信息
func (c *TMainConf) GetCluster(clustName ...string) (ICluster, error) {
	return nil, nil
}

//GetSubObject 获取子配置信息
func (c *TMainConf) GetSubObject(name string, v interface{}) (int32, error) {
	conf, err := c.GetSubConf(name)
	if err != nil {
		return 0, err
	}

	if err := conf.ToStruct(&v); err != nil {
		err = fmt.Errorf("获取%s配置失败:%v", name, err)
		return 0, err
	}
	return conf.GetVersion(), nil
}

//Has 是否存在子配置
func (c *TMainConf) Has(names ...string) bool {
	return false
}

//Iter 迭代所有配置
func (c *TMainConf) Iter(f func(path string, conf *RawConf) bool) {
}

//Close 关闭清理资源
func (c *TMainConf) Close() error {
	return nil
}

func (c *TMainConf) GetClusterNames() []string {
	return c.ClusterNames
}

func TestComparer_Update(t *testing.T) {
	type fields struct {
		oconf      IServerConf
		nconf      IServerConf
		valueNames []string
		subNames   []string
	}
	type args struct {
		n IServerConf
	}

	tests := []struct {
		name   string
		fields fields
		args   args
	}{
		{name: "1. compare-update-不存在nconf数据用例", fields: fields{oconf: nil, nconf: nil}, args: args{n: NewTMainConf(1, map[string]interface{}{"xx": "11"})}},
		{name: "2. compare-update-存在nconf数据用例", fields: fields{oconf: nil, nconf: NewTMainConf(1, map[string]interface{}{"xx": "11"})}, args: args{n: NewTMainConf(2, map[string]interface{}{"xx": "22"})}},
		{name: "3. compare-update-存在oconf,nconf数据用例", fields: fields{oconf: NewTMainConf(0, map[string]interface{}{"xx": "1212"}), nconf: NewTMainConf(1, map[string]interface{}{"xx": "11"})}, args: args{n: NewTMainConf(2, map[string]interface{}{"xx": "22"})}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := &Comparer{
				oconf:      tt.fields.oconf,
				nconf:      tt.fields.nconf,
				valueNames: tt.fields.valueNames,
				subNames:   tt.fields.subNames,
			}
			s.Update(tt.args.n)
			switch tt.name[:2] {
			case "1.":
				if s.oconf != nil || s.nconf == nil {
					t.Errorf("用例[%s]更新nconf,nil判断失败", tt.name)
				}
				if s.nconf != tt.args.n {
					t.Errorf("用例[%s]更新nconf数据失败", tt.name)
				}
			case "2.":
				if s.oconf == nil || s.nconf == nil {
					t.Errorf("用例[%s]更新nconf,nil判断失败", tt.name)
				}
				if s.nconf != tt.args.n {
					t.Errorf("用例[%s]更新nconf数据失败", tt.name)
				}
			case "3.":
				if s.oconf == nil || s.nconf == nil {
					t.Errorf("用例[%s]更新nconf,nil判断失败", tt.name)
				}
				if s.oconf != tt.fields.nconf {
					t.Errorf("用例[%s]更新oconf记录数据失败", tt.name)
				}
				if s.nconf != tt.args.n {
					t.Errorf("用例[%s]更新nconf数据失败", tt.name)
				}
			default:
				t.Errorf("用例[%s]没有做断言判断结果", tt.name)
			}
		})
	}
}

func TestComparer_IsChanged(t *testing.T) {
	//该方法是通过版本号进行比较   所以只需要mock版本号信息
	type fields struct {
		oconf      IServerConf
		nconf      IServerConf
		valueNames []string
		subNames   []string
	}
	tests := []struct {
		name   string
		fields fields
		want   bool
	}{
		{name: "1. compare-IsChanged-配置不存在", fields: fields{oconf: nil, nconf: nil}, want: false},
		{name: "2. compare-IsChanged-new配置不存在", fields: fields{oconf: NewTMainConf(1, map[string]interface{}{"xxx": "ww"}), nconf: nil}, want: false},
		{name: "3. compare-IsChanged-版本号相同,内容不同", fields: fields{oconf: NewTMainConf(1, map[string]interface{}{"xxx": "ww"}), nconf: NewTMainConf(1, map[string]interface{}{"xxx": "ww11"})}, want: false},
		{name: "4. compare-IsChanged-版本号相同,内容相同", fields: fields{oconf: NewTMainConf(1, map[string]interface{}{"xxx": "ww11"}), nconf: NewTMainConf(1, map[string]interface{}{"xxx": "ww11"})}, want: false},
		{name: "5. compare-IsChanged-版本号不同,内容不同", fields: fields{oconf: NewTMainConf(0, map[string]interface{}{"xxx": "ww"}), nconf: NewTMainConf(1, map[string]interface{}{"xxx": "ww11"})}, want: true},
		{name: "6. compare-IsChanged-版本号不同,内容相同", fields: fields{oconf: NewTMainConf(0, map[string]interface{}{"xxx": "ww11"}), nconf: NewTMainConf(1, map[string]interface{}{"xxx": "ww11"})}, want: true},
	}
	for _, tt := range tests {
		s := &Comparer{
			oconf:      tt.fields.oconf,
			nconf:      tt.fields.nconf,
			valueNames: tt.fields.valueNames,
			subNames:   tt.fields.subNames,
		}
		if got := s.IsChanged(); got != tt.want {
			t.Errorf("Comparer.IsChanged() = %v, want %v", got, tt.want)
		}
	}
}

//都是要比较所有的监控数组的值是否发生变化
func TestComparer_IsSubConfChanged(t *testing.T) {
	type fields struct {
		oconf      IServerConf
		nconf      IServerConf
		valueNames []string
		subNames   []string
	}
	type args struct {
		names []string
	}
	tests := []struct {
		name          string
		fields        fields
		args          args
		wantIsChanged bool
	}{
		{name: "1.1. compare-IsSubConfChanged-版本号相同,内容不同", fields: fields{oconf: NewTMainConf1(1, map[string]string{"xx": "123455"}), nconf: NewTMainConf1(1, map[string]string{"xx": "44444"})}, args: args{names: []string{"xx"}}, wantIsChanged: false},
		{name: "1.2. compare-IsSubConfChanged-版本号相同,内容相同", fields: fields{oconf: NewTMainConf1(1, map[string]string{"xx": "123455"}), nconf: NewTMainConf1(1, map[string]string{"xx": "123455"})}, args: args{names: []string{"xx"}}, wantIsChanged: false},
		{name: "1.3. compare-IsSubConfChanged-版本号相同,内容都为空", fields: fields{oconf: NewTMainConf1(1, map[string]string{}), nconf: NewTMainConf1(1, map[string]string{})}, args: args{names: []string{}}, wantIsChanged: false},
		{name: "1.4. compare-IsSubConfChanged-版本号相同,内容都为空1", fields: fields{oconf: NewTMainConf1(1, map[string]string{}), nconf: NewTMainConf1(1, map[string]string{})}, args: args{names: []string{"xx"}}, wantIsChanged: false},

		{name: "2.1. compare-IsSubConfChanged-版本号不同,oldversion > newversion,内容不同", fields: fields{oconf: NewTMainConf1(1, map[string]string{"xx": "123455"}), nconf: NewTMainConf1(0, map[string]string{"xx": "44444"})}, args: args{names: []string{"xx"}}, wantIsChanged: false},
		{name: "2.2. compare-IsSubConfChanged-版本号不同,oldversion > newversion,内容相同", fields: fields{oconf: NewTMainConf1(1, map[string]string{"xx": "123455"}), nconf: NewTMainConf1(0, map[string]string{"xx": "123455"})}, args: args{names: []string{"xx"}}, wantIsChanged: false},
		{name: "2.3. compare-IsSubConfChanged-版本号不同,oldversion > newversion,内容都为空", fields: fields{oconf: NewTMainConf1(1, map[string]string{}), nconf: NewTMainConf1(0, map[string]string{})}, args: args{names: []string{}}, wantIsChanged: false},
		{name: "2.4. compare-IsSubConfChanged-版本号不同,oldversion > newversion,内容都为空1", fields: fields{oconf: NewTMainConf1(1, map[string]string{}), nconf: NewTMainConf1(0, map[string]string{})}, args: args{names: []string{"xx"}}, wantIsChanged: false},
		{name: "2.5. compare-IsSubConfChanged-版本号不同,oldversion < newversion,内容不同", fields: fields{oconf: NewTMainConf1(1, map[string]string{"xx": "123455"}), nconf: NewTMainConf1(0, map[string]string{"xx": "44444"})}, args: args{names: []string{"xx"}}, wantIsChanged: false},
		{name: "2.6. compare-IsSubConfChanged-版本号不同,oldversion < newversion,内容相同", fields: fields{oconf: NewTMainConf1(1, map[string]string{"xx": "123455"}), nconf: NewTMainConf1(0, map[string]string{"xx": "123455"})}, args: args{names: []string{"xx"}}, wantIsChanged: false},
		{name: "2.7. compare-IsSubConfChanged-版本号不同,oldversion < newversion,内容都为空", fields: fields{oconf: NewTMainConf1(1, map[string]string{}), nconf: NewTMainConf1(0, map[string]string{})}, args: args{names: []string{}}, wantIsChanged: false},
		{name: "2.8. compare-IsSubConfChanged-版本号不同,oldversion < newversion,内容都为空1", fields: fields{oconf: NewTMainConf1(1, map[string]string{}), nconf: NewTMainConf1(0, map[string]string{})}, args: args{names: []string{"xx"}}, wantIsChanged: false},
	}
	for _, tt := range tests {
		s := &Comparer{
			oconf:      tt.fields.oconf,
			nconf:      tt.fields.nconf,
			valueNames: tt.fields.valueNames,
			subNames:   tt.fields.subNames,
		}
		if gotIsChanged := s.IsSubConfChanged(tt.args.names...); gotIsChanged != tt.wantIsChanged {
			t.Errorf("Comparer.IsSubConfChanged() = %v, want %v", gotIsChanged, tt.wantIsChanged)
		}
	}
}

//都是要比较所有的监控数组的值是否发生变化
func TestComparer_IsValueChanged(t *testing.T) {
	type fields struct {
		oconf      IServerConf
		nconf      IServerConf
		valueNames []string
		subNames   []string
	}
	type args struct {
		names []string
	}
	tests := []struct {
		name          string
		fields        fields
		args          args
		wantIsChanged bool
	}{
		{name: "1.1. compare-IsValueChanged-版本号相同,内容不同", fields: fields{oconf: NewTMainConf(1, map[string]interface{}{"xx": "123455"}), nconf: NewTMainConf(1, map[string]interface{}{"xx": "44444"})}, args: args{names: []string{"xx"}}, wantIsChanged: false},
		{name: "1.2. compare-IsValueChanged-版本号相同,内容相同", fields: fields{oconf: NewTMainConf(1, map[string]interface{}{"xx": "123455"}), nconf: NewTMainConf(1, map[string]interface{}{"xx": "123455"})}, args: args{names: []string{"xx"}}, wantIsChanged: false},
		{name: "1.3. compare-IsValueChanged-版本号相同,内容都为空", fields: fields{oconf: NewTMainConf(1, map[string]interface{}{}), nconf: NewTMainConf(1, map[string]interface{}{})}, args: args{names: []string{}}, wantIsChanged: false},
		{name: "1.4. compare-IsValueChanged-版本号相同,内容都为空1", fields: fields{oconf: NewTMainConf(1, map[string]interface{}{}), nconf: NewTMainConf(1, map[string]interface{}{})}, args: args{names: []string{"xx"}}, wantIsChanged: false},

		{name: "2.1. compare-IsValueChanged-版本号不同,oldversion > newversion,内容不同", fields: fields{oconf: NewTMainConf(1, map[string]interface{}{"xx": "123455"}), nconf: NewTMainConf(0, map[string]interface{}{"xx": "44444"})}, args: args{names: []string{"xx"}}, wantIsChanged: false},
		{name: "2.2. compare-IsValueChanged-版本号不同,oldversion > newversion,内容相同", fields: fields{oconf: NewTMainConf(1, map[string]interface{}{"xx": "123455"}), nconf: NewTMainConf(0, map[string]interface{}{"xx": "123455"})}, args: args{names: []string{"xx"}}, wantIsChanged: false},
		{name: "2.3. compare-IsValueChanged-版本号不同,oldversion > newversion,内容都为空", fields: fields{oconf: NewTMainConf(1, map[string]interface{}{}), nconf: NewTMainConf(0, map[string]interface{}{})}, args: args{names: []string{}}, wantIsChanged: false},
		{name: "2.4. compare-IsValueChanged-版本号不同,oldversion > newversion,内容都为空1", fields: fields{oconf: NewTMainConf(1, map[string]interface{}{}), nconf: NewTMainConf(0, map[string]interface{}{})}, args: args{names: []string{"xx"}}, wantIsChanged: false},
		{name: "2.5. compare-IsValueChanged-版本号不同,oldversion < newversion,内容不同", fields: fields{oconf: NewTMainConf(1, map[string]interface{}{"xx": "123455"}), nconf: NewTMainConf(0, map[string]interface{}{"xx": "44444"})}, args: args{names: []string{"xx"}}, wantIsChanged: false},
		{name: "2.6. compare-IsValueChanged-版本号不同,oldversion < newversion,内容相同", fields: fields{oconf: NewTMainConf(1, map[string]interface{}{"xx": "123455"}), nconf: NewTMainConf(0, map[string]interface{}{"xx": "123455"})}, args: args{names: []string{"xx"}}, wantIsChanged: false},
		{name: "2.7. compare-IsValueChanged-版本号不同,oldversion < newversion,内容都为空", fields: fields{oconf: NewTMainConf(1, map[string]interface{}{}), nconf: NewTMainConf(0, map[string]interface{}{})}, args: args{names: []string{}}, wantIsChanged: false},
		{name: "2.8. compare-IsValueChanged-版本号不同,oldversion < newversion,内容都为空1", fields: fields{oconf: NewTMainConf(1, map[string]interface{}{}), nconf: NewTMainConf(0, map[string]interface{}{})}, args: args{names: []string{"xx"}}, wantIsChanged: false},
	}
	for _, tt := range tests {
		s := &Comparer{
			oconf:      tt.fields.oconf,
			nconf:      tt.fields.nconf,
			valueNames: tt.fields.valueNames,
			subNames:   tt.fields.subNames,
		}
		if gotIsChanged := s.IsValueChanged(tt.args.names...); gotIsChanged != tt.wantIsChanged {
			t.Errorf("Comparer.IsValueChanged() = %v, want %v", gotIsChanged, tt.wantIsChanged)
		}
	}
}

func TestNewComparer(t *testing.T) {
	type args struct {
		oconf      IServerConf
		valueNames []string
		subNames   []string
	}
	tests := []struct {
		name    string
		args    args
		want    *Comparer
		wantErr string
	}{
		{name: "1. NewComparer-oldconf为空", args: args{oconf: nil, valueNames: []string{"valueNames1"}, subNames: []string{"subNames1"}}, want: nil, wantErr: "配置不能为空"},
		{name: "2. NewComparer-实体对象初始化", args: args{oconf: NewTMainConf(0, map[string]interface{}{"xx": "123455"}), valueNames: []string{"valueNames1"}, subNames: []string{"subNames1"}},
			want: &Comparer{oconf: NewTMainConf(0, map[string]interface{}{"xx": "123455"}), valueNames: []string{"valueNames1"}, subNames: []string{"subNames1"}}, wantErr: "配置不能为空"},
	}
	for _, tt := range tests {
		func() {
			defer func() {
				e := recover()
				if e != nil {
					assert.Equal(t, tt.wantErr, types.GetString(e), tt.name+",pinic")
				}
			}()
			got := NewComparer(tt.args.oconf, tt.args.valueNames, tt.args.subNames...)
			assert.Equal(t, tt.want, got, tt.name)
		}()
	}
}
