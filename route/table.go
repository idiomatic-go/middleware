package route

import "sync"

type routes interface {
	lookup(name string) *Route
	add(r *Route) bool
	update(r *Route) bool
	remove(name string) bool
}

type table struct {
	mu     sync.RWMutex
	routes map[string]*Route
}

func NewTable() routes {
	return new(table)
}

func (t *table) lookup(name string) *Route {
	if t == nil || name == "" {
		return nil
	}
	t.mu.RLock()
	defer t.mu.RUnlock()
	return t.routes[name]
}

func (t *table) add(r *Route) bool {
	if t == nil || r == nil || r.Name == "" {
		return false
	}
	t.mu.Lock()
	defer t.mu.Unlock()
	if _, ok := t.routes[r.Name]; ok {
		return false
	}
	t.routes[r.Name] = r
	return true
}

func (t *table) update(r *Route) bool {
	if t == nil || r == nil || r.Name == "" {
		return false
	}
	t.mu.Lock()
	t.routes[r.Name] = r
	t.mu.Unlock()
	return true
}

func (t *table) remove(name string) bool {
	if t == nil || name == "" {
		return false
	}
	t.mu.Lock()
	delete(t.routes, name)
	t.mu.Unlock()
	return true
}
