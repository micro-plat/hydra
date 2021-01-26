package ctx

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/url"
	"reflect"
	"strings"

	"github.com/micro-plat/hydra/context"
	"github.com/micro-plat/lib4go/encoding"
	"github.com/micro-plat/lib4go/types"
	"gopkg.in/yaml.v3"
)

type params struct {
	body  []byte
	query string
}

type valueReader struct {
	hasRead bool
	value   interface{}
	err     error
}

//body 用于处理http请求的body读取
type body struct {
	ctx      context.IInnerContext
	encoding string
	rawBody  valueReader
	fullBody valueReader
	mapBody  valueReader
}

//NewBody 构建body处理工具
func NewBody(c context.IInnerContext, encoding string) *body {
	return &body{ctx: c, encoding: encoding}
}

//GetMap 读取并合并请求参数
func (w *body) GetMap() (data map[string]interface{}, err error) {

	//从缓存中读取数据
	if w.mapBody.hasRead {
		if w.mapBody.err != nil {
			return nil, w.mapBody.err
		}
		if m, ok := w.mapBody.value.(map[string]interface{}); ok {
			return m, nil
		}
		return nil, nil
	}

	//从body中读取原处理流
	w.mapBody.hasRead = true
	var body []byte
	ctp := strings.ToLower(w.ctx.ContentType())

	if strings.Contains(ctp, "__raw__") {
		return w.ctx.GetRawForm(), nil
	}
	if body, _, w.mapBody.err = w.GetFullRaw(); w.mapBody.err != nil {
		return nil, w.mapBody.err
	}
	if body, w.mapBody.err = urlDecode(body, w.encoding); w.mapBody.err != nil {
		return nil, w.mapBody.err
	}
	//处理body数据
	data = make(map[string]interface{})
	if len(body) != 0 {
		switch {
		case strings.Contains(ctp, "/xml"):
			data, err = types.NewXMapByXML(types.BytesToString(body))
			if err != nil {
				return nil, fmt.Errorf("xml转换为map失败:%w", err)
			}
		case strings.Contains(ctp, "/yaml") || strings.Contains(ctp, "/x-yaml"):
			w.mapBody.err = yaml.Unmarshal(body, &data)
		case strings.Contains(ctp, "/json"):
			d := json.NewDecoder(bytes.NewReader(body))
			d.UseNumber()
			w.mapBody.err = d.Decode(&data)
		case strings.Contains(ctp, "/x-www-form-urlencoded") || strings.Contains(ctp, "/form-data"):
			var values url.Values
			values, w.mapBody.err = url.ParseQuery(types.BytesToString(body))
			if w.mapBody.err != nil {
				break
			}
			for k, v := range values {
				//处理编码问题
				var buff []byte
				buff, w.mapBody.err = encoding.Decode(strings.Join(v, ","), w.encoding)
				if w.mapBody.err != nil {
					break
				}
				data[k] = types.BytesToString(buff)
			}
		}
	}
	if w.mapBody.err != nil {
		w.mapBody.err = fmt.Errorf("将%s转换为map失败:%w", body, w.mapBody.err)
		return nil, w.mapBody.err
	}

	//处理URL参数
	values := w.ctx.GetURL().Query()
	for k, v := range values {
		vs := make([]string, 0, 1)
		for _, tp := range v {
			if x, err := urlDecode([]byte(tp), w.encoding); err == nil {
				vs = append(vs, types.BytesToString(x))
			} else {
				vs = append(vs, tp)
			}
		}

		//合并body,url参数
		if x, ok := data[k]; ok {
			s := reflect.ValueOf(x)
			switch s.Kind() {
			case reflect.Map, reflect.Struct, reflect.Ptr, reflect.UnsafePointer, reflect.Func:
				return nil, fmt.Errorf("body与url中存在相同的参数，且类型为map或struct,无法进行参数合并")
			case reflect.Array, reflect.Slice: //传入数据为数组
				slice := make([]string, 0, s.Len())
				for i := 0; i < s.Len(); i++ {
					slice = append(slice, fmt.Sprint(s.Index(i).Interface()))
				}
				vs = append(vs, slice...)
			default:
				vs = append(vs, fmt.Sprint(s.Interface()))
			}
		}
		if len(vs) > 1 {
			data[k] = vs
			continue
		}
		data[k] = vs[0]
	}
	w.mapBody.value = data
	return data, nil
}

//GetFullRaw 读取所有请求Get,POST,PUT,DELETET等提交的数据
func (w *body) GetFullRaw() (body []byte, queryString string, err error) {
	//从缓存中读取
	if w.fullBody.hasRead {
		if w.fullBody.err != nil {
			return nil, "", w.fullBody.err
		}
		if p, ok := w.fullBody.value.(*params); ok {
			return p.body, p.query, nil
		}
		return nil, "", nil
	}

	//从原串中读取
	w.fullBody.hasRead = true
	p := &params{query: w.ctx.GetURL().RawQuery}
	p.body, w.fullBody.err = w.GetBody()
	if w.fullBody.err != nil {
		return nil, "", w.fullBody.err
	}
	w.fullBody.value = p
	return p.body, p.query, nil
}

//GetBody 获取POST,PUT,DELETET等提交的数据
func (w *body) GetBody() (s []byte, err error) {
	if w.rawBody.hasRead {
		if w.rawBody.err != nil {
			return nil, w.rawBody.err
		}
		return w.rawBody.value.([]byte), nil
	}
	w.rawBody.hasRead = true
	if w.ctx.ContentType() == "multipart/form-data" { //文件上专时已经进行了包体转换
		w.rawBody.value = []byte(w.ctx.GetPostForm().Encode())
		return w.rawBody.value.([]byte), nil
	}

	w.rawBody.value, w.rawBody.err = ioutil.ReadAll(w.ctx.GetBody())
	if w.rawBody.err != nil {
		return nil, fmt.Errorf("获取body发生错误:%w", w.rawBody.err)
	}
	return w.rawBody.value.([]byte), nil

}

func urlDecode(v []byte, c string) ([]byte, error) {
	if strings.EqualFold(c, encoding.UTF8) {
		return v, nil
	}
	buff, err := encoding.DecodeBytes(v, c)
	if err != nil {
		return nil, fmt.Errorf("DecodeBytes.err:%w", err)
	}
	return buff, nil
}
