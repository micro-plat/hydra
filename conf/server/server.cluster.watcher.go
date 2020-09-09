package server

import (
	"github.com/micro-plat/hydra/conf"
	"github.com/micro-plat/lib4go/utility"
)

type cwatcher struct {
	c          *Cluster
	notifyChan chan conf.ICNode
	id         string
}

func newCWatcher(c *Cluster) *cwatcher {
	return &cwatcher{
		c:          c,
		id:         utility.GetGUID(),
		notifyChan: make(chan conf.ICNode, 1),
	}
}

func (w *cwatcher) Notify() chan conf.ICNode {
	return w.notifyChan
}

func (w *cwatcher) notify(n conf.ICNode) {
	select {
	case w.notifyChan <- n:
	default:
	}
}
func (w *cwatcher) Close() error {
	w.c.removeWatcher(w.id)
	return nil
}
