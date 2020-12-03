package redis

import (
	"fmt"
	"strings"

	"github.com/micro-plat/hydra/registry"
)

func (r *Redis) GetValue(path string) (data []byte, version int32, err error) {
	buff, err := r.client.Get(swapKey(path)).Result()
	if err != nil {
		if err.Error() == "redis: nil" {
			return nil, 0, fmt.Errorf("数据不存在")
		}
		return nil, 0, err
	}
	value, err := newValueByJSON(buff)
	if err != nil {
		return nil, 0, err
	}
	return []byte(value.Data), value.Version, nil
}
func (r *Redis) GetChildren(path string) (paths []string, version int32, err error) {
	key := swapKey(path)
	npaths, err := r.client.Keys(key + ":*").Result()
	if err != nil {
		return nil, 0, err
	}

	exclude := swapKey(path, "watch")
	paths = make([]string, 0, len(npaths))
	for _, p := range npaths {
		if strings.HasPrefix(p, exclude) {
			continue
		}
		rpath := registry.Trim(swapPath(strings.TrimPrefix(p, key)))
		paths = append(paths, rpath)
	}
	return paths, 0, nil
}

func swapKey(elem ...string) string {
	var builder strings.Builder
	for _, v := range elem {
		if v == "/" || v == "\\" || strings.TrimSpace(v) == "" {
			continue
		}
		builder.WriteString(strings.Trim(v, "/"))
		builder.WriteString(":")
	}

	str := strings.ReplaceAll(builder.String(), "/", ":")
	return strings.TrimSuffix(str, ":")
}
func splitKey(key string) []string {
	return strings.Split(key, ":")
}

func swapPath(elem ...string) string {
	var builder strings.Builder
	for _, v := range elem {
		if v == "/" || v == "\\" || strings.TrimSpace(v) == "" {
			continue
		}
		builder.WriteString(strings.Trim(v, "/"))
		builder.WriteString("/")
	}

	str := strings.ReplaceAll(builder.String(), ":", "/")
	return registry.Format(str)
}
