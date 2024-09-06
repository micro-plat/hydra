package global

//db 数据库处理逻辑
type local struct {
	handlers []func() error
}

//AddHandler 添加处理函数
func (d *local) AddHandler(fs ...interface{}) {
	for _, fn := range fs {
		var nfunc func() error
		hasMatch := false
		if fx, ok := fn.(func()); ok {
			hasMatch = true
			nfunc = func() error {
				fx()
				return nil
			}
		}
		if fx, ok := fn.(func() error); ok {
			hasMatch = true
			nfunc = fx
		}
		if !hasMatch {
			panic("函数签名格式不正确，支持的格式有func(){} 或 func()error{}")
		}
		d.handlers = append(d.handlers, nfunc)
	}
}

//GetHandlers 获取所有处理函数
func (d *local) GetHandlers() []func() error {
	return d.handlers
}
