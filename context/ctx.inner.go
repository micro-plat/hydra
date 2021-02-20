package context

import (
	"io"
	"net/http"
	"net/url"
)

type IInnerContext interface {
	ClientIP() string
	GetBody() io.ReadCloser     //.Request.Body
	GetMethod() string          //.Request.Method
	GetURL() *url.URL           //.Request.URL.Path
	Header(string, string)      //context.Header
	GetHeaders() http.Header    //Request.Header
	GetCookies() []*http.Cookie //Request.Cookies()
	GetParams() map[string]interface{}
	GetRouterPath() string //Context.FullPath()
	GetPostForm() url.Values
	GetRawForm() map[string]interface{}
	ContentType() string

	Abort()
	WStatus(int)           //c.Context.Writer.WriteHeader(s)
	Status() int           //Context.Writer.Status()
	Written() bool         //Context.Writer.Written()
	WHeaders() http.Header //c.Context.Writer.Header()
	WHeader(string) string //c.Context.Writer.Header().Get
	// File(string)           //Context.File(path)
	ServeContent(filepath string, fs http.FileSystem) int
	Data(int, string, []byte) //c.Context.Data(status, tpName, v)
	Redirect(int, string)
	GetService() string
	GetFile(fileKey string) (string, io.ReadCloser, int64, error)
	GetHTTPReqResp() (*http.Request, http.ResponseWriter)
	ClearAuth(c ...bool) bool
}
