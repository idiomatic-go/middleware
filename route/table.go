package route

import (
	"golang.org/x/time/rate"
	"net/http"
	"sync"
)

type Routes interface {
	SetDefault(r Route)
	SetMatcher(fn Matcher)

	Lookup(req *http.Request) Route
	LookupByName(name string) Route

	Exists(name string) bool
	Add(r Route) bool
	Remove(name string) bool

	SetTimeout(name string, timeout int) bool
	ResetTimeout(name string) bool
	DisableTimeout(name string) bool

	SetLimiter(name string, max rate.Limit, burst int) bool
	ResetLimiter(name string) bool
	DisableLimiter(name string) bool

	t() *table
}

type table struct {
	mu           sync.RWMutex
	routes       map[string]route
	defaultRoute Route
	match        Matcher
}

func NewTable() Routes {
	t := new(table)
	// TODO : rework
	t.defaultRoute, _ = NewRoute(DefaultName)
	t.match = func(req *http.Request) (name string) {
		return ""
	}
	return t
}

func (t *table) SetDefault(r Route) {
	if t == nil || r == nil {
		return
	}
	t.mu.Lock()
	t.defaultRoute = r
	t.mu.Unlock()
}

func (t *table) SetMatcher(fn Matcher) {
	if t == nil || fn == nil {
		return
	}
	t.mu.Lock()
	t.match = fn
	t.mu.Unlock()
}

func (t *table) Lookup(req *http.Request) Route {
	name := t.match(req)
	if name != "" {
		if r := t.LookupByName(name); r != nil {
			return r
		}
	}
	t.mu.RLock()
	defer t.mu.RUnlock()
	return t.defaultRoute
}

func (t *table) LookupByName(name string) Route {
	if t == nil || name == "" {
		return nil
	}
	t.mu.RLock()
	defer t.mu.RUnlock()
	if r, ok := t.routes[name]; ok {
		return &r
	}
	return nil
}

func (t *table) Exists(name string) bool {
	if t == nil || IsEmpty(name) {
		return false
	}
	t.mu.RLock()
	defer t.mu.RUnlock()
	if _, ok := t.routes[name]; ok {
		return true
	}
	return false
}

func (t *table) Add(r Route) bool {
	if t == nil || r == nil || IsEmpty(r.Name()) {
		return false
	}
	t.mu.Lock()
	defer t.mu.Unlock()
	if _, ok := t.routes[r.Name()]; ok {
		return false
	}
	route := r.t()
	//if route.original.limit > 0 {
	//	route.rateLimiter = rate.NewLimiter(route.original.limit, route.original.burst)
	//}
	t.routes[r.Name()] = route
	return true
}

func (t *table) Remove(name string) bool {
	if t == nil || IsEmpty(name) {
		return false
	}
	t.mu.Lock()
	delete(t.routes, name)
	t.mu.Unlock()
	return true
}

func (t *table) count() int {
	t.mu.RLock()
	defer t.mu.RUnlock()
	return len(t.routes)
}

func (t *table) isEmpty() bool {
	return t.count() == 0
}

func (t *table) t() *table {
	return t
}
