package ctx

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/url"
	"strings"

	"github.com/clbanning/mxj"
	"github.com/micro-plat/hydra/context"
	"github.com/micro-plat/lib4go/encoding"
	"gopkg.in/yaml.v2"
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

//GetBodyMap 读取body原串并返回map
func (w *body) GetBodyMap() (map[string]interface{}, error) {

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

	if body, _, w.mapBody.err = w.GetRequestParams(); w.mapBody.err != nil {
		return nil, w.mapBody.err
	}
	if body, w.mapBody.err = urlDecode(body, w.encoding); w.mapBody.err != nil {
		return nil, w.mapBody.err
	}
	//处理body数据
	data := make(map[string]interface{})
	if len(body) != 0 {
		ctp := strings.ToLower(w.ctx.ContentType())
		switch {
		case strings.Contains(ctp, "/xml"):
			mxj.PrependAttrWithHyphen(false) //修改成可以转换成多层map
			data, w.mapBody.err = mxj.NewMapXml([]byte(body))
		case strings.Contains(ctp, "/yaml") || strings.Contains(ctp, "/x-yaml"):
			w.mapBody.err = yaml.Unmarshal([]byte(body), &data)
		case strings.Contains(ctp, "/json"):
			w.mapBody.err = json.Unmarshal([]byte(body), &data)
		case strings.Contains(ctp, "/x-www-form-urlencoded") || strings.Contains(ctp, "/form-data"):
			var values url.Values
			values, w.mapBody.err = url.ParseQuery(string(body))
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
				data[k] = string(buff)
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
		vs := make([]byte, 0, 10)
		for _, tp := range v {
			if x, err := urlDecode([]byte(tp), w.encoding); err == nil {
				vs = append(vs, x...)
			} else {
				vs = append(vs, tp...)
			}
		}

		if x, ok := data[k]; ok {
			vs = append(vs, []byte(",")...)
			vs = append(vs, []byte(x.(string))...)
		}
		data[k] = string(vs)
	}
	w.mapBody.value = data
	return data, nil
}

//GetBody 读取所有请求Get,POST,PUT,DELETET等提交的数据
func (w *body) GetRequestParams() (body []byte, queryString string, err error) {
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
