package wss

import (
	"fmt"
	"sync"
)

type tunnelPool struct {
	mu      sync.RWMutex
	groups  map[string]map[string]*session
	pending map[string]chan *Frame
	next    int
}

func newTunnelPool() *tunnelPool {
	return &tunnelPool{
		groups:  make(map[string]map[string]*session),
		pending: make(map[string]chan *Frame),
	}
}

func (p *tunnelPool) add(s *session) {
	p.mu.Lock()
	defer p.mu.Unlock()
	if _, ok := p.groups[s.group]; !ok {
		p.groups[s.group] = make(map[string]*session)
	}
	if old := p.groups[s.group][s.id]; old != nil {
		old.close()
	}
	p.groups[s.group][s.id] = s
}

func (p *tunnelPool) remove(s *session) {
	p.mu.Lock()
	defer p.mu.Unlock()
	if g, ok := p.groups[s.group]; ok {
		if g[s.id] == s {
			delete(g, s.id)
		}
		if len(g) == 0 {
			delete(p.groups, s.group)
		}
	}
}

func (p *tunnelPool) pick(group string) (*session, error) {
	p.mu.Lock()
	defer p.mu.Unlock()
	g := p.groups[group]
	if len(g) == 0 {
		return nil, fmt.Errorf("group %s has no available client", group)
	}
	i := 0
	target := p.next % len(g)
	p.next++
	for _, s := range g {
		if i == target {
			return s, nil
		}
		i++
	}
	return nil, fmt.Errorf("group %s has no available client", group)
}

func (p *tunnelPool) wait(id string) chan *Frame {
	p.mu.Lock()
	defer p.mu.Unlock()
	ch := make(chan *Frame, 32)
	p.pending[id] = ch
	return ch
}

func (p *tunnelPool) send(id string, frame *Frame) (ok bool) {
	defer func() {
		if recover() != nil {
			ok = false
		}
	}()
	p.mu.RLock()
	ch := p.pending[id]
	p.mu.RUnlock()
	if ch == nil {
		return false
	}
	ch <- frame
	return true
}

func (p *tunnelPool) done(id string, frame *Frame) {
	p.mu.Lock()
	ch := p.pending[id]
	delete(p.pending, id)
	p.mu.Unlock()
	if ch != nil {
		ch <- frame
		close(ch)
	}
}

func (p *tunnelPool) cancel(id string) {
	p.mu.Lock()
	ch := p.pending[id]
	delete(p.pending, id)
	p.mu.Unlock()
	if ch != nil {
		close(ch)
	}
}

func (p *tunnelPool) closeAll() {
	p.mu.Lock()
	sessions := make([]*session, 0)
	for _, g := range p.groups {
		for _, s := range g {
			sessions = append(sessions, s)
		}
	}
	p.groups = make(map[string]map[string]*session)
	pending := p.pending
	p.pending = make(map[string]chan *Frame)
	p.mu.Unlock()

	for _, s := range sessions {
		s.close()
	}
	for _, ch := range pending {
		close(ch)
	}
}
