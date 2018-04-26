package context

//Redirect 设置页面转跳
func (r *Response) Redirect(code int, url string) {
	r.Params["Status"] = code
	r.Params["Location"] = url
	r.Status = code
	return
}

//IsRedirect 是否是URL转跳
func (r *Response) IsRedirect() (string, bool) {
	location, ok := r.Params["Location"]
	if !ok {
		return "", false
	}
	url, ok := location.(string)
	if !ok {
		return url, false
	}
	status := r.Params["Status"]
	return url, status == 301 || status == 302 || status == 303 || status == 307 || status == 309
}
