package http

import (
	"compress/gzip"
	"crypto/rand"
	"crypto/tls"
	"crypto/x509"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"strings"

	varhttp "github.com/micro-plat/hydra/conf/vars/http"
	"github.com/micro-plat/hydra/context"
	"github.com/micro-plat/lib4go/encoding"
)

func (c *Client) Request(method string, url string, params string, charset string, header http.Header, cookies ...*http.Cookie) (content []byte, status int, err error) {
	r, _, s, err := c.HRequest(method, url, params, charset, header, cookies...)
	return r, s, err
}

// Request 发送http请求, method:http请求方法包括:get,post,delete,put等 url: 请求的HTTP地址,不包括参数,params:请求参数,
// header,http请求头多个用/n分隔,每个键值之前用=号连接
func (c *Client) HRequest(method string, url string, params string, charset string, header http.Header, cookies ...*http.Cookie) (content []byte, rspHeader http.Header, status int, err error) {
	method = strings.ToUpper(method)
	req, err := http.NewRequest(method, url, encoding.GetEncodeReader([]byte(params), charset))
	if err != nil {
		return
	}

	for _, cookie := range cookies {
		req.AddCookie(cookie)
	}
	req.Close = true
	if header == nil {
		header = http.Header{}
	}
	if c := header.Get("Content-Type"); (method == "POST" || method == "PUT" || method == "DELETE") && c == "" {
		header.Set("Content-Type", fmt.Sprintf("application/x-www-form-urlencoded;charset=%s", charset))
	}
	for i, v := range header {
		req.Header.Set(i, strings.Join(v, ","))
	}

	if ctx, ok := context.GetContext(); ok {
		req.Header.Set(context.XRequestID, ctx.User().GetTraceID())
	}
	response, err := c.client.Do(req)
	if response != nil {
		defer response.Body.Close()
	}
	if err != nil {
		return nil, nil, 0, fmt.Errorf("client.Do err:%v", err)
	}

	rawBody := response.Body
	if strings.EqualFold(response.Header.Get("Content-Encoding"), "gzip") {
		rawBody, err = gzip.NewReader(response.Body)
		if err != nil {
			err = fmt.Errorf("http resp unzip is failed,err: %w", err)
			return nil, nil, status, err
		}
	}

	body, err := ioutil.ReadAll(rawBody)
	if err != nil {
		return nil, nil, 0, fmt.Errorf("body ReadAll err:%v", err)
	}
	rspHeader = response.Header
	status = response.StatusCode
	ct, err := encoding.DecodeBytes(body, charset)
	if err != nil {
		return nil, nil, 0, fmt.Errorf("body charset err:%v", err)
	}
	content = ct
	return
}

func getCert(c *varhttp.HTTPConf) (*tls.Config, error) {
	ssl := &tls.Config{InsecureSkipVerify: true}
	if len(c.Certs) == 2 {
		cert, err := tls.LoadX509KeyPair(c.Certs[0], c.Certs[1])
		if err != nil {
			return nil, fmt.Errorf("cert证书(pem:%s,key:%s),加载失败:%v", c.Certs[0], c.Certs[1], err)
		}
		ssl.Certificates = []tls.Certificate{cert}
	}
	if c.Ca != "" {
		caData, err := ioutil.ReadFile(c.Ca)
		if err != nil {
			return nil, fmt.Errorf("ca证书(%s)读取错误:%v", c.Ca, err)
		}
		pool := x509.NewCertPool()
		pool.AppendCertsFromPEM(caData)
		ssl.RootCAs = pool
	}
	if len(ssl.Certificates) == 0 && ssl.RootCAs == nil {
		return ssl, nil
	}
	ssl.Rand = rand.Reader
	return ssl, nil

}
func getProxy(c *varhttp.HTTPConf) func(*http.Request) (*url.URL, error) {
	if c.Proxy != "" {
		return func(_ *http.Request) (*url.URL, error) {
			return url.Parse(c.Proxy) //根据定义Proxy func(*Request) (*url.URL, error)这里要返回url.URL
		}
	}
	return nil
}
func getCharset(charset ...string) (encoding string) {
	if len(charset) > 0 {
		return strings.ToUpper(charset[0])
	}
	return "UTF-8"
}

func (c *Client) print(p func(...interface{}), h string, r ...interface{}) {
	if c.Trace {
		line := make([]interface{}, 0, len(r)+1)
		line = append(line, h)
		line = append(line, r...)
		p(line...)
	}
}
