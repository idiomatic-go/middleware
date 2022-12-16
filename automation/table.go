package automation

import (
	"net/http"
	"sync"
)

var this *table

type table struct {
	mu         sync.RWMutex
	pmu        sync.RWMutex
	defaultAct *actuator
	match      Matcher
	actuators  map[string]*actuator
	pings      map[string]*pingAction
}

func NewTable() Automation {
	this = newTable()
	return this
}

func newTable() *table {
	t := new(table)
	t.pings = make(map[string]*pingAction, 100)
	t.actuators = make(map[string]*actuator, 100)
	t.defaultAct = &actuator{name: DefaultName, timeout: newTimeout(NewTimeoutConfig(NilValue, NilValue)), limit: nil}
	t.match = func(req *http.Request) (name string) {
		return ""
	}
	return t
}

func (t *table) SetDefault(name string, tc *TimeoutConfig, rc *RateLimitConfig) {
	t.mu.Lock()
	defer t.mu.Unlock()
	if name == "" {
		name = DefaultName
	}
	t.defaultAct = &actuator{name: name, timeout: newTimeout(tc), limit: nil}
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

func (t *table) IsPingEnabled(name string) bool {
	if name == "" {
		return false
	}
	t.pmu.RLock()
	if p, ok := t.pings[name]; ok {
		return p.IsEnabled()
	}
	t.pmu.Unlock()
	return false
}

func (t *table) Lookup(req *http.Request) Actuator {
	name := t.match(req)
	if name != "" {
		if r := t.LookupByName(name); r != nil {
			return r
		}
	}
	t.mu.RLock()
	defer t.mu.RUnlock()
	return t.defaultAct
}

func (t *table) LookupByName(name string) Actuator {
	if name == "" {
		return nil
	}
	t.mu.RLock()
	defer t.mu.RUnlock()
	if r, ok := t.actuators[name]; ok {
		return r
	}
	return nil
}

func (t *table) Add(name string, pc *PingConfig, tc *TimeoutConfig, rc *RateLimitConfig) bool {
	if IsEmpty(name) {
		return false
	}
	if pc != nil {
		t.pmu.Lock()
		t.pings[name] = newPingAction(pc.enabled)
		t.pmu.Unlock()
	}
	t.mu.Lock()
	defer t.mu.Unlock()
	t.actuators[name] = &actuator{name: name, timeout: newTimeout(tc), limit: nil}
	return false
}

func configureTimeout(timeout int, statusCode int) {
	/*t.mu.Lock()
	defer t.mu.Unlock()
	if len(v) == 0 {
		a.current = a.Default
		return
	}
	if timeout, ok := v[0].(int); ok {
		if timeout <= 0 {
			timeout = NilValue
		}
		a.current = timeout
	}

	*/
}
