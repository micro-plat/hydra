package balancer

import (
	"errors"
	"testing"
	"time"

	"github.com/micro-plat/hydra/registry"
	"github.com/micro-plat/lib4go/ut"
)

type testfileChecker struct {
	modTime time.Time
	apis    map[string]bool
	files   map[string]string
}

func (f testfileChecker) Exists(filename string) bool {
	if v, ok := f.apis[filename]; ok {
		return v
	}
	return false
}
func (f testfileChecker) LastModeTime(path string) (t time.Time, err error) {
	return f.modTime, nil
}
func (f testfileChecker) ReadDir(path string) (r []string, err error) {
	return []string{"merchant"}, nil
}
func (f testfileChecker) ReadAll(path string) (buf []byte, err error) {
	if _, ok := f.files[path]; ok {
		return []byte(""), nil
	}
	return nil, errors.New("file not exists")
}
func (f testfileChecker) CreateFile(fileName string, data string) error {
	return nil
}
func (f testfileChecker) Delete(fileName string) error {
	return nil
}
func (f testfileChecker) WriteFile(fileName string, data string) error {
	return nil
}
func TestWatcher1(t *testing.T) {
	client, err := registry.NewLocalRegistryWithChcker(testfileChecker{})
	ut.ExpectSkip(t, err, nil)
	w := &Watcher{client: client, service: "", closeCh: make(chan struct{})}
	updaters, err := w.Next()
	ut.ExpectSkip(t, err, nil)
	ut.ExpectSkip(t, len(updaters), 1)
	ut.Expect(t, updaters[0].Addr, "merchant")
}
