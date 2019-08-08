package conf

//WebServerConf web服务器配置
type WebServerConf = APIServerConf

//NewWebServerConf 构建api server配置信息
func NewWebServerConf(address string) *WebServerConf {
	return &WebServerConf{
		Address: address,
	}
}
