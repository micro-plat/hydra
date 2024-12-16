package dbr

import (
	"fmt"

	_ "github.com/go-sql-driver/mysql"
	_ "github.com/mattn/go-oci8"
	r "github.com/micro-plat/hydra/registry"
)

func getRegistryForTest(provider string) (r.IRegistry, error) {
	fact := &dbrFactory{proto: provider, opts: &r.Options{}}
	if provider == "mysql" {
		return fact.Create(r.WithAuthCreds("hbsv2x_dev", "123456dev"), r.Addrs("192.168.0.36"), r.Metadata("db", "hbsv2x_dev"))
	}
	if provider == "oracle" {
		return fact.Create(r.WithAuthCreds("ims17_v1_dev", "123456dev"), r.Addrs("orcl136"))
	}

	return nil, fmt.Errorf("获取注册中心provider错误,%s", provider)
}
