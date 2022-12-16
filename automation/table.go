package automation

import (
	"golang.org/x/time/rate"
	"net/http"
	"sync"
)

type directory[T any] struct {
	entry map[string]T
}

type table struct {
	mu        sync.RWMutex
	Default   Actuator
	match     Matcher
	actuators map[string]*actuator
	pings     map[string]*pingAction
}

func NewTable() Automation {
	return newTable()
}

func newTable() *table {
	t := new(table)
	t.pings = make(map[string]*pingAction, 100)
	t.actuators = make(map[string]*actuator, 100)
	//t.limits.entry= make(map[string]*rateLimitAction, 100)
	//t.actuators = make(map[string]Actuator, 100)
	//t.defaultRoute = newRoute(DefaultName)
	t.match = func(req *http.Request) (name string) {
		return ""
	}
	return t
}

func (t *table) SetDefault(a Actuator) {
	if a == nil {
		return
	}
	t.mu.Lock()
	t.Default = a
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

func (t *table) IsPing(name string) bool {
	return false
}

func (t *table) NewTimeout(timeout int) TimeoutAction {
	return newTimeoutAction(timeout, t)
}

func (t *table) NewPing(enable bool) Action {
	return nil
}

func (t *table) NewRateLimit(max rate.Limit, burst int) RateLimitAction {
	return nil
}

func (t *table) Lookup(req *http.Request) Actuator {
	return nil
}
func (t *table) LookupByName(name string) Actuator {
	return nil
}

func (t *table) Exists(name string) bool {
	return false
}
func (t *table) Add(a Actuator) bool {
	return false
}

func (t *table) Remove(name string) {

}

func (t *table) configureTimeout(v ...any) {
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
