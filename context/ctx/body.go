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
	"github.com/micro-plat/lib4go/types"
	"gopkg.in/yaml.v2"
)

//body 用于处理http请求的body读取
type body struct {
	ctx         context.IInnerContext
	path        *rpath
	body        []byte
	bodyReadErr error
	hasReadBody bool
}

func NewBody(c context.IInnerContext, path *rpath) *body {
	return &body{ctx: c, path: path}
}

//GetRawBodyMap 读取body原串并返回map
func (w *body) GetRawBodyMap(encoding ...string) (map[string]interface{}, error) {
	body, err := w.GetRawBody(encoding...)
	if err != nil || body == "" {
		return nil, err
	}
	data := make(map[string]interface{})
	ctp := w.ctx.ContentType()
	switch {
	case strings.Contains(ctp, "xml"):
		mxj.PrependAttrWithHyphen(false) //修改成可以转换成多层map
		data, err = mxj.NewMapXml([]byte(body))
	case strings.Contains(ctp, "yaml"):
		err = yaml.Unmarshal([]byte(body), &data)
	case strings.Contains(ctp, "json"):
		err = json.Unmarshal([]byte(body), &data)
	default:
		data["__body_"] = body
	}
	if err != nil {
		panic(fmt.Errorf("将%s转换为map失败:%w", body, err))
	}
	return data, nil
}

//GetBody 读取body @todo 【test】 特殊字符，中文,gbk
func (w *body) GetBody(e ...string) (s string, err error) {
	body, err := w.GetRawBody(e...)
	if err != nil {
		return "", err
	}
	if body != "" {
		return body, nil
	}
	return w.ctx.GetForm().Encode(), nil
}

//GetRawBody 返回body原字符串
func (w *body) GetRawBody(e ...string) (s string, err error) {
	// routerObj, err := w.path.GetRouter()
	// if err != nil {
	// 	return "", err
	// }
	encode := types.GetStringByIndex(e, 0, w.path.GetEncoding())
	if w.hasReadBody {
		if w.bodyReadErr != nil {
			return "", w.bodyReadErr
		}
		buff, err := encoding.DecodeBytes(w.body, encode)
		return string(buff), err
	}
	w.hasReadBody = true
	w.body, w.bodyReadErr = ioutil.ReadAll(w.ctx.GetBody())
	if w.bodyReadErr != nil {
		return "", fmt.Errorf("获取body发生错误:%w", w.bodyReadErr)
	}
	s, w.bodyReadErr = url.QueryUnescape(string(w.body))
	if w.bodyReadErr != nil {
		return "", fmt.Errorf("url.unescape出错:%w", w.bodyReadErr)
	}
	w.body = []byte(s)
	var buff []byte
	buff, w.bodyReadErr = encoding.DecodeBytes(w.body, encode)
	return string(buff), w.bodyReadErr
}
