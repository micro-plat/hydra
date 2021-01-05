package global

var funcs = make([]func() error, 0, 1)

//OnReady 系统准备好后稍后执行
func OnReady(fs ...interface{}) {
	for _, fn := range fs {
		nfunc := getFuncs(fn)
		if !isReady {
			funcs = append(funcs, nfunc)
			continue
		}
		if err := nfunc(); err != nil {
			panic(err)
		}
	}
}

//OnReadying 在系统准备好后首先执行
func OnReadying(fs ...interface{}) {
	for _, fn := range fs {
		nfunc := getFuncs(fn)
		if !isReady {
			nlist := make([]func() error, 0, len(funcs)+1)
			nlist = append(nlist, nfunc)
			nlist = append(nlist, funcs...)
			funcs = nlist
			continue
		}
		if err := nfunc(); err != nil {
			panic(err)
		}
	}
}

func getFuncs(fn interface{}) func() error {
	var nfunc func() error
	hasMatchType := false
	if fx, ok := fn.(func()); ok {
		hasMatchType = true
		nfunc = func() error {
			fx()
			return nil
		}
	}
	if fx, ok := fn.(func() error); ok {
		hasMatchType = true
		nfunc = fx
	}
	if !hasMatchType {
		panic("函数签名格式不正确，支持的格式有func(){} 或 func()error{}")
	}
	return nfunc
}

func doReadyFuncs() error {
	for _, f := range funcs {
		if err := f(); err != nil {
			return err
		}
	}
	return nil
}
