package container

import (
	"fmt"
	"testing"
)

func TestVerAdd(t *testing.T) {
	verObj := newVers()
	tp := "test"
	name := "taosy"
	key := "201314"
	verObj.Add(tp, name, key)
	tps := fmt.Sprintf("%s-%s", tp, name)
	if _, ok := verObj.keys[tps]; !ok {
		t.Error("保存Ver失败")
		return
	}

	if len(verObj.keys) != 1 || verObj.keys[tps].current != key || len(verObj.keys[tps].keys) != 1 {
		t.Error("保存的ver数量不正确")
		return
	}

	key = "201313"
	verObj.Add(tp, name, key)
	if len(verObj.keys) != 1 || verObj.keys[tps].current != key || len(verObj.keys[tps].keys) != 2 {
		t.Error("保存的ver数量不正确11")
		return
	}

	tp = "test1"
	name = "taosy1"
	key = "201312"
	verObj.Add(tp, name, key)
	if len(verObj.keys) != 2 {
		t.Error("保存的ver数量不正确11")
		return
	}

	for _, item := range verObj.keys {
		if len(item.keys) <= 0 {
			t.Error("保存的ver数量不正确11")
			return
		}

		if len(item.keys) == 1 && item.current != key {
			t.Error("保存的ver数量不正确11")
			continue
		}

		res := false
		for _, str := range item.keys {
			if str == item.current {
				res = true
				continue
			}
		}
		if !res {
			t.Error("保存的ver数量不正确12")
			return
		}
	}
}

func TestVerRemove(t *testing.T) {
	verObj := newVers()
	tp := "test"
	name := "taosy"
	key := "201314"
	verObj.Add(tp, name, key)
	key = "201313"
	verObj.Add(tp, name, key)
	tps := fmt.Sprintf("%s-%s", tp, name)
	t.Logf("add key:%s,len:%v", verObj.keys[tps].current, verObj.keys[tps].keys)
	verObj.Remove(func(key string) bool {
		t.Log("移除key")
		return true
	})
	if _, ok := verObj.keys[tps]; !ok {
		t.Error("移除key失败")
		return
	}

	if len(verObj.keys[tps].keys) != 1 {
		t.Error("移除key失败1")
		return
	}

	if verObj.keys[tps].keys[0] != "201313" {
		t.Error("移除key失败2")
		return
	}

	if verObj.keys[tps].keys[0] != verObj.keys[tps].current {
		t.Error("移除key失败3")
		return
	}

	return
}
