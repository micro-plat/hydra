package mock

import (
	"testing"

	"github.com/micro-plat/hydra"
	"github.com/stretchr/testify/assert"
)

type ts struct {
	Name string `json:"name"`
}

func TestGetObject(t *testing.T) {
	hydra.Conf.Vars().Custom("a", "b", map[string]interface{}{
		"name": "colin",
	})

	app, err := NewAPPConf()
	assert.Equal(t, nil, err)
	assert.Equal(t, nil, err)
	v := &ts{}
	app.GetVarConf().GetObject("a", "b", v)
	assert.Equal(t, "colin", v.Name)

	c := &ts{}
	app.GetVarConf().GetObject("a", "b", c)
	assert.Equal(t, "colin", c.Name)

}
func BenchmarkGetObject(t *testing.B) {
	hydra.Conf.Vars().Custom("a", "b", map[string]interface{}{
		"name": "colin",
	})

	app, err := NewAPPConf()
	assert.Equal(t, nil, err)
	assert.Equal(t, nil, err)

	for i := 0; i < t.N; i++ {
		v := &ts{}
		app.GetVarConf().GetObject("a", "b", v)
		assert.Equal(t, "colin", v.Name)
	}

}
