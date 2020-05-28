package render

//Render 响应模板信息
type Render struct {

	//模板项
	Content string `json:"template,omitempty" valid:"required" toml:"content,omitempty"`

	//Disable 禁用
	Disable bool `json:"disable,omitemptye" toml:"disable,omitempty"`
}

//Tmplt 输出控件模板
type Tmplt map[string]string

//NewTmplt 构建模板
func NewTmplt(opts ...Option) Tmplt {
	r := make(Tmplt)
	for _, opt := range opts {
		opt(r)
	}
	return r
}

//Get 获取转换结果
func (r *Render) Get(funcs map[string]interface{}, i interface{}) (string, error) {
	return translate(r.Content, funcs, i)
}
