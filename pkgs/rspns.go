package pkgs

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"strings"

	"github.com/asaskevich/govalidator"
	"github.com/micro-plat/lib4go/encoding"
	"github.com/micro-plat/lib4go/errs"
	"github.com/micro-plat/lib4go/types"
	"gopkg.in/yaml.v2"
)

//------------------------RPC响应---------------------------------------

//Rspns 请求响应
type Rspns struct {
	status   int
	header   types.XMap
	data     types.XMap
	err      error
	result   interface{}
	encoding string
}

//NewRspns 根据状态构建响应
func NewRspns(result interface{}) *Rspns {
	r := NewRspnsByHD(0, "{}", result)
	return r
}

//NewRspnsByHD 请求响应
func NewRspnsByHD(status int, header string, result interface{}) (r *Rspns) {
	r = &Rspns{
		encoding: "utf-8",
		result:   result,
		status:   types.DecodeInt(status, 0, http.StatusOK),
		header:   make(map[string]interface{}),
		data:     make(map[string]interface{}),
	}
	//处理请求头
	if header != "" {
		if err := json.Unmarshal([]byte(header), &r.header); err != nil {
			r.header["header"] = header
		}
	}
	//处理Content-Type
	if _, ok := r.header["Content-Type"]; !ok {
		r.header["Content-Type"] = "application/json"
	}
	switch v := result.(type) {
	case error:
		if r.status >= http.StatusOK && r.status < http.StatusBadRequest {
			r.status = http.StatusBadRequest
		}
		r.err = fmt.Errorf("请求发生错误：%w", v)
		r.data["__body__"] = v
	case errs.IError:
		if r.status >= http.StatusOK && r.status < http.StatusBadRequest {
			r.status = v.GetCode()
		}
		r.result = v.GetError().Error()
		r.data["__body__"] = v
	case string:
		//转换数据
		r.data, r.err = r.getMap(fmt.Sprint(r.header["Content-Type"]), types.StringToBytes(v))
	case map[string]interface{}:
		r.data = v
	default:
		panic("不支持的数据类型")
	}
	return r
}

//IsSuccess 请求是否成功(状态码是否为２００)
func (r *Rspns) IsSuccess() bool {
	return r.status == http.StatusOK
}

//GetStatus 获取响应的状态码
func (r *Rspns) GetStatus() int {
	return r.status
}

//GetMap 获取响应的map对象
func (r *Rspns) GetMap() types.XMap {
	return r.data
}

//GetHeaders 获取响应头信息
func (r *Rspns) GetHeaders() types.XMap {
	return r.header
}

//GetResult 获取响应的原串
func (r *Rspns) GetResult() interface{} {
	return r.result
}

//GetError 获取远程请求的error或result转换为map
func (r *Rspns) GetError() error {
	return r.err
}

//Bind 根据输入参数绑定对象
func (r *Rspns) Bind(obj interface{}) error {
	if r.err != nil {
		return r.err
	}
	//处理数据结构转换
	if err := r.data.ToAnyStruct(obj); err != nil {
		return errs.NewError(http.StatusNotAcceptable, fmt.Errorf("对象转换有误 %v", err))
	}

	//验证数据格式
	if _, err := govalidator.ValidateStruct(obj); err != nil {
		return errs.NewError(http.StatusNotAcceptable, fmt.Errorf("响应参数有误 %v", err))
	}
	return nil
}

func (r *Rspns) getMap(ctp string, body []byte) (data map[string]interface{}, err error) {
	data = make(map[string]interface{})
	switch {
	case strings.Contains(ctp, "/xml"):
		data, err = types.NewXMapByXML(types.BytesToString(body))
		if err != nil {
			return nil, fmt.Errorf("xml转换为map失败:%w", err)
		}
	case strings.Contains(ctp, "/yaml") || strings.Contains(ctp, "/x-yaml"):
		err = yaml.Unmarshal(body, &data)
	case strings.Contains(ctp, "/json"):
		newbuff := body
		if bytes.HasPrefix(newbuff, []byte(`"{\"`)) {
			var s string
			err = json.Unmarshal(body, &s)
			if err == nil {
				newbuff = types.StringToBytes(s)
			}
		}
		d := json.NewDecoder(bytes.NewReader(newbuff))
		d.UseNumber()
		err = d.Decode(&data)
		r.err = err
	case strings.Contains(ctp, "/x-www-form-urlencoded") || strings.Contains(ctp, "/form-data"):
		var values url.Values
		values, err = url.ParseQuery(types.BytesToString(body))
		if err != nil {
			break
		}
		for k, v := range values {
			//处理编码问题
			var buff []byte
			buff, err = encoding.Decode(strings.Join(v, ","), r.encoding)
			if err != nil {
				break
			}
			data[k] = types.BytesToString(buff)
		}
	default:
		data["__body__"] = types.BytesToString(body)
	}
	return data, nil
}
