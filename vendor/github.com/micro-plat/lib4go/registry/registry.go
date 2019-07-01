package registry

type ValueWatcher interface {
	GetValue() ([]byte, int32)
	GetError() error
	GetPath() string
}
type ChildrenWatcher interface {
	GetValue() ([]string, int32)
	GetPath() string
	GetError() error
}
