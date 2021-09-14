package river

import (
	"sync"

	"github.com/micro-plat/hydra/components/dbs"
	"github.com/micro-plat/hydra/registry/registry/mysql/internal/client"
	"github.com/micro-plat/lib4go/db"
	"github.com/micro-plat/lib4go/types"

	log "github.com/lib4dev/cli/logger"
)

type River struct {
	c *Config

	wg sync.WaitGroup

	client *client.Client

	closeCh chan struct{}

	db dbs.IDB

	sourceData map[string]types.XMap
}

func NewRiver(c *Config) (*River, error) {
	r := new(River)

	r.c = c
	r.closeCh = make(chan struct{})
	var err error
	r.db, err = db.NewDB(c.DBConf.Provider, c.DBConf.ConnString, c.DBConf.MaxOpen, c.DBConf.MaxIdle, c.DBConf.LifeTime)
	if err != nil {
		return nil, err
	}

	r.client = client.NewClient()

	return r, nil
}

func (r *River) GetClient() *client.Client {
	return r.client
}

func (r *River) Run() error {
	r.wg.Add(1)
	go r.syncLoop()

	return nil
}

// Close closes the River
func (r *River) Close() {
	log.Log.Infof("closing river")

	r.wg.Wait()
	close(r.closeCh)
}
