package creator

import "github.com/micro-plat/lib4go/errs"

type iCustomerBuilder interface {
	Load()
	ISUB
	Map() map[string]interface{}
}

var _ iCustomerBuilder = &CustomerBuilder{}

type CustomerBuilder struct {
	httpBuilder
	loader func() (string, interface{})
}

//newCustomerBuilder 构建http生成器
func newCustomerBuilder(s ...interface{}) *CustomerBuilder {
	b := &CustomerBuilder{httpBuilder: httpBuilder{BaseBuilder: BaseBuilder{}}}
	if len(s) == 0 {
		b.httpBuilder.BaseBuilder[ServerMainNodeName] = make(map[string]interface{})
		return b
	}
	b.httpBuilder.BaseBuilder[ServerMainNodeName] = s[0]
	return b
}

//AddLoader 添加加载函数
func (b *CustomerBuilder) AddLoader(loader func() (string, interface{})) {
	b.loader = loader
}

//Load 加载配置
func (b *CustomerBuilder) Load() {
	if b.loader == nil {
		return
	}
	node, data := b.loader()
	if err := errs.GetError(data); err != nil {
		panic(err)
	}
	b.BaseBuilder[node] = data
}
