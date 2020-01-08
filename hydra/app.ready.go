package hydra

type funs struct {
	list []func() error
}

func newFuns() *funs {
	return &funs{list: make([]func() error, 0, 1)}
}

//Call 所有准备函数
func (f *funs) Call() error {
	for _, f := range f.list {
		if err := f(); err != nil {
			return err
		}
	}
	return nil
}

//Ready 等待app准备好，正式启动时退出
func (m *MicroApp) Ready(f func() error) {
	m.funs.list = append(m.funs.list, f)
}
