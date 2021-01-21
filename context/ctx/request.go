package ctx

import (
	"fmt"
	"net/http"
	"strings"

	"github.com/asaskevich/govalidator"
	"github.com/micro-plat/hydra/conf"
	"github.com/micro-plat/hydra/conf/app"
	"github.com/micro-plat/hydra/context"
	"github.com/micro-plat/lib4go/errs"
	"github.com/micro-plat/lib4go/types"
)

type request struct {
	ctx     context.IInnerContext
	appConf app.IAPPConf
	cookies types.XMap
	headers types.XMap
	types.XMap
	readMapErr error
	*body
	path *rpath
	*file
}

//NewRequest 构建请求的Request
//自动对请求进行解码，响应结果进行编码。
//当指定为gbk,gb2312后,请求方式为application/x-www-form-urlencoded或application/xml、application/json时内容必须编码为指定的格式，否则会解码失败
func NewRequest(c context.IInnerContext, s app.IAPPConf, meta conf.IMeta) *request {
	rpath := NewRpath(c, s, meta)
	req := &request{
		ctx:  c,
		body: NewBody(c, rpath.GetEncoding()),
		XMap: make(map[string]interface{}),
		path: rpath,
		file: NewFile(c, meta),
	}
	req.XMap, req.readMapErr = req.body.GetMap()
	if req.XMap == nil {
		req.XMap = make(map[string]interface{})
	}
	if req.readMapErr != nil {
		req.readMapErr = errs.NewError(http.StatusNotAcceptable, req.readMapErr)
	}
	return req
}

//GetHTTPRequest 获取http request原生对象
func (r *request) GetHTTPRequest() *http.Request {
	req, _ := r.ctx.GetHTTPReqResp()
	return req
}

//Path 获取请求路径信息
func (r *request) Path() context.IPath {
	return r.path
}

//Bind 根据输入参数绑定对象，并使用govalidator.ValidateMap进行参数验证
func (r *request) Bind(obj interface{}) error {
	if r.readMapErr != nil {
		return r.readMapErr
	}

	//处理数据结构转换
	if err := r.XMap.ToAnyStruct(obj); err != nil {
		return errs.NewError(http.StatusNotAcceptable, fmt.Errorf("对象转换有误 %v", err))
	}

	//验证数据格式
	if _, err := govalidator.ValidateStruct(obj); err != nil {
		return errs.NewError(http.StatusNotAcceptable, fmt.Errorf("输入参数有误 %v", err))
	}
	return nil
}

//Check 检查输入参数和配置参数是否为空
func (r *request) Check(field ...string) error {
	for _, key := range field {
		if v := r.GetString(key); v == "" {
			return errs.NewError(http.StatusNotAcceptable, fmt.Errorf("输入参数:%s值不能为空", key))
		}
	}
	return nil
}

//CheckMap 传入验证Map[字段名]验证规则，并使用govalidator.ValidateMap进行参数验证
func (r *request) CheckMap(vdt map[string]interface{}) error {
	if len(vdt) == 0 {
		return nil
	}
	input := make(map[string]interface{})
	for k := range vdt {
		input[k] = r.XMap[k]
	}

	//验证数据格式
	if _, err := govalidator.ValidateMap(input, vdt); err != nil {
		return errs.NewError(http.StatusNotAcceptable, fmt.Errorf("输入参数有误 %v", err))
	}
	return nil
}

//GetMap 获取请求的参数信息
func (r *request) GetMap() types.XMap {
	return r.XMap
}
func (r *request) GetError() error {
	return r.readMapErr
}

//GetPlayload 获取trace信息
func (r *request) GetPlayload() string {
	if r.readMapErr != nil {
		return fmt.Errorf("err:%w", r.readMapErr).Error()
	}
	body := r.GetMap()
	return fmt.Sprintf("%+v", body)
}

//Headers 获取请求的header
func (r *request) Headers() types.XMap {
	if r.headers != nil {
		return r.headers
	}
	hds := r.ctx.GetHeaders()
	r.headers = make(map[string]interface{})
	for k, v := range hds {
		r.headers[k] = strings.Join(v, ",")
	}
	return r.headers
}

//Cookies 获取请求的header
func (r *request) Cookies() types.XMap {
	if r.cookies != nil {
		return r.cookies
	}
	r.cookies = make(map[string]interface{})
	cookies := r.ctx.GetCookies()
	for _, cookie := range cookies {
		r.cookies[cookie.Name] = cookie.Value
	}
	return r.cookies
}
