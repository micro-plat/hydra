package dbr

type valueEntity struct {
	Value   []byte
	version int32
	path    string
	Err     error
}
type childrenEntity struct {
	children []string
	version  int32
	path     string
	Err      error
}

func (v *valueEntity) GetPath() string {
	return v.path
}
func (v *valueEntity) GetValue() ([]byte, int32) {
	return v.Value, v.version
}
func (v *valueEntity) GetError() error {
	return v.Err
}

func (v *childrenEntity) GetValue() ([]string, int32) {
	return v.children, v.version
}
func (v *childrenEntity) GetError() error {
	return v.Err
}
func (v *childrenEntity) GetPath() string {
	return v.path
}
