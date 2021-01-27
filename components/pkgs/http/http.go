package http

import (
	"bytes"
	"fmt"
	"io"
	"mime/multipart"
	"net"
	"net/http"
	"os"

	"time"

	varhttp "github.com/micro-plat/hydra/conf/vars/http"
)

//Client HTTP客户端
type Client struct {
	*varhttp.HTTPConf
	client *http.Client
}

//ClientRequest  http请求
type ClientRequest struct {
	headers  map[string]string
	client   *http.Client
	method   string
	url      string
	params   string
	encoding string
}

// NewClient 构建HTTP客户端，用于发送GET POST等请求
func NewClient(opts ...varhttp.Option) (client *Client, err error) {
	httpconf := varhttp.New(opts...)
	return NewClientByConf(httpconf)
}

//NewClientByConf 通过配置对象获取客户端
func NewClientByConf(conf *varhttp.HTTPConf) (client *Client, err error) {
	client = &Client{}
	client.HTTPConf = conf
	tlsConf, err := getCert(client.HTTPConf)
	if err != nil {
		return nil, err
	}
	orginalClient := &http.Client{
		Transport: &http.Transport{
			DisableKeepAlives: client.HTTPConf.Keepalive,
			TLSClientConfig:   tlsConf,
			Proxy:             getProxy(client.HTTPConf),
			Dial: func(netw, addr string) (net.Conn, error) {
				c, err := net.DialTimeout(netw, addr, time.Second*time.Duration(client.HTTPConf.ConnectionTimeout))
				if err != nil {
					return nil, err
				}
				c.SetDeadline(time.Now().Add(time.Duration(client.HTTPConf.RequestTimeout) * time.Second))
				return c, nil
			},
			MaxIdleConnsPerHost:   0,
			ResponseHeaderTimeout: 0,
		},
	}
	client.client = orginalClient
	return
}

// Get http get请求
func (c *Client) Get(url string, charset ...string) (content string, status int, err error) {
	ncharset := getCharset(charset...)
	r, s, err := c.Request(http.MethodGet, url, "", ncharset, http.Header{})
	return string(r), s, err
}

// Post http Post请求
func (c *Client) Post(url string, params string, charset ...string) (content string, status int, err error) {
	ncharset := getCharset(charset...)
	r, s, err := c.Request(http.MethodPost, url, params, ncharset, http.Header{})
	return string(r), s, err
}

//Upload 文件上传
func (c *Client) Upload(url string, params map[string]string, files map[string]string, charset string, header http.Header, cookies ...*http.Cookie) (content string, status int, err error) {
	ncharset := getCharset(charset)
	bodyBuffer := &bytes.Buffer{}
	bodyWriter := multipart.NewWriter(bodyBuffer)

	//字段处理
	for k, v := range params {
		err = bodyWriter.WriteField(k, v)
		if err != nil {
			return "", 0, fmt.Errorf("设置字段失败:%s(%v)", k, v)
		}
	}

	//文件流处理
	for k, v := range files {
		fw1, err := bodyWriter.CreateFormFile(k, v)
		if err != nil {
			return "", 0, fmt.Errorf("无法创建文件流:%v", v)
		}
		f1, err := os.Open(v)
		if err != nil {
			return "", 0, fmt.Errorf("无法读取文件:%s", v)
		}
		defer f1.Close()
		io.Copy(fw1, f1)
	}

	contentType := bodyWriter.FormDataContentType()
	bodyWriter.Close()
	header.Set("Content-Type", contentType)
	r, s, err := c.Request(http.MethodPost, url, bodyBuffer.String(), ncharset, header, cookies...)
	return string(r), s, err
}

// SaveAs 将请求内容以文件方式保存
func (c *Client) SaveAs(method string, url string, params string, path string, charset string, header http.Header, cookies ...*http.Cookie) (status int, err error) {
	body, status, err := c.Request(method, url, params, charset, header, cookies...)
	if err != nil {
		return
	}
	fl, err := os.OpenFile(path, os.O_RDWR|os.O_APPEND|os.O_CREATE, 0664)
	if err != nil {
		return
	}
	defer fl.Close()
	n, err := fl.Write(body)
	if err == nil && n < len(body) {
		err = io.ErrShortWrite
	}
	return
}
