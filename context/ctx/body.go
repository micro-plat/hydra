package ctx

import (
	"encoding/json"
	"fmt"
	"io/ioutil"

	"github.com/micro-plat/hydra/context"
	"github.com/micro-plat/lib4go/encoding"
	"github.com/micro-plat/lib4go/types"
)

//body 用于处理http请求的body读取
type body struct {
	ctx         context.IInnerContext
	body        []byte
	bodyReadErr error
	hasReadBody bool
}

func newBody(c context.IInnerContext) *body {
	return &body{ctx: c}
}

//GetBodyMap 读取body并返回map
func (w *body) GetBodyMap(encoding ...string) (map[string]interface{}, error) {
	body, err := w.GetBody(encoding...)
	if err != nil || body == "" {
		return nil, err
	}
	data := make(map[string]interface{})
	err = json.Unmarshal([]byte(body), &data)
	if err != nil {
		return nil, fmt.Errorf("将body转换为map失败:%w", err)
	}
	return data, nil
}

//GetBody 读取body返回body原字符串
func (w *body) GetBody(e ...string) (s string, err error) {
	encode := types.GetStringByIndex(e, 0, "utf-8")
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
	var buff []byte
	buff, w.bodyReadErr = encoding.DecodeBytes(w.body, encode)
	if w.bodyReadErr != nil {
		return "", fmt.Errorf("获取body发生错误:%w", w.bodyReadErr)
	}
	return string(buff), nil
}
