// Package testing 提供 Hydra 容器测试辅助函数
package testing

// RegisterTestComponent 注册测试组件
// 用于注册任意类型的测试组件到容器
//
// 示例:
//   mockComponent := &MyMockComponent{}
//   RegisterTestComponent("cache", "redis", mockComponent)
func RegisterTestComponent(typ, name string, component interface{}) {
	setTestComponent(typ, name, component)
}

// UnregisterTestComponent 取消注册测试组件
func UnregisterTestComponent(typ, name string) {
	clearTestComponent(typ, name)
}

// HasTestComponents 检查是否有测试组件注册
func HasTestComponents() bool {
	return hasTestComponents()
}
