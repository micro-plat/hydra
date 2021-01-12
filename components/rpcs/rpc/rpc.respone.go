package rpc

// import (
// 	"encoding/json"
// 	"fmt"
// 	"net/http"
// 	"strings"

// 	"github.com/micro-plat/lib4go/jsons"
// )

// //------------------------RPC响应---------------------------------------

// //Response 请求结果
// type Response struct {
// 	Status int
// 	Header string
// 	hdMap  map[string]string
// 	Result string
// }

// //NewResponse 请求响应
// func NewResponse(status int, header string, result string) (res *Response) {
// 	res = &Response{
// 		Status: status,
// 		Result: result,
// 		Header: header,
// 	}
// 	res.hdMap, _ = getHeader(header)
// 	return res
// }

// //NewResponseByStatus 根据状态构建响应
// func NewResponseByStatus(status int, err error) *Response {
// 	r := NewResponse(status, "{}", err.Error())
// 	return r
// }

// //IsSuccess 请求是否成功
// func (r *Response) IsSuccess() bool {
// 	return r.Status == http.StatusOK
// }

// //IsJSON 结果是否是json串
// func (r *Response) IsJSON() bool {
// 	ctp := r.GetHeader("Content-Type")
// 	buff := []byte(r.Result)
// 	if (ctp == "" || strings.Contains(ctp, "json")) && json.Valid(buff) && (strings.HasPrefix(r.Result, "{") ||
// 		strings.HasPrefix(r.Result, "[")) {
// 		return true
// 	}
// 	return false
// }

// //GetResult 获取请求结果
// func (r *Response) GetResult() (map[string]interface{}, error) {
// 	return jsons.Unmarshal([]byte(r.Result))
// }

// //GetHeader 根据KEY获取参数
// func (r *Response) GetHeader(key string) string {
// 	return r.hdMap[key]
// }

// func getHeader(h string) (map[string]string, error) {
// 	hd := make(map[string]string)
// 	if h == "" {
// 		return hd, nil
// 	}
// 	mh, err := jsons.Unmarshal([]byte(h))
// 	if err != nil {
// 		return nil, err
// 	}
// 	for k, v := range mh {
// 		switch t := v.(type) {
// 		case []string:
// 			hd[k] = strings.Join(t, ",")
// 		case string:
// 			hd[k] = t
// 		default:
// 			hd[k] = fmt.Sprint(t)
// 		}
// 	}
// 	return hd, nil
// }
