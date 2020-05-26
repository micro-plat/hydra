package context

import (
	"io"
	"net/http"
	"net/url"
)

type IInnerContext interface {
	ClientIP() string
	GetBody() io.ReadCloser            //.Request.Body
	GetMethod() string                 //.Request.Method
	UrlQuery() url.Values              //url.Query
	GetURL() *url.URL                  //.Request.URL.Path
	Header(string, string)             //context.Header
	GetHeaders() http.Header           //Request.Header
	GetCookies() []*http.Cookie        //Request.Cookies()
	Param(string) string               //Context.Param(key)
	ShouldBind(interface{}) error      //Context.ShouldBind(&obj)
	GetPostForm(string) (string, bool) //context GetPostForm)
	GetQuery(string) (string, bool)    //context.GetQuery
	PostForm() url.Values              //Context.Request.PostForm

	Abort()
	WStatus(int)              //c.Context.Writer.WriteHeader(s)
	Status() int              //Context.Writer.Status()
	Written() bool            //Context.Writer.Written()
	WHeader(string) string    //c.Context.Writer.Header().Get
	File(string)              //Context.File(path)
	Data(int, string, []byte) //c.Context.Data(status, tpName, v)
	XML(int, interface{})
	YAML(int, interface{})
	JSON(int, interface{})
	Redirect(int, string)
}
