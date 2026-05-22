// Package testing 提供 Hydra 组件测试辅助函数
package testing

import (
	"sync"

	"github.com/micro-plat/hydra/components/container"
	"github.com/micro-plat/lib4go/db"
)

var (
	// globalTestRegistry 全局测试组件注册表
	globalTestRegistry = make(map[string]interface{})
	registryMutex       sync.RWMutex
)

// init 在包导入时初始化测试组件获取函数
func init() {
	container.SetTestComponentGetter(getTestComponent)
}

// getTestComponent 获取测试组件（内部使用）
func getTestComponent(typ, name string) (interface{}, bool) {
	registryMutex.RLock()
	defer registryMutex.RUnlock()

	key := typ + "_" + name
	component, ok := globalTestRegistry[key]
	return component, ok
}

// SetTestDB 设置测试用 Mock 数据库
// 会注册到全局测试注册表，供所有 Container 使用
//
// 示例:
//   mockDB := db.NewMockDB()
//   SetTestDB(mockDB)
//   defer ClearTestDB()
func SetTestDB(mockDB db.IDB, names ...string) {
	name := "db"
	if len(names) > 0 {
		name = names[0]
	}

	registryMutex.Lock()
	defer registryMutex.Unlock()

	globalTestRegistry["db_"+name] = mockDB
}

// ClearTestDB 清除指定名称的测试数据库
func ClearTestDB(names ...string) {
	name := "db"
	if len(names) > 0 {
		name = names[0]
	}

	registryMutex.Lock()
	defer registryMutex.Unlock()

	key := "db_" + name
	if component, ok := globalTestRegistry[key]; ok {
		// 关闭可关闭的组件
		if closer, ok := component.(interface{ Close() error }); ok {
			_ = closer.Close()
		}
		delete(globalTestRegistry, key)
	}
}

// ClearAllTestComponents 清除所有测试组件
func ClearAllTestComponents() {
	registryMutex.Lock()
	defer registryMutex.Unlock()

	// 关闭可关闭的组件
	for key, v := range globalTestRegistry {
		if closer, ok := v.(interface{ Close() error }); ok {
			_ = closer.Close()
		}
		delete(globalTestRegistry, key)
	}
}

// ========== 内部辅助函数 ==========

// setTestComponent 注册测试组件（内部使用）
func setTestComponent(typ, name string, component interface{}) {
	registryMutex.Lock()
	defer registryMutex.Unlock()

	key := typ + "_" + name
	globalTestRegistry[key] = component
}

// clearTestComponent 清除测试组件（内部使用）
func clearTestComponent(typ, name string) {
	registryMutex.Lock()
	defer registryMutex.Unlock()

	key := typ + "_" + name
	if component, ok := globalTestRegistry[key]; ok {
		// 关闭可关闭的组件
		if closer, ok := component.(interface{ Close() error }); ok {
			_ = closer.Close()
		}
		delete(globalTestRegistry, key)
	}
}

// hasTestComponents 检查是否有测试组件注册（内部使用）
func hasTestComponents() bool {
	registryMutex.RLock()
	defer registryMutex.RUnlock()

	return len(globalTestRegistry) > 0
}
