package gin

import (
	"encoding/json"
	"fmt"
	"io/ioutil"

	"github.com/gin-gonic/gin"
	"github.com/micro-plat/lib4go/encoding"
	"github.com/micro-plat/lib4go/types"
)

//body 用于处理http请求的body读取
type body struct {
	*gin.Context
	body        string
	bodyReadErr error
	hasReadBody bool
}

func newBody(c *gin.Context) *body {
	return &body{Context: c}
}

//GetBodyMap 读取body并返回map
func (w *body) GetBodyMap(encoding ...string) (map[string]interface{}, error) {
	body, err := w.GetBody(encoding...)
	if err != nil || body == "" {
		return nil, err
	}
	data := make(map[string]interface{})
	err = json.Unmarshal([]byte(body), &data)
	return data, fmt.Errorf("将body转换为map失败:%w", err)
}

//GetBody 读取body返回body原字符串
func (w *body) GetBody(e ...string) (string, error) {
	encode := types.GetStringByIndex(e, 0, "utf-8")
	if w.hasReadBody {
		return w.body, w.bodyReadErr
	}
	buff, err := ioutil.ReadAll(w.Context.Request.Body)
	if err != nil {
		return "", fmt.Errorf("获取body发生错误:%w", err)
	}
	nbuff, err := encoding.DecodeBytes(buff, encode)
	if err != nil {
		return "", fmt.Errorf("获取body发生错误:%w", err)
	}
	return string(nbuff), nil
}
