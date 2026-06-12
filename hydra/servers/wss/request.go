package wss

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/url"
	"strings"
)

type dispatchRequest struct {
	method string
	path   string
	query  string
	form   map[string]interface{}
	header map[string]string
}

func newDispatchRequest(method string, path string, query string, body []byte, header map[string]string) *dispatchRequest {
	form := make(map[string]interface{})
	if query != "" {
		fillValues(form, query)
	}
	if len(body) > 0 {
		fillBody(form, body, header)
	}
	if header == nil {
		header = make(map[string]string)
	}
	return &dispatchRequest{method: method, path: path, query: query, form: form, header: header}
}

func fillValues(form map[string]interface{}, raw string) {
	values, err := url.ParseQuery(raw)
	if err != nil {
		return
	}
	for k, v := range values {
		if len(v) > 0 {
			form[k] = v[0]
		}
	}
}

func fillBody(form map[string]interface{}, body []byte, header map[string]string) {
	contentType := strings.ToLower(headerValue(header, "Content-Type"))
	if strings.Contains(contentType, "application/x-www-form-urlencoded") {
		fillValues(form, string(body))
	} else if strings.Contains(contentType, "multipart/form-data") {
		req, err := http.NewRequest(http.MethodPost, "/", bytes.NewReader(body))
		if err == nil {
			req.Header.Set("Content-Type", headerValue(header, "Content-Type"))
			if req.ParseMultipartForm(32<<20) == nil {
				for k, v := range req.PostForm {
					if len(v) > 0 {
						form[k] = v[0]
					}
				}
			}
		}
	} else {
		var m map[string]interface{}
		if json.Unmarshal(body, &m) == nil {
			for k, v := range m {
				form[k] = v
			}
		}
	}
	form["__body__"] = body
}

func headerValue(header map[string]string, key string) string {
	for k, v := range header {
		if strings.EqualFold(k, key) {
			return v
		}
	}
	return ""
}

func (r *dispatchRequest) GetName() string {
	return r.path
}

func (r *dispatchRequest) GetService() string {
	return r.path
}

func (r *dispatchRequest) GetQuery() string {
	return r.query
}

func (r *dispatchRequest) GetMethod() string {
	return r.method
}

func (r *dispatchRequest) GetForm() map[string]interface{} {
	return r.form
}

func (r *dispatchRequest) GetHeader() map[string]string {
	return r.header
}
