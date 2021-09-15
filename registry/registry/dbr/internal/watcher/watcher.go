package watcher

import (
	log "github.com/lib4dev/cli/logger"
	"github.com/micro-plat/hydra/components/dbs"
	"github.com/micro-plat/hydra/registry/registry/mysql/internal/river"
)

type Watcher struct {
	closeCh chan struct{}
	db      dbs.IDB
}

func NewWatcher() (*Watcher, error) {
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
