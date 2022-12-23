package automation

import (
	"errors"
	"net/http"
	"sync"
)

type table struct {
	egress     bool
	mu         sync.RWMutex
	match      Matcher
	hostAct    *actuator
	defaultAct *actuator
	actuators  map[string]*actuator
}

func NewEgressTable() Automation {
	this := newTable(true)
	return this
}

func NewIngressTable() Automation {
	this := newTable(false)
	return this
}

func newTable(egress bool) *table {
	t := new(table)
	t.egress = egress
	t.actuators = make(map[string]*actuator, 100)
	t.hostAct = newDefaultActuator(HostActuatorName)
	t.defaultAct = newDefaultActuator(DefaultActuatorName)
	/*if egress {
		t.defaultAct = newActuator(DefaultActuatorName, t, newTimeout(DefaultActuatorName, t, nil),
			newCircuitBreaker(DefaultActuatorName, t, nil),
			newFailover(DefaultActuatorName, t, nil))
	} else {
		t.defaultAct = newActuator(DefaultActuatorName, t, newTimeout(DefaultActuatorName, t, nil), newRateLimiter(DefaultActuatorName, t, nil))
	}

	*/
	t.match = func(req *http.Request) (name string) {
		return ""
	}
	return t
}

func (t *table) isEgress() bool { return t.egress }

func (t *table) SetHostActuator(config ...any) error {
	t.mu.Lock()
	defer t.mu.Unlock()
	act := newActuator(HostActuatorName, t, config...)
	err := act.validate(t.egress)
	if err != nil {
		return err
	}
	t.hostAct = act
	return nil
}

func (t *table) SetDefaultActuator(name string, config ...any) error {
	t.mu.Lock()
	defer t.mu.Unlock()
	if name == "" {
		name = DefaultActuatorName
	}
	act := newActuator(name, t, config...)
	err := act.validate(t.egress)
	if err != nil {
		return err
	}
	t.defaultAct = act
	return nil
}

func (t *table) SetMatcher(fn Matcher) {
	if fn == nil {
		return
	}
	t.mu.Lock()
	t.match = fn
	t.mu.Unlock()
}

func (t *table) Host() Actuator {
	return t.hostAct
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

func (t *table) Add(name string, config ...any) error {
	if IsEmpty(name) {
		return errors.New("invalid argument: name is empty")
	}
	t.mu.Lock()
	defer t.mu.Unlock()
	act := newActuator(name, t, config...)
	err := act.validate(t.egress)
	if err != nil {
		return err
	}
	t.actuators[name] = act
	return nil
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

func (t *table) update(name string, act *actuator) {
	if name == "" || act == nil {
		return
	}
	//t.mu.Lock()
	//defer t.mu.Unlock()
	delete(t.actuators, name)
	t.actuators[name] = act
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
