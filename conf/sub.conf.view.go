package conf

//View web服务器view解析配置
type View struct {
	Path    string   `json:"path" valid:"ascii,required"`
	Left    string   `json:"left,omitempty" valid:"ascii"`
	Right   string   `json:"right,omitempty" valid:"ascii"`
	Files   []string `json:"-"`
	Disable bool     `json:"disable,omitempty"`
}

//NewView 构建web服务view配置
func NewView(path string) *View {
	return &View{
		Path: path,
	}
}

//WithExpression 设置view绑定变量的左，右解析符
func (a *View) WithExpression(left string, right string) *View {
	a.Left = left
	a.Right = right
	return a
}

//WithEnable 启用配置
func (a *View) WithEnable() *View {
	a.Disable = false
	return a
}

//WithDisable 禁用配置
func (a *View) WithDisable() *View {
	a.Disable = false
	return a
}
