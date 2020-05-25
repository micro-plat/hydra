package uuid

//IUUID UUID接口
type IUUID interface {
	GetString(pre ...string) string
	Get() int64
}
