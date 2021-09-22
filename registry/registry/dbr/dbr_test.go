/*
 * @Description:
 * @Autor: taoshouyin
 * @Date: 2021-09-18 09:36:32
 * @LastEditors: taoshouyin
 * @LastEditTime: 2021-09-18 10:59:27
 */
package dbr

import (
	"testing"

	_ "github.com/go-sql-driver/mysql"
	r "github.com/micro-plat/hydra/registry"
	"github.com/micro-plat/lib4go/assert"
)

func getRegistry(t *testing.T) r.IRegistry {
	fact := &dbrFactory{proto: MYSQL, opts: &r.Options{}}
	r, err := fact.Create(r.WithAuthCreds("hbsv2x_dev", "123456dev"), r.Addrs("192.168.0.36"), r.Metadata("db", "hbsv2x_dev"))
	assert.Equal(t, nil, err, err)
	return r
}

func Test_dbrFactory_Create(t *testing.T) {
	tests := []struct {
		name    string
		proto   string
		opts    []r.Option
		want    bool
		wantErr bool
	}{
		{name: "mysql注册中心初始化-错误的地址", proto: "mysql", opts: []r.Option{r.Addrs("125215"), r.WithAuthCreds("hbsv2x_dev", "123456dev"), r.WithMetadata(map[string]string{"db": "hbsv2x_dev"})}, wantErr: true, want: false},
		{name: "mysql注册中心初始化-错误帐号密码", proto: "mysql", opts: []r.Option{r.Addrs("192.168.0.36:3306"), r.WithAuthCreds("hbsv2x_dev", "11111"), r.WithMetadata(map[string]string{"db": "hbsv2x_dev"})}, wantErr: true, want: false},
		{name: "mysql注册中心初始化-不存在的库名", proto: "mysql", opts: []r.Option{r.Addrs("192.168.0.36:3306"), r.WithAuthCreds("hbsv2x_dev", "123456dev"), r.WithMetadata(map[string]string{"db": "xxx"})}, wantErr: true, want: false},
		{name: "mysql注册中心初始化-正确链接", proto: "mysql", opts: []r.Option{r.Addrs("192.168.0.36:3306"), r.WithAuthCreds("hbsv2x_dev", "123456dev"), r.WithMetadata(map[string]string{"db": "hbsv2x_dev"})}, wantErr: false, want: true},
		{name: "mysql注册中心初始化-错误的注册枚举", proto: "xxx", opts: []r.Option{r.Addrs("192.168.0.36:3306"), r.WithAuthCreds("hbsv2x_dev", "123456dev"), r.WithMetadata(map[string]string{"db": "hbsv2x_dev"})}, wantErr: true, want: false},

		{name: "oracle注册中心初始化-错误的地址", proto: "oracle", opts: []r.Option{r.Addrs("125215"), r.WithAuthCreds("hbsv2x_dev", "123456dev")}, wantErr: true, want: false},
		{name: "oracle注册中心初始化-错误帐号密码", proto: "oracle", opts: []r.Option{r.Addrs("orcl136"), r.WithAuthCreds("hbsv2x_dev", "11111")}, wantErr: true, want: false},
		{name: "oracle注册中心初始化-不存在的库名", proto: "oracle", opts: []r.Option{r.Addrs("orcl136"), r.WithAuthCreds("hbsv2x_dev", "123456dev")}, wantErr: true, want: false},
		{name: "oracle注册中心初始化-正确链接", proto: "oracle", opts: []r.Option{r.Addrs("orcl136"), r.WithAuthCreds("hbsv2x_dev", "123456dev")}, wantErr: false, want: true},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			z := &dbrFactory{
				proto: tt.proto,
				opts:  &r.Options{},
			}
			got, err := z.Create(tt.opts...)
			assert.Equal(t, tt.wantErr, (err != nil))
			assert.Equal(t, tt.want, (got != nil))
		})
	}
}

func TestCreatePersistentNode(t *testing.T) {
	r := getRegistry(t)
	err := r.CreatePersistentNode("/node", `{"id":100}`)
	assert.Equal(t, nil, err, err)

	buff, ver, err := r.GetValue("/node")
	assert.Equal(t, nil, err, err)
	assert.Equal(t, `{"id":100}`, string(buff))
	assert.Equal(t, int32(1), ver)

}
func TestCreateTempNode(t *testing.T) {
	r := getRegistry(t)
	err := r.CreateTempNode("/node/t", `{"id":100}`)
	assert.Equal(t, nil, err, err)

	buff, ver, err := r.GetValue("/node/t")
	assert.Equal(t, nil, err, err)
	assert.Equal(t, `{"id":100}`, string(buff))
	assert.Equal(t, int32(1), ver)

}
