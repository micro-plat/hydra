package render

import (
	"errors"
	"fmt"

	"github.com/micro-plat/hydra/conf"
	"github.com/micro-plat/hydra/global"
	"github.com/micro-plat/lib4go/tgo"
	"github.com/micro-plat/lib4go/types"
)

const (
	//TypeNodeName render配置节点名
	TypeNodeName = "render"

	//scriptName 脚本中render结果值
	scriptName = "render"
)

//Render 响应模板信息
type Render struct {
	//Disable 禁用
	Disable bool `json:"-"`
	tengo   *tgo.VM
}

//GetConf 设置GetRender配置
func GetConf(cnf conf.IServerConf) (rsp *Render, err error) {
	script, err := cnf.GetSubConf(TypeNodeName)
	if errors.Is(err, conf.ErrNoSetting) {
		return &Render{Disable: true}, nil
	}
	if err != nil {
		return nil, fmt.Errorf("%s配置有误:%v", TypeNodeName, err)
	}
	render := &Render{}
	render.tengo, err = tgo.New(string(script.GetRaw()), tgo.WithModule(global.GetTGOModules()...))
	if err != nil {
		return nil, fmt.Errorf("%s脚本错误:%v", TypeNodeName, err)
	}
	return render, nil
}

//Get 获取转换结果
func (r *Render) Get() (*Result, bool, error) {

	//执行脚本，获取render结果
	result, err := r.tengo.Run()
	if err != nil {
		return nil, false, err
	}

	//获取脚本执行结果
	sresult := result.GetArray(scriptName)
	if len(sresult) >= 2 {
		ct := ""
		if len(sresult) > 2 {
			ct = types.GetString(sresult[2])
		}
		return &Result{
			Status:      types.GetInt(sresult[0]),
			Content:     types.GetString(sresult[1]),
			ContentType: ct,
		}, true, nil
	}
	return nil, false, nil

}
