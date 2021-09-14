package river

import (
	"time"

	"github.com/micro-plat/hydra/conf/vars/db"
)

// Config is the configuration
type Config struct {
	DBConf        *db.DB        `json:"db_conf"`
	FlushBulkTime time.Duration `json:"flush_bulk_time"`
}
