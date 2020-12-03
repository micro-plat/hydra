package redis

import (
	"fmt"
)

func (r *Redis) Update(path string, data string) (err error) {

	//获取原数据
	b, err := r.Exists(path)
	if err != nil {
		return fmt.Errorf("查询出错:%w", err)
	}
	if !b {
		return fmt.Errorf("节点不存在%s", swapKey(path))
	}

	//获取原数据
	buff, err := r.client.Get(swapKey(path)).Result()
	if err != nil {
		return err
	}

	//解析并判断节点类型
	ovalue, err := newValueByJSON(buff)
	if err != nil {
		return err
	}

	//超时时长
	exp := r.maxExpiration
	if ovalue.IsTemp {
		exp = r.tmpExpiration
	}

	//构建新对象，并修改
	value := newValue(data, ovalue.IsTemp)
	_, err = r.client.Set(swapKey(path), value.String(), exp).Result() //? timeout
	if err != nil {
		return err
	}

	//通知变更
	r.notifyValueChange(path, value)
	return nil

}

//Delete 删除节点
func (r *Redis) Delete(path string) error {
	key := swapKey(path)
	keys := []string{key, swapKey(key, "seq")}
	for _, k := range keys {
		_, err := r.client.Del(k).Result()
		if err != nil {
			return fmt.Errorf("%v(%s)", err, k)
		}
	}
	r.notifyParentChange(path, 0)
	return nil
}
