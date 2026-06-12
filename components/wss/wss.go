package wss

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"sync"
)

type Header map[string]string

type Request struct {
	Method string
	Path   string
	Query  string
	Body   []byte
	Header Header
}

type Response struct {
	Status int
	Header Header
	Body   []byte
}

type StartHandler func(resp *Response) error
type ChunkHandler func(data []byte) error

type Caller func(method string, path string, query string, body []byte, header map[string]string) (int, map[string]string, []byte, error)
type StreamCaller func(req Request, onChunk ChunkHandler) (*Response, error)
type OpenStreamCaller func(req Request, onStart StartHandler, onChunk ChunkHandler) (*Response, error)

type IComponentWSS interface {
	Request(req Request) (*Response, error)
	Call(method string, path string, query string, body []byte, header map[string]string) (int, map[string]string, []byte, error)

	Get(path string, query string, header ...map[string]string) (int, map[string]string, []byte, error)
	Post(path string, query string, body []byte, header ...map[string]string) (int, map[string]string, []byte, error)
	Put(path string, query string, body []byte, header ...map[string]string) (int, map[string]string, []byte, error)
	Delete(path string, query string, body []byte, header ...map[string]string) (int, map[string]string, []byte, error)
	Patch(path string, query string, body []byte, header ...map[string]string) (int, map[string]string, []byte, error)

	GET(path string, query ...map[string]string) (*Response, error)
	DELETE(path string, query ...map[string]string) (*Response, error)
	POST(path string, body []byte) (*Response, error)
	PUT(path string, body []byte) (*Response, error)
	PATCH(path string, body []byte) (*Response, error)
	PostJSON(path string, data interface{}) (*Response, error)
	PutJSON(path string, data interface{}) (*Response, error)
	PatchJSON(path string, data interface{}) (*Response, error)
	PostForm(path string, data map[string]string) (*Response, error)
	PutForm(path string, data map[string]string) (*Response, error)
	PatchForm(path string, data map[string]string) (*Response, error)
	PostText(path string, text string) (*Response, error)
	OpenStream(req Request, onStart StartHandler, onChunk ChunkHandler) (*Response, error)
	Stream(req Request, onChunk ChunkHandler) (*Response, error)
}

type ComponentWSS struct{}

var (
	lock             sync.RWMutex
	caller           Caller
	streamCaller     StreamCaller
	openStreamCaller OpenStreamCaller
)

func New() *ComponentWSS {
	return &ComponentWSS{}
}

func (r *Response) OK() bool {
	return r != nil && r.Status >= http.StatusOK && r.Status < http.StatusBadRequest
}

func (r *Response) String() string {
	if r == nil {
		return ""
	}
	return string(r.Body)
}

func (r *Response) JSON(v interface{}) error {
	if r == nil {
		return fmt.Errorf("wss response is nil")
	}
	return json.Unmarshal(r.Body, v)
}

func Register(c Caller) {
	lock.Lock()
	defer lock.Unlock()
	caller = c
}

func RegisterStream(c StreamCaller) {
	lock.Lock()
	defer lock.Unlock()
	streamCaller = c
}

func RegisterOpenStream(c OpenStreamCaller) {
	lock.Lock()
	defer lock.Unlock()
	openStreamCaller = c
}

func (c *ComponentWSS) Request(req Request) (*Response, error) {
	if req.Method == "" {
		req.Method = http.MethodGet
	}
	status, header, body, err := c.Call(req.Method, req.Path, req.Query, req.Body, map[string]string(req.Header))
	if header == nil {
		header = map[string]string{}
	}
	return &Response{Status: status, Header: Header(header), Body: body}, err
}

func (c *ComponentWSS) Call(method string, path string, query string, body []byte, header map[string]string) (int, map[string]string, []byte, error) {
	lock.RLock()
	current := caller
	lock.RUnlock()
	if current == nil {
		return http.StatusServiceUnavailable, nil, nil, fmt.Errorf("wss.server is not started")
	}
	return current(method, path, query, body, header)
}

func (c *ComponentWSS) Get(path string, query string, header ...map[string]string) (int, map[string]string, []byte, error) {
	return c.Call(http.MethodGet, path, query, nil, firstHeader(header...))
}

func (c *ComponentWSS) Post(path string, query string, body []byte, header ...map[string]string) (int, map[string]string, []byte, error) {
	return c.Call(http.MethodPost, path, query, body, firstHeader(header...))
}

func (c *ComponentWSS) Put(path string, query string, body []byte, header ...map[string]string) (int, map[string]string, []byte, error) {
	return c.Call(http.MethodPut, path, query, body, firstHeader(header...))
}

func (c *ComponentWSS) Delete(path string, query string, body []byte, header ...map[string]string) (int, map[string]string, []byte, error) {
	return c.Call(http.MethodDelete, path, query, body, firstHeader(header...))
}

func (c *ComponentWSS) Patch(path string, query string, body []byte, header ...map[string]string) (int, map[string]string, []byte, error) {
	return c.Call(http.MethodPatch, path, query, body, firstHeader(header...))
}

func (c *ComponentWSS) GET(path string, query ...map[string]string) (*Response, error) {
	return c.Request(Request{Method: http.MethodGet, Path: path, Query: firstQuery(query...)})
}

func (c *ComponentWSS) DELETE(path string, query ...map[string]string) (*Response, error) {
	return c.Request(Request{Method: http.MethodDelete, Path: path, Query: firstQuery(query...)})
}

func (c *ComponentWSS) POST(path string, body []byte) (*Response, error) {
	return c.Request(Request{Method: http.MethodPost, Path: path, Body: body})
}

func (c *ComponentWSS) PUT(path string, body []byte) (*Response, error) {
	return c.Request(Request{Method: http.MethodPut, Path: path, Body: body})
}

func (c *ComponentWSS) PATCH(path string, body []byte) (*Response, error) {
	return c.Request(Request{Method: http.MethodPatch, Path: path, Body: body})
}

func (c *ComponentWSS) PostJSON(path string, data interface{}) (*Response, error) {
	return c.doJSON(http.MethodPost, path, data)
}

func (c *ComponentWSS) PutJSON(path string, data interface{}) (*Response, error) {
	return c.doJSON(http.MethodPut, path, data)
}

func (c *ComponentWSS) PatchJSON(path string, data interface{}) (*Response, error) {
	return c.doJSON(http.MethodPatch, path, data)
}

func (c *ComponentWSS) PostForm(path string, data map[string]string) (*Response, error) {
	return c.doForm(http.MethodPost, path, data)
}

func (c *ComponentWSS) PutForm(path string, data map[string]string) (*Response, error) {
	return c.doForm(http.MethodPut, path, data)
}

func (c *ComponentWSS) PatchForm(path string, data map[string]string) (*Response, error) {
	return c.doForm(http.MethodPatch, path, data)
}

func (c *ComponentWSS) PostText(path string, text string) (*Response, error) {
	return c.Request(Request{
		Method: http.MethodPost,
		Path:   path,
		Body:   []byte(text),
		Header: Header{"Content-Type": "text/plain; charset=utf-8"},
	})
}

func (c *ComponentWSS) Stream(req Request, onChunk ChunkHandler) (*Response, error) {
	return c.OpenStream(req, nil, onChunk)
}

func (c *ComponentWSS) OpenStream(req Request, onStart StartHandler, onChunk ChunkHandler) (*Response, error) {
	lock.RLock()
	current := openStreamCaller
	legacy := streamCaller
	lock.RUnlock()
	if current == nil && legacy == nil {
		return &Response{Status: http.StatusServiceUnavailable}, fmt.Errorf("wss.server is not started")
	}
	if req.Method == "" {
		req.Method = http.MethodGet
	}
	if req.Header == nil {
		req.Header = Header{}
	}
	req.Header["Accept"] = "text/event-stream"
	if current != nil {
		return current(req, onStart, onChunk)
	}
	resp, err := legacy(req, onChunk)
	if onStart != nil && resp != nil {
		if startErr := onStart(resp); startErr != nil {
			return resp, startErr
		}
	}
	return resp, err
}

func (c *ComponentWSS) doForm(method string, path string, data map[string]string) (*Response, error) {
	body := []byte(firstQuery(data))
	return c.Request(Request{
		Method: method,
		Path:   path,
		Body:   body,
		Header: Header{"Content-Type": "application/x-www-form-urlencoded"},
	})
}

func (c *ComponentWSS) doJSON(method string, path string, data interface{}) (*Response, error) {
	body, err := json.Marshal(data)
	if err != nil {
		return nil, err
	}
	return c.Request(Request{
		Method: method,
		Path:   path,
		Body:   body,
		Header: Header{"Content-Type": "application/json; charset=utf-8"},
	})
}

func firstHeader(header ...map[string]string) map[string]string {
	if len(header) == 0 {
		return nil
	}
	return header[0]
}

func firstQuery(query ...map[string]string) string {
	if len(query) == 0 || query[0] == nil {
		return ""
	}
	values := url.Values{}
	for k, v := range query[0] {
		values.Set(k, v)
	}
	return values.Encode()
}
