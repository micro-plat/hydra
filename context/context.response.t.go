package context

import (
	"strings"
)

const (
	CT_DEF = iota + 1
	CT_JSON
	CT_XML
	CT_YMAL
	CT_PLAIN
	CT_HTML
)

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

//GetContentType  0：自动 1:json 2:xml 3:plain
func (r *Response) GetContentType() int {
	cc, ok1 := r.Content.(string)
	tp, ok2 := r.Params["Content-Type"].(string)
	if ok2 {
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
	if ok1 {
		if strings.HasPrefix(cc, "[") && strings.HasSuffix(cc, "]") {
			return CT_JSON
		} else if strings.HasPrefix(cc, "{") && strings.HasSuffix(cc, "}") {
			return CT_JSON
		} else if strings.HasPrefix(cc, "<?xml") {
			return CT_XML
		} else if strings.HasPrefix(cc, "<html") {
			return CT_HTML
		}
		return CT_PLAIN
	}
	return CT_DEF
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
