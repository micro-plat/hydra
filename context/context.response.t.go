package context

import (
	"bytes"
	"encoding/json"
	"encoding/xml"
	"fmt"
	"reflect"
	"strings"

	yaml "gopkg.in/yaml.v2"
)

const (
	CT_DEF = iota
	CT_JSON
	CT_XML
	CT_YMAL
	CT_HTML
	CT_PLAIN
)

//ContentTypes content type map
var ContentTypes = map[int]string{
	CT_JSON:  "application/json; charset=UTF-8",
	CT_XML:   "text/xml; charset=UTF-8",
	CT_YMAL:  "text/ymal; charset=UTF-8",
	CT_HTML:  "text/html; charset=UTF-8",
	CT_PLAIN: "text/plain; charset=UTF-8",
}

//SeTextJSON 将content type设置为application/json; charset=UTF-8
func (r *Response) SeTextJSON() {
	r.Params["Content-Type"] = "application/json; charset=UTF-8"
}

//SetTextXML 将content type设置为text/xml; charset=UTF-8
func (r *Response) SetTextXML() {
	r.Params["Content-Type"] = "text/xml; charset=UTF-8"
}

//SetTextHTML 将content type设置为text/html; charset=UTF-8
func (r *Response) SetTextHTML() {
	r.Params["Content-Type"] = "text/html; charset=UTF-8"
}

//SetTextPlain 将content type设置为text/plain; charset=UTF-8
func (r *Response) SetTextPlain() {
	r.Params["Content-Type"] = "text/plain; charset=UTF-8"
}

//GetRenderContent  0：自动 1:json 2:xml 3:plain
func (r *Response) getContentType() int {
	tp, ok := r.Params["Content-Type"].(string)
	if ok {
		if strings.Contains(tp, "json") {
			return CT_JSON
		} else if strings.Contains(tp, "xml") {
			return CT_XML
		} else if strings.Contains(tp, "plain") {
			return CT_PLAIN
		} else if strings.Contains(tp, "yaml") {
			return CT_YMAL
		}
	}
	return CT_DEF
}

//GetRenderContent 获取用于render的content type 和 内容
func (r *Response) GetRenderContent(df int) (int, interface{}, error) {
	data := r.GetContent()
	t := r.getContentType()
	if data == nil {
		return t, nil, nil
	}
	val := reflect.ValueOf(data)
	if val.Kind() == reflect.Interface || val.Kind() == reflect.Ptr {
		val = val.Elem()
	}
	switch val.Kind() {
	case reflect.Struct, reflect.Slice, reflect.Map, reflect.Array:
		switch {
		case t == CT_JSON || (t == CT_DEF && df == CT_JSON):
			return CT_JSON, data, nil
		case t == CT_XML || (t == CT_DEF && df == CT_XML):
			buff, err := xml.Marshal(data)
			if err != nil {
				return t, nil, err
			}
			return CT_XML, buff, nil
		case t == CT_YMAL || (t == CT_DEF && df == CT_YMAL):
			buff, err := yaml.Marshal(data)
			if err != nil {
				return 0, nil, err
			}
			return CT_YMAL, buff, nil
		case t == CT_DEF:
			return df, data, nil
		default: // CT_PLAIN, CT_HTML:
			return t, fmt.Sprintf("%+v", data), nil
		}
	case reflect.String:
		value := []byte(data.(string))
		switch {
		case (t == CT_JSON || t == CT_DEF) && json.Valid(value):
			return CT_JSON, data, nil
		case (t == CT_XML || t == CT_DEF) && bytes.HasPrefix(value, []byte("<?xml")):
			return CT_XML, value, nil
		case (t == CT_HTML || t == CT_DEF) && bytes.HasPrefix(value, []byte("<!DOCTYPE html")):
			return CT_HTML, data, nil
		}
		switch {
		case t == CT_JSON || t == CT_DEF && df == CT_JSON:
			return CT_JSON, map[string]interface{}{"data": data}, nil
		case t == CT_XML || t == CT_DEF && df == CT_XML:
			return CT_XML, data, nil
		case t == CT_DEF:
			return df, data, nil
		default:
			return t, data, nil
		}

	default:
		switch {
		case t == CT_PLAIN || t == CT_HTML:
			return t, data, nil
		case t == CT_YMAL || t == CT_DEF && df == CT_YMAL:
			buff, err := yaml.Marshal(map[string]interface{}{
				"data": data,
			})
			if err != nil {
				return t, nil, err
			}
			return t, buff, nil
		default:
			return df, map[string]interface{}{"data": data}, nil
		}
	}
}

func (r *Response) GetContent() interface{} {
	return r.Content
}

func (r *Response) ShouldContent(content interface{}) {
	switch v := content.(type) {
	case IError:
		r.err = v.GetError()
		r.Status = v.GetCode()
	case error:
		r.err = content.(error)
	}
	r.Status = r.getStatus(content)
	r.Content = content
	return
}

func (r *Response) MustContent(status int, content interface{}) {
	r.Status = status
	r.ShouldContent(content)
}
func (r *Response) getStatus(c interface{}) int {
	switch c.(type) {
	case IError, error:
		if r.Status < 400 {
			return 400
		}
		return r.Status
	default:
		if r.Status == 0 {
			r.Status = 200
		}
		return r.Status
	}
}

//JSON 按json格式输入
func (r *Response) JSON(content interface{}) {
	r.SeTextJSON()
	r.ShouldContent(content)
}

//XML 按xml格式输入
func (r *Response) XML(content interface{}) {
	r.SetTextXML()
	r.ShouldContent(content)
}

//Text 按text/plain格式输入
func (r *Response) Text(content interface{}) {
	r.SetTextPlain()
	r.ShouldContent(content)
}

//HTML 按text/HTML
func (r *Response) HTML(content interface{}) {
	r.SetTextHTML()
	r.ShouldContent(content)
}
