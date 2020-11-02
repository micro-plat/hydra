package cache

//TypeNodeName 分类节点名
const TypeNodeName = "cache"

//Cache 消息队列配置
type Cache struct {
	Proto string `json:"proto"  toml:"proto" valid:"required" `
	Raw   []byte `json:"-"`
}
