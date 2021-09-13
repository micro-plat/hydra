package river

import (
	"reflect"
	"time"

	"github.com/micro-plat/hydra/registry/registry/mysql/internal/client"
	"github.com/micro-plat/hydra/registry/registry/mysql/internal/sql"
	"github.com/micro-plat/lib4go/types"

	log "github.com/lib4dev/cli/logger"
)

func (r *River) syncLoop() {

	interval := r.c.FlushBulkTime
	if interval == 0 {
		interval = 1000 * time.Millisecond
	}

	ticker := time.NewTicker(interval)
	defer ticker.Stop()
	defer r.wg.Done()

	reqs := make([]*client.BulkRequest, 0, 1024)

	for {
		needFlush := false
		select {
		case <-ticker.C:
			list, err := r.db.Query(sql.GetData, nil)
			if err != nil {
				log.Log.Errorf("pull data err %v, close sync", err)
				return
			}
			diff := r.diff(list)
			if len(diff) > 0 {
				reqs = append(reqs, diff...)
				needFlush = true
			}
			r.sourceData = getDataMap(list)

		case <-r.closeCh:
			return
		}

		if needFlush {
			if err := r.doBulk(reqs); err != nil {
				return
			}
			reqs = reqs[0:0]
		}
	}
}

func (r *River) doBulk(reqs []*client.BulkRequest) error {
	if len(reqs) == 0 {
		return nil
	}
	if err := r.client.Bulk(reqs); err != nil {
		return err
	}

	return nil
}

func getDataMap(data types.XMaps) map[string]types.XMap {
	r := map[string]types.XMap{}
	for _, v := range data {
		r[v.GetString("path")] = v
	}

	return r
}

func (r *River) diff(t types.XMaps) []*client.BulkRequest {
	tempSource := r.sourceData
	target := getDataMap(t)
	diff := make([]*client.BulkRequest, 0)

	//新增
	for name, data := range target {
		if _, ok := tempSource[name]; !ok {
			item := &client.BulkRequest{
				Action:  client.ActionInsert,
				PkValue: data.GetString("path"),
			}
			diff = append(diff, item)
		}
	}

	//减少
	for name, data := range tempSource {
		if _, ok := target[name]; !ok {
			item := &client.BulkRequest{
				Action:  client.ActionDelete,
				PkValue: data.GetString("path"),
			}
			diff = append(diff, item)
			delete(tempSource, name)
		}
	}

	//变动
	for name, data := range tempSource {
		if !reflect.DeepEqual(data, target[name]) {
			item := &client.BulkRequest{
				Action:  client.ActionUpdate,
				PkValue: data.GetString("path"),
			}
			diff = append(diff, item)
		}
	}

	return diff
}
