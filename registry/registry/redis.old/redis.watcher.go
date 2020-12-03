package redis

import (
	"fmt"
	"strings"
	"time"

	"github.com/micro-plat/lib4go/registry"
	"github.com/micro-plat/lib4go/security/md5"
)

type valueEntity struct {
	Value   []byte
	version int32
	path    string
	Err     error
}

type valuesEntity struct {
	values  []string
	version int32
	path    string
	Err     error
}

func (v *valueEntity) GetPath() string {
	return v.path
}
func (v *valueEntity) GetValue() ([]byte, int32) {
	return v.Value, v.version
}
func (v *valueEntity) GetError() error {
	return v.Err
}

func (v *valuesEntity) GetValue() ([]string, int32) {
	return v.values, v.version
}
func (v *valuesEntity) GetError() error {
	return v.Err
}
func (v *valuesEntity) GetPath() string {
	return v.path
}

func (r *redisRegistry) WatchChildren(path string) (ch chan registry.ChildrenWatcher, err error) {
	if !r.isConnect {
		err = ErrColientCouldNotConnect
		return
	}

	ch = make(chan registry.ChildrenWatcher, 1)
	changeCh := make(chan valuesEntity, 1)
	errCh := make(chan error, 1)
	for {
		select {
		case <-r.CloseCh:
			ch <- &valuesEntity{path: path, Err: ErrClientConnClosing}
			return
		case e := <-errCh:
			ch <- &valuesEntity{path: path, Err: e}
			return
		case arry := <-changeCh:
			ch <- &arry
			return
		case <-time.After(time.Second * 3):
			if r.done {
				ch <- &valuesEntity{path: path, Err: ErrClientConnClosing}
				return
			}
			go func(ch chan valuesEntity, eCh chan error) {
				ps, vers, e := r.GetChildren(path)
				if e != nil {
					eCh <- e
					return
				}
				cuurtC := map[string]interface{}{}
				r.watchMap.Range(func(key, val interface{}) bool {
					rpath := key.(string)
					if strings.HasPrefix(rpath, path) && len(rpath) > len(path) {
						cuurtC[rpath] = val
					}
					return true
				})

				b := false
				for _, str := range ps {
					d, _, e := r.GetValue(str)
					if e != nil {
						eCh <- e
						return
					}
					val := md5.Encrypt(string(d))
					if _, ok := cuurtC[str]; ok {
						if val != fmt.Sprintf("%v", cuurtC[str]) {
							//修改了字节点的值
							r.watchMap.Store(str, val)
							b = true
						}
						delete(cuurtC, str)
					} else {
						//添加了字节点
						r.watchMap.Store(str, val)
						b = true
					}
				}
				if cuurtC != nil && len(cuurtC) > 0 {
					//删除了字节点
					for k, _ := range cuurtC {
						r.watchMap.Delete(k)
					}
					b = true
				}

				if b {
					ch <- valuesEntity{path: path, Err: err, values: ps, version: vers}
					return
				}
			}(changeCh, errCh)
		}
	}
}
func (r *redisRegistry) WatchValue(path string) (ch chan registry.ValueWatcher, err error) {
	if !r.isConnect {
		err = ErrColientCouldNotConnect
		return
	}
	if r.done {
		err = ErrClientConnClosing
		return
	}
	ch = make(chan registry.ValueWatcher, 1)
	changeCh := make(chan valueEntity, 1)
	errCh := make(chan error, 1)
	for {
		select {
		case <-r.CloseCh:
			ch <- &valueEntity{path: path, Err: ErrClientConnClosing}
			return
		case e := <-errCh:
			ch <- &valueEntity{path: path, Err: e}
			return
		case v := <-changeCh:
			ch <- &v
			return
		case <-time.After(time.Second * 3):
			// fmt.Println("yyyyyyyyyyyyyy:", path)
			go func(changeCh chan valueEntity, eCh chan error) {
				pathKey := joinR(path)
				res, vers, err := r.GetValue(pathKey)
				if err != nil && !strings.Contains(err.Error(), "不存在") {
					eCh <- err
					return
				}

				b := false
				rval, ok := r.watchMap.Load(path)
				if err == nil {
					val := md5.Encrypt(string(res))
					r.watchMap.Store(path, val)
					if !ok || val != fmt.Sprintf("%v", rval) {
						//新增或修改节点
						r.watchMap.Store(path, val)
						b = true
					}
				} else {
					if ok {
						//删除节点
						r.watchMap.Delete(path)
						b = true
					}
				}
				if b {
					// fmt.Println("zzzzzzzzzzz:", string(res))
					changeCh <- valueEntity{path: path, Err: nil, Value: res, version: vers}
					return
				}

			}(changeCh, errCh)
		}
	}
}

func (r *redisRegistry) eventWatch() {
	for {
		if r.done {
			return
		}
		select {
		case <-r.CloseCh:
			return
		case <-time.After(time.Second * 3):
			if _, err := r.client.Ping().Result(); err != nil {
				r.isConnect = false
				r.Log.Warnf("reids已断开连接:%v", r.options.Addrs)
			} else {
				r.isConnect = true
			}
		}
	}
}
