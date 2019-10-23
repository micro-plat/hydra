package context

import (
	"strings"
)

const (
	CT_DEF = iota
	CT_JSON
	CT_XML
	CT_YMAL
	CT_HTML
	CT_PLAIN
	CT_OTHER
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
func (r *Response) SetJSON() {
	r.Params["Content-Type"] = "application/json; charset=UTF-8"
}
func (r *Response) SetContentType(p string) {
	r.Params["Content-Type"] = p
}

//SetTextXML 将content type设置为text/xml; charset=UTF-8
func (r *Response) SetXML() {
	r.Params["Content-Type"] = "text/xml; charset=UTF-8"
}

//SetTextHTML 将content type设置为text/html; charset=UTF-8
func (r *Response) SetHTML() {
	r.Params["Content-Type"] = "text/html; charset=UTF-8"
}

//SetTextPlain 将content type设置为text/plain; charset=UTF-8
func (r *Response) SetPlain() {
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
		} else if strings.Contains(tp, "html") {
			return CT_HTML
		}
		return CT_OTHER
	}
	return CT_DEF
}
func (r *Response) GetContent() interface{} {
	return r.Content
}

func (r *Response) ShouldContent(content interface{}) {
	switch v := content.(type) {
	case IResult:
		r.Status = v.GetCode()
		r.Content = v.GetResult()
		return
	case IError:
		r.err = v.GetError()
		r.Status = v.GetCode()
	case error:
		r.err = content.(error)
	}
	r.Status = r.GetCode(content)
	r.Content = content
	return
}

func (r *Response) MustContent(status int, content interface{}) {
	r.Status = status
	r.ShouldContent(content)
}
func (r *Response) GetCode(c interface{}) int {
	switch v := c.(type) {
	case IResult:
		return v.GetCode()
	case IError:
		return v.GetCode()
	case error:
		if r.Status == 0 {
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
	r.SetJSON()
	r.ShouldContent(content)
}

//XML 按xml格式输入
func (r *Response) XML(content interface{}) {
	r.SetXML()
	r.ShouldContent(content)
}

//Text 按text/plain格式输入
func (r *Response) Text(content interface{}) {
	r.SetPlain()
	r.ShouldContent(content)
}

//HTML 按text/HTML
func (r *Response) HTML(content interface{}) {
	r.SetHTML()
	r.ShouldContent(content)
}
