package static

//Static 设置静态文件配置
type Static struct {
	Dir       string   `json:"dir,omitempty" valid:"ascii"`
	Archive   string   `json:"archive,omitempty" valid:"ascii"`
	Prefix    string   `json:"prefix,omitempty" valid:"ascii"`
	Exts      []string `json:"exts,omitempty" valid:"ascii"`
	Exclude   []string `json:"exclude,omitempty" valid:"ascii"`
	FirstPage string   `json:"first-page,omitempty" valid:"ascii"`
	Rewriters []string `json:"rewriters,omitempty" valid:"ascii"`
	Disable   bool     `json:"disable,omitempty"`
}
