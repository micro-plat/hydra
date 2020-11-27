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

type bodyValue struct {
	hasRead bool
	value   interface{}
	err     error
}

//body 用于处理http请求的body读取
type body struct {
	ctx          context.IInnerContext
	encoding     string
	rawBody      bodyValue
	fullBody     bodyValue
	fullFormBody string
	mapBody      bodyValue
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
		return w.mapBody.value.(map[string]interface{}), w.mapBody.err
	}

	//从body中读取原处理流
	w.mapBody.hasRead = true
	var body string
	var noNeedReadURLQuery bool
	body, w.mapBody.err = w.GetBody()
	data := make(map[string]interface{})
	if w.mapBody.err != nil {
		return data, w.mapBody.err
	}

	if body != "" {
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
			noNeedReadURLQuery = true
			var values url.Values
			values, w.mapBody.err = url.ParseQuery(w.fullFormBody)
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
	if !noNeedReadURLQuery {
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
				vs = append(vs, []byte(x.(string))...)
			}
			data[k] = string(vs)
		}
	}
	//处理返回结果
	w.mapBody.value = data
	return data, nil
}

//GetBody 读取所有请求Get,POST,PUT,DELETET等提交的数据
func (w *body) GetBody() (s string, err error) {
	//从缓存中读取
	if w.fullBody.hasRead {
		if w.fullBody.err != nil {
			return "", w.fullBody.err
		}
		return w.fullBody.value.(string), w.fullBody.err
	}

	//从原串中读取
	w.fullBody.hasRead = true
	var buff []byte
	buff, w.fullBody.err = w.GetRawBody()
	if w.fullBody.err != nil {
		return "", w.fullBody.err
	}

	ctp := strings.ToLower(w.ctx.ContentType())
	if strings.Contains(ctp, "/x-www-form-urlencoded") || strings.Contains(ctp, "/form-data") {

		formRaw := w.ctx.GetForm().Encode()
		replaceMent := "%s%s"
		if len(formRaw) > 0 && len(buff) > 0 {
			replaceMent = "%s&%s"
		}
		w.fullFormBody = fmt.Sprintf(replaceMent, string(buff), formRaw)
		s, err := url.QueryUnescape(w.fullFormBody)
		if err != nil {
			w.fullBody.err = fmt.Errorf("GetBody.QueryUnescape.err:%w", err)
			return "", w.fullBody.err
		}
		buff = []byte(s)
	}

	//处理编码问题
	buff, err = encoding.DecodeBytes(buff, w.encoding)
	if err != nil {
		w.fullBody.err = fmt.Errorf("GetBody.DecodeBytes.err:%w", err)
		return "", w.fullBody.err
	}
	w.fullBody.value = string(buff)
	return w.fullBody.value.(string), nil

}

//GetRawBody 获取POST,PUT,DELETET等提交的数据
func (w *body) GetRawBody() (s []byte, err error) {
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
	// s, err := url.QueryUnescape(string(v))
	// if err != nil {
	// 	return nil, fmt.Errorf("QueryUnescape.err:%w", err)
	// }
	if strings.EqualFold(c, encoding.UTF8) {
		return v, nil
	}
	buff, err := encoding.DecodeBytes(v, c)
	if err != nil {
		return nil, fmt.Errorf("DecodeBytes.err:%w", err)
	}
	return buff, nil
}
