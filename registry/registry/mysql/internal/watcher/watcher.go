package watcher

import (
	"runtime"

	log "github.com/lib4dev/cli/logger"
	"github.com/micro-plat/hydra/registry/registry/mysql/internal/river"
)

type Watcher struct {
	closeCh chan struct{}
	*river.River
}

func NewWatcher(c *river.Config) (*Watcher, error) {
	r, err := river.NewRiver(c)
	if err != nil {
		return nil, err
	}
	return &Watcher{
		River:   r,
		closeCh: make(chan struct{}),
	}, nil
}

func (w *Watcher) Watch() {
	runtime.GOMAXPROCS(runtime.NumCPU())

	done := make(chan struct{}, 1)
	go func() {
		w.Run()
		done <- struct{}{}
	}()

	select {
	case n := <-w.closeCh:
		log.Log.Infof("receive signal %v, closing", n)
	}

	w.River.Close()
	<-done
}

func (w *Watcher) Close() {
	close(w.closeCh)
}
