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
	Remove(name string)

	SetTimeout(name string, timeout int)
	ResetTimeout(name string)
	DisableTimeout(name string)

	//SetLimiter(name string, max rate.Limit, burst int)
	SetLimit(name string, max rate.Limit)
	SetBurst(name string, burst int)
	ResetLimiter(name string)
	DisableLimiter(name string)
}

type table struct {
	mu           sync.RWMutex
	routes       map[string]*route
	defaultRoute Route
	match        Matcher
}

func NewTable() Routes {
	return newTable()
}

func newTable() *table {
	t := new(table)
	t.routes = make(map[string]*route, 100)
	t.defaultRoute = newRoute(DefaultName)
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
	if fn == nil {
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
	if name == "" {
		return nil
	}
	t.mu.RLock()
	defer t.mu.RUnlock()
	if r, ok := t.routes[name]; ok {
		return r
	}
	return nil
}

func (t *table) Exists(name string) bool {
	if name == "" {
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
	if r == nil || IsEmpty(r.Name()) {
		return false
	}
	t.mu.Lock()
	defer t.mu.Unlock()
	if _, ok := t.routes[r.Name()]; ok {
		return false
	}
	route := r.(*route)
	t.routes[r.Name()] = route
	return true
}

func (t *table) Remove(name string) {
	if name == "" {
		return
	}
	t.mu.Lock()
	delete(t.routes, name)
	t.mu.Unlock()
}

func (t *table) count() int {
	t.mu.RLock()
	defer t.mu.RUnlock()
	return len(t.routes)
}

func (t *table) isEmpty() bool {
	return t.count() == 0
}
