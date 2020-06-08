package http

import "net/http"

//IClient http请求
type IClient interface {
	Get(url string, charset ...string) (content string, status int, err error)
	Post(url string, params string, charset ...string) (content string, status int, err error)
	Request(method string, url string, params string, charset string, header http.Header, cookies ...*http.Cookie) (content []byte, status int, err error)
	SaveAs(method string, url string, params string, path string, charset string, header http.Header, cookies ...*http.Cookie) (status int, err error)
	Upload(url string, params map[string]string, files map[string]string, charset string, header http.Header, cookies ...*http.Cookie) (content string, status int, err error)
}

//IComponentHTTPClient http请求组件
type IComponentHTTPClient interface {
	GetRegularClient(names ...string) (d IClient)
	GetClient(names ...string) (d IClient, err error)
}
