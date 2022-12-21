package automation

import (
	"errors"
	"net/http"
	"sync"
)

type table struct {
	mu         sync.RWMutex
	match      Matcher
	defaultAct *actuator
	actuators  map[string]*actuator
}

func NewTable() Automation {
	this := newTable()
	return this
}

func newTable() *table {
	t := new(table)
	t.actuators = make(map[string]*actuator, 100)
	t.defaultAct = &actuator{name: DefaultName,
		logger:      defaultLogger,
		timeout:     newTimeout(DefaultName, nil, t),
		rateLimiter: newRateLimiter(DefaultName, nil, t),
	}
	t.match = func(req *http.Request) (name string) {
		return ""
	}
	return t
}

func (t *table) SetDefault(name string, tc *TimeoutConfig, rlc *RateLimiterConfig, cbc *CircuitBreakerConfig, fc *FailoverConfig) {
	t.mu.Lock()
	defer t.mu.Unlock()
	if name == "" {
		name = DefaultName
	}
	t.defaultAct = &actuator{name: name,
		logger:         defaultLogger,
		timeout:        newTimeout(name, tc, t),
		rateLimiter:    newRateLimiter(name, rlc, t),
		circuitBreaker: newCircuitBreaker(name, cbc, t),
		failover:       newFailover(name, fc, t)}
}

func (t *table) SetMatcher(fn Matcher) {
	if fn == nil {
		return
	}
	t.mu.Lock()
	t.match = fn
	t.mu.Unlock()
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

func (t *table) Add(name string, tc *TimeoutConfig, rlc *RateLimiterConfig, cbc *CircuitBreakerConfig, fc *FailoverConfig) bool {
	if IsEmpty(name) {
		return false
	}
	t.mu.Lock()
	defer t.mu.Unlock()
	t.actuators[name] = &actuator{name: name,
		logger:         defaultLogger,
		timeout:        newTimeout(name, tc, t),
		rateLimiter:    newRateLimiter(name, rlc, t),
		circuitBreaker: newCircuitBreaker(name, cbc, t),
		failover:       newFailover(name, fc, t)}
	return true
}

func (t *table) exists(name string) bool {
	if name == "" {
		return false
	}
	t.mu.RLock()
	defer t.mu.RUnlock()
	if _, ok := t.actuators[name]; ok {
		return true
	}
	return false
}

func (t *table) update(name string, act *actuator) error {
	if name == "" || act == nil {
		return errors.New("invalid argument : name or actuator is nil or empty")
	}
	//t.mu.Lock()
	//defer t.mu.Unlock()
	delete(t.actuators, name)
	t.actuators[name] = act
	return nil
	//}
	//return errors.New(fmt.Sprintf("invalid argument : actuator not found [%v]", name))
}

func (t *table) count() int {
	t.mu.RLock()
	defer t.mu.RUnlock()
	return len(t.actuators)
}

func (t *table) isEmpty() bool {
	return t.count() == 0
}

func (t *table) remove(name string) {
	if name == "" {
		return
	}
	t.mu.Lock()
	delete(t.actuators, name)
	t.mu.Unlock()
}
