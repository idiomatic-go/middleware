package route

import (
	"golang.org/x/time/rate"
	"net/http"
	"reflect"
	"sync"
)

type Routes interface {
	SetDefault(r *Route)
	SetMatchFn(fn MatchFn)
	Lookup(req *http.Request) Route
	LookupByName(name string) (Route, bool)
	Add(r Route) bool
	//AddWithLimiter(r *Route, max rate.Limit, b int) bool
	UpdateTimeout(name string, timeout int) bool
	UpdateLimit(name string, max rate.Limit) bool
	UpdateBurst(name string, b int) bool
	Remove(name string) bool
	RemoveLimiter(name string) bool
}

type table struct {
	mu           sync.RWMutex
	routes       map[string]route
	defaultRoute Route
	match        MatchFn
}

func NewTable() Routes {
	t := new(table)
	t.defaultRoute = NewRoute(DefaultName)
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
	//if rt,ok :=
	t.mu.Unlock()
}

func (t *table) SetMatchFn(fn MatchFn) {
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
		if r, ok := t.LookupByName(name); ok {
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
		return *r, true
	}
	return Route{}, false
}

func (t *table) Add(r Route) bool {
	if t == nil || r == nil || IsEmpty(r.Name()) {
		return false
	}
	t.mu.Lock()
	defer t.mu.Unlock()
	a := reflect.TypeOf(r)
	if rt, ok := a.(route); ok {

	}
	if _, ok := t.routes[r.Name()]; ok {
		return false
	}
	t.routes[r.Name()] = r
	return true
}

func (t *table) AddWithLimiter(r Route, max rate.Limit, b int) bool {
	if t == nil || r == nil || IsEmpty(r.Name()) {
		return false
	}
	t.mu.Lock()
	defer t.mu.Unlock()
	if _, ok := t.routes[r.Name]; ok {
		return false
	}
	r.Current.Limit = max
	r.Current.Burst = b
	r.Original.Limit = max
	r.Original.Burst = b
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
		r.Current.Timeout = timeout
	}
	t.mu.Unlock()
	return true
}

func (t *table) UpdateLimit(name string, max rate.Limit) bool {
	if t == nil || name == "" {
		return false
	}
	t.mu.Lock()
	if r, ok := t.routes[name]; ok {
		if r.IsRateLimiter() {
			r.Limit = max
			r.rateLimiter.SetLimit(max)
		}
	}
	t.mu.Unlock()
	return true
}

func (t *table) UpdateBurst(name string, b int) bool {
	if t == nil || name == "" || b < 0 {
		return false
	}
	t.mu.Lock()
	if r, ok := t.routes[name]; ok {
		if r.IsRateLimiter() {
			r.Burst = b
			r.rateLimiter.SetBurst(b)
		}
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
