package route

import "sync"

type Routes interface {
	Lookup(name string) *Route
	Add(r *Route) bool
	Update(r *Route) bool
	Remove(name string) bool
}

type table struct {
	mu     sync.RWMutex
	routes map[string]*Route
}

func NewTable() Routes {
	return new(table)
}

func (t *table) Lookup(name string) *Route {
	if t == nil || name == "" {
		return nil
	}
	t.mu.RLock()
	defer t.mu.RUnlock()
	return t.routes[name]
}

func (t *table) Add(r *Route) bool {
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

func (t *table) Update(r *Route) bool {
	if t == nil || r == nil || r.Name == "" {
		return false
	}
	t.mu.Lock()
	t.routes[r.Name] = r
	t.mu.Unlock()
	return true
}

func (t *table) Remove(name string) bool {
	if t == nil || name == "" {
		return false
	}
	t.mu.Lock()
	delete(t.routes, name)
	t.mu.Unlock()
	return true
}
