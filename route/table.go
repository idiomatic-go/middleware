package route

import (
	"golang.org/x/time/rate"
	"net/http"
	"sync"
)

type Routes interface {
	SetDefault(r *Route)
	SetMatchFn(fn MatchFn)
	Lookup(req *http.Request) Route
	LookupByName(name string) (Route, bool)
	Add(r *Route) bool
	AddWithLimiter(r *Route, max rate.Limit, b int) bool
	UpdateTimeout(name string, timeout int) bool
	UpdateLimiter(name string, max rate.Limit, b int) bool
	Remove(name string) bool
	RemoveLimiter(name string) bool
}

type table struct {
	mu           sync.RWMutex
	routes       map[string]*Route
	defaultRoute *Route
	match        MatchFn
}

func NewTable(r *Route) Routes {
	t := new(table)
	t.defaultRoute = r
	t.match = func(req *http.Request) (name string) {
		return ""
	}
	return t
}

func (t *table) SetDefault(r *Route) {
	if t != nil && r != nil {
		t.defaultRoute = r
	}
}
func (t *table) SetMatchFn(fn MatchFn) {
	if fn != nil {
		t.match = fn
	}
}

func (t *table) Lookup(req *http.Request) Route {
	name := t.match(req)
	if name == "" {
		return *t.defaultRoute
	}
	if r, ok := t.LookupByName(name); ok {
		return r
	}
	return *t.defaultRoute
}

func (t *table) LookupByName(name string) (Route, bool) {
	if t == nil || name == "" {
		return Route{}, false
	}
	t.mu.RLock()
	defer t.mu.RUnlock()
	if r, ok := t.routes[name]; ok {
		return *r, true
	}
	return Route{}, false
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

func (t *table) AddWithLimiter(r *Route, max rate.Limit, b int) bool {
	if t == nil || r == nil || r.Name == "" {
		return false
	}
	t.mu.Lock()
	defer t.mu.Unlock()
	if _, ok := t.routes[r.Name]; ok {
		return false
	}
	r.rateLimiter = rate.NewLimiter(max, b)
	t.routes[r.Name] = r
	return true
}

func (t *table) UpdateTimeout(name string, timeout int) bool {
	if t == nil || name == "" || timeout <= 0 {
		return false
	}
	t.mu.Lock()
	if r, ok := t.routes[name]; ok {
		r.Timeout = timeout
	}
	t.mu.Unlock()
	return true
}

func (t *table) UpdateLimiter(name string, max rate.Limit, b int) bool {
	if t == nil || name == "" || b < 0 {
		return false
	}
	t.mu.Lock()
	if r, ok := t.routes[name]; ok {
		r.rateLimiter = rate.NewLimiter(max, b)
	}
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

func (t *table) RemoveLimiter(name string) bool {
	if t == nil || name == "" {
		return false
	}
	t.mu.Lock()
	if r, ok := t.routes[name]; ok {
		r.rateLimiter = nil
		return true
	}
	t.mu.Unlock()
	return false
}
