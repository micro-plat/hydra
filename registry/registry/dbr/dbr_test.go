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
