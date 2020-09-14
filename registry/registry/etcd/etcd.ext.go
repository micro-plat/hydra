package etcd

import (
	"context"
	"errors"
	"fmt"
	"strings"
	"time"

	hr "github.com/micro-plat/hydra/registry"
	"github.com/micro-plat/lib4go/registry"
	"github.com/zhiyunliu/etcd/clientv3"
	"github.com/zhiyunliu/etcd/etcdserver/api/v3rpc/rpctypes"
)

var (
	ErrColientCouldNotConnect = errors.New("etcd: could not connect to the server")
	ErrClientConnClosing      = errors.New("etcd: the client connection is closing")
)

const LEASE_TTL = 30

func (e *etcdRegistry) WatchChildren(path string) (data chan registry.ChildrenWatcher, err error) {

	//fmt.Println("WatchChildren:", path)
	data = make(chan registry.ChildrenWatcher, 1)

	watcher, _ := e.Watch(hr.WatchService(path), hr.WatchDomain(e.options.Domain))

	watchChan := make(chan *valuesEntity, 1)

	go func(w hr.Watcher, c chan *valuesEntity) {
		for {
			item, err := w.Next()

			//fmt.Println("WatchChildren.event:", path, item.Data)
			if err != nil {

				c <- &valuesEntity{path: path, Err: err}
			}
			switch item.Action {
			case "create":
				fallthrough
			case "update":
				fallthrough
			case "delete":
				paths, version, err := e.GetChildren(path)
				c <- &valuesEntity{path: path, Err: err, values: paths, version: version}
			}
		}

	}(watcher, watchChan)

	go func(ch chan registry.ChildrenWatcher) {
		select {
		case <-e.CloseCh:
			ch <- &valuesEntity{path: path, Err: ErrClientConnClosing}
			watcher.Stop()
			return
		case result := <-watchChan:
			ch <- result
		}
	}(data)

	return
}

func (e *etcdRegistry) WatchValue(path string) (data chan registry.ValueWatcher, err error) {
	//fmt.Println("WatchValue:", path)
	data = make(chan registry.ValueWatcher, 1)
	watcher, _ := e.Watch(hr.WatchService(path), hr.WatchDomain(e.options.Domain))

	watchChan := make(chan *valueEntity, 1)

	go func(w hr.Watcher, c chan *valueEntity) {
		for {
			item, err := w.Next()

			//fmt.Println("Watch.event:", path, item.Data)
			if err != nil {
				c <- &valueEntity{path: path, Err: err}
			}
			switch item.Action {
			case "create":
				fallthrough
			case "update":
				fallthrough
			case "delete":
				paths, version, err := e.GetValue(path)
				c <- &valueEntity{path: path, Err: err, Value: paths, version: version}
			}
		}

	}(watcher, watchChan)

	go func(data chan registry.ValueWatcher) {
		for {
			select {
			case <-e.CloseCh:
				data <- &valueEntity{path: path, Err: ErrClientConnClosing}
				return
			case result := <-watchChan:
				data <- result
			}
		}
	}(data)
	return
}
func (e *etcdRegistry) GetChildren(path string) (paths []string, version int32, err error) {

	//fmt.Println("etcdRegistry.GetChildren:", path)
	getchildrenOpts := []clientv3.OpOption{
		clientv3.WithSerializable(),
		clientv3.WithSort(clientv3.SortByKey, clientv3.SortAscend),
		clientv3.WithPrefix(),
		clientv3.WithKeysOnly(),
		clientv3.WithFragment(),
		clientv3.WithPrevKV(),
	}
	ctx, cancel := context.WithTimeout(context.Background(), e.options.Timeout)
	defer cancel()

	if !strings.HasSuffix(path, "/") {
		path = path + "/"
	}

	rsp, err := e.client.Get(ctx, path, getchildrenOpts...)
	if err != nil {
		return
	}
	if rsp.Count <= 0 {
		return
	}

	tmpMap := map[string]bool{}

	for _, v := range rsp.Kvs {
		//fmt.Println("rsp.Kvs:", i, string(v.Key))
		tmpv := strings.Replace(string(v.Key), path, "", -1)
		//fmt.Println("tmpv:", tmpv)
		if idx := strings.Index(tmpv, "/"); idx > 0 {
			tmpv = string(tmpv[0:idx])
		}
		//fmt.Println("tmpv.2:", tmpv)
		if len(tmpv) <= 0 {
			continue
		}
		tmpMap[tmpv] = true
	}

	for k := range tmpMap {
		paths = append(paths, k)
	}

	version = int32(rsp.Header.GetRevision())

	//fmt.Println("subPath:", paths)

	return
}
func (e *etcdRegistry) GetValue(path string) (data []byte, version int32, err error) {
	ctx, cancel := context.WithTimeout(context.Background(), e.options.Timeout)
	defer cancel()

	rsp, err := e.client.Get(ctx, path, clientv3.WithSerializable())
	if err != nil {
		return
	}
	version = int32(0)
	if len(rsp.Kvs) == 0 {
		return
	}
	kv := rsp.Kvs[0]
	data = kv.Value
	return
}
func (e *etcdRegistry) CreatePersistentNode(path string, data string) (err error) {

	//fmt.Println("etcdRegistry.CreatePersistentNode:", path, data)
	// create an entry for the node
	var putOpts []clientv3.OpOption
	if rsp, err := e.client.Put(context.Background(), path, data, putOpts...); err != nil {
		fmt.Println("etcdRegistry.client.Put.1:", rsp.Header, rsp.PrevKv)
		//fmt.Println("etcdRegistry.client.Put.2:", path, data, err)
		return err
	}

	return
}
func (e *etcdRegistry) CreateTempNode(path string, data string) (err error) {
	_, err = e.createNode(path, data)
	return
}
func (e *etcdRegistry) CreateSeqNode(path string, data string) (rpath string, err error) {
	return e.createNode(path, data)
}

func (e *etcdRegistry) Update(path string, data string) (err error) {

	// create an entry for the node
	putOpts := []clientv3.OpOption{}
	if _, err = e.client.Put(context.Background(), path, data, putOpts...); err != nil {
		return err
	}

	return
}
func (e *etcdRegistry) Delete(path string) (err error) {
	// create an entry for the node
	putOpts := []clientv3.OpOption{}
	if _, err = e.client.Delete(context.Background(), path, putOpts...); err != nil {
		//fmt.Println("etcdRegistry.Delete:", err)
		return err
	}
	return
}
func (e *etcdRegistry) Exists(path string) (exists bool, err error) {
	//fmt.Println("etcdRegistry.Exists：", path)
	// create an entry for the node
	putOpts := []clientv3.OpOption{}

	rsp, err := e.client.Get(context.Background(), path, putOpts...)
	if err != nil {
		return false, err
	}
	// fmt.Println(rsp.Header)
	// fmt.Println(rsp.Kvs)
	// fmt.Println(rsp.More)
	// fmt.Println(rsp.Count)
	if rsp.Count <= 0 {
		return false, nil
	}
	return true, nil
}
func (e *etcdRegistry) Close() error {
	close(e.CloseCh)
	return e.client.Close()
}

func (e *etcdRegistry) createNode(path, data string) (rpath string, err error) {
	var lgr *clientv3.LeaseGrantResponse

	//fmt.Println("etcdRegistry.CreateTempNode", path, data)
	leaseClt := clientv3.NewLease(e.client)
	// get a lease used to expire keys since we have a ttl
	lgr, err = leaseClt.Grant(context.Background(), int64(LEASE_TTL))
	if err != nil {
		//fmt.Println("etcdRegistry.client.Grant", path, data)
		return
	}

	rpath = fmt.Sprintf("%s%d", path, lgr.Revision)

	//fmt.Println("path.2:", rpath)
	// create an entry for the node
	var putOpts []clientv3.OpOption
	if lgr != nil {
		putOpts = append(putOpts, clientv3.WithLease(lgr.ID))
	}
	if _, err = e.client.Put(context.Background(), rpath, data, putOpts...); err != nil {
		return
	}
	e.leases.Store(rpath, lgr.ID)
	return
}

func (e *etcdRegistry) leaseRemain() {
	for {
		//fmt.Println("leaseRemain.", time.Now().Format("20060102150405"))
		select {
		case <-e.CloseCh:
			return
		case <-time.After(time.Second * (LEASE_TTL - 1)):
			e.leases.Range(func(key, val interface{}) bool {
				path := key.(string)
				//fmt.Println("leases.Range.", path, key)
				go e.leaseKeepAliveOnce(path)
				return true
			})
		}
	}
}

func (e *etcdRegistry) leaseKeepAliveOnce(path string) (err error) {
	// missing lease, check if the key exists
	ctx, cancel := context.WithTimeout(context.Background(), e.options.Timeout)
	defer cancel()

	rsp, err := e.client.Get(ctx, path, clientv3.WithSerializable())
	if err != nil {
		return err
	}
	var leaseID clientv3.LeaseID
	// get the existing lease
	for _, kv := range rsp.Kvs {
		if kv.Lease > 0 {
			leaseID = clientv3.LeaseID(kv.Lease)
			break
		}
	}

	if _, err := e.client.KeepAliveOnce(context.TODO(), leaseID); err != nil {
		if err != rpctypes.ErrLeaseNotFound {
			return err
		}
	}
	return
}
