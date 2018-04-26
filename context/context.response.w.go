package context

import (
	"bytes"
	"fmt"
	"time"
)

//SetView 设置view
func (r *Response) SetView(name string) {
	r.Params["__view"] = name
}

//NoView 设置view
func (r *Response) NoView() {
	r.Params["__view"] = "NONE"
}

//SetCookie 设置cookie
func (r *Response) SetCookie(name string, value string, timeout int, domain string) {
	list := make([]string, 0, 2)
	if v, ok := r.Params["Set-Cookie"]; ok {
		list = v.([]string)
	}
	list = append(list, r.getSetCookie(name, value, timeout, domain))
	r.Params["Set-Cookie"] = list
}
func (r *Response) getSetCookie(name string, value string, timeout interface{}, domain string) string {
	var b bytes.Buffer
	fmt.Fprintf(&b, "%s=%s", name, value)
	var maxAge int64
	switch v := timeout.(type) {
	case int:
		maxAge = int64(v)
	case int32:
		maxAge = int64(v)
	case int64:
		maxAge = v
	}
	switch {
	case maxAge > 0:
		if len(domain) > 0 {
			fmt.Fprintf(&b, "; Expires=%s; Max-Age=%d;path=/;domain=%s", time.Now().Add(time.Duration(maxAge)*time.Second).UTC().Format(time.RFC1123), maxAge, domain)
			return b.String()
		}
		fmt.Fprintf(&b, "; Expires=%s; Max-Age=%d;path=/", time.Now().Add(time.Duration(maxAge)*time.Second).UTC().Format(time.RFC1123), maxAge)

	case maxAge < 0:
		fmt.Fprintf(&b, "; Max-Age=0")
	}
	return b.String()
}

func (r *Response) Success(v interface{}) {
	r.Status = 200
	r.Content = v
	return

}
