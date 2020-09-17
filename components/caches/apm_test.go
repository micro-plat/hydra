package caches

import (
	"testing"
)

type MockCache struct {
}

func (t *MockCache) Get(key string) (string, error) {
	return "mockval", nil
}
func (t *MockCache) Decrement(key string, delta int64) (n int64, err error) {
	return 1, nil
}
func (t *MockCache) Increment(key string, delta int64) (n int64, err error) {
	return 1, nil
}
func (t *MockCache) Gets(key ...string) (r []string, err error) {
	return []string{"a", "b"}, nil
}
func (t *MockCache) Add(key string, value string, expiresAt int) error {
	return nil
}
func (t *MockCache) Set(key string, value string, expiresAt int) error {
	return nil
}
func (t *MockCache) Delete(key string) error {
	return nil
}
func (t *MockCache) Exists(key string) bool {
	return true
}
func (t *MockCache) Delay(key string, expiresAt int) error {
	return nil
}
func (t *MockCache) Close() error {
	return nil
}

type APMMockCache struct {
	*MockCache
}

func (t *APMMockCache) GetProto() string {
	return "mock"
}

func (t *APMMockCache) GetServers() []string {
	return []string{"ip1", "ip2"}
}

func TestNewCache(t *testing.T) {

	org := &MockCache{}

	r, err := NewAPMCache("test", org)
	if err != nil {
		t.Error("构建失败")
	}

	if r != org {
		t.Error("构建逻辑错误")
	}

	v, err := r.Get("t")
	if v != "mockval" || err != nil {
		t.Error("Caches Get 测试失败", err)
	}
	dv, err := r.Decrement("t", 1)
	if dv != 1 || err != nil {
		t.Error("Caches Decrement 测试失败", err)
	}
	iv, err := r.Increment("t", 1)
	if iv != 1 || err != nil {
		t.Error("Caches Increment 测试失败", err)
	}

	getvs, err := r.Gets("a", "b")
	if len(getvs) != 2 || err != nil {
		t.Error("Caches Gets 测试失败", err)
	}

	err = r.Add("t", "v", 1)
	if err != nil {
		t.Error("Caches Add 构建失败")
	}

	err = r.Set("t", "v", 1)
	if err != nil {
		t.Error("Caches Set 构建失败")
	}

	err = r.Delete("t")
	if err != nil {
		t.Error("Caches Delete 构建失败")
	}

	b := r.Exists("t")
	if !b {
		t.Error("Caches Exists 构建失败")
	}

	err = r.Delay("t", 1)
	if err != nil {
		t.Error("Caches Delay 构建失败")
	}

	err = r.Close()
	if err != nil {
		t.Error("Caches Close 构建失败")
	}

}

func TestNewAPMCache(t *testing.T) {
	org := &APMMockCache{}
	apmr, err := NewAPMCache("testapm", org)
	if err != nil {
		t.Error("构建失败")
	}

	if apmr == org {
		t.Error("apm构建逻辑错误")
	}

	v, err := apmr.Get("t")
	if v != "mockval" || err != nil {
		t.Error("Caches Get 测试失败", err)
	}
	dv, err := apmr.Decrement("t", 1)
	if dv != 1 || err != nil {
		t.Error("Caches Decrement 测试失败", err)
	}
	iv, err := apmr.Increment("t", 1)
	if iv != 1 || err != nil {
		t.Error("Caches Increment 测试失败", err)
	}

	getvs, err := apmr.Gets("a", "b")
	if len(getvs) != 2 || err != nil {
		t.Error("Caches Gets 测试失败", err)
	}

	err = apmr.Add("t", "v", 1)
	if err != nil {
		t.Error("Caches Add 构建失败")
	}

	err = apmr.Set("t", "v", 1)
	if err != nil {
		t.Error("Caches Set 构建失败")
	}

	err = apmr.Delete("t")
	if err != nil {
		t.Error("Caches Delete 构建失败")
	}

	b := apmr.Exists("t")
	if !b {
		t.Error("Caches Exists 构建失败")
	}

	err = apmr.Delay("t", 1)
	if err != nil {
		t.Error("Caches Delay 构建失败")
	}

	err = apmr.Close()
	if err != nil {
		t.Error("Caches Close 构建失败")
	}

}
