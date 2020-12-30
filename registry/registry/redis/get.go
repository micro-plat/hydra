package redis

import (
	"fmt"
	"strings"

	"github.com/micro-plat/hydra/registry/registry/redis/internal"
)

//GetValue 获取节点值
func (r *Redis) GetValue(path string) (data []byte, version int32, err error) {
	buff, err := r.client.Get(internal.SwapKey(path)).Result()
	if err != nil {
		if err.Error() == "redis: nil" {
			existsChild, err := r.client.ExistsChildren(internal.SwapKey(path) + ":*")
			if !existsChild || err != nil {
				return nil, 0, fmt.Errorf("数据不存在")
			}
			return []byte{}, 0, nil
		}
		return nil, 0, err
	}
	value, err := newValueByJSON(buff)
	if err != nil {
		return nil, 0, err
	}
	return []byte(value.Data), value.Version, nil
}

//GetChildren 获取所有子节点
func (r *Redis) GetChildren(path string) (paths []string, version int32, err error) {
	key := internal.SwapKey(path)

	//npaths, err := r.client.Keys(key + ":*").Result()
	npaths, err := r.client.SearchChildren(key + ":*")
	if err != nil {
		return nil, 0, err
	}

	exclude := internal.SwapKey(path, "watch")
	paths = make([]string, 0, len(npaths))
	cache := map[string]bool{}

	for _, p := range npaths {

		if strings.HasPrefix(p, exclude) {
			continue
		}

		p = strings.TrimPrefix(p, key+":")
		if idx := strings.Index(p, ":"); idx > 0 {
			p = p[:idx]
		}
		if p == "" {
			continue
		}

		if ok, _ := cache[p]; ok {
			continue
		}
		cache[p] = true
		paths = append(paths, p)
	}

	return paths, 0, nil
}
