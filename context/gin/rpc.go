package gin

type rpc struct {
	user *user
}

func (r *rpc) Request(service string, form map[string]interface{}, opts ...rpc.RequestOption) (res *rpc.Response, err error) {
	// if header == nil {
	// 	header = map[string]string{}
	// }
	// if input == nil {
	// 	input = map[string]interface{}{}
	// }
	// if _, ok := header["X-Request-Id"]; !ok {
	// 	header["X-Request-Id"] = r.user.requestID
	// }

	// method, ok := header["method"]
	// if !ok {
	// 	method = "get"
	// }
	// status, r, param, err = cr.rpc.Request(service, strings.ToUpper(method), header, form, failFast)
	// if err != nil || status != 200 {
	// 	return
	// }
	return nil, nil
}
