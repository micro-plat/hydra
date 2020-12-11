package redis

import "github.com/micro-plat/hydra/registry/registry/redis/internal"

//Exists 检查节点是否存在
func (r *Redis) Exists(path string) (bool, error) {
	key := internal.SwapKey(path)
	e, err := r.client.Exists(key).Result()
	if err != nil {
		return false, err
	}
	if err == nil && e == 1 {
		return true, nil
	}
	//npaths, err := r.client.Keys(key + ":*").Result()
	exists, err := r.client.ExistsChildren(key + ":*")
	if err != nil {
		return false, err
	}
	return exists, err
}
