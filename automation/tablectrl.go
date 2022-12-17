package automation

func (t *table) setTimeout(name string, timeout int) {
	if name == "" {
		return
	}
	t.mu.Lock()
	if a, ok := t.actuators[name]; ok {
		if timeout <= 0 {
			timeout = NilValue
		}

		clone := cloneTimeout(a)
		clone.current.timeout = timeout
		newAct := &actuator{name: name, timeout: clone}
		t.update(name, newAct)
	}
	t.mu.Unlock()
}

func (t *table) resetTimeout(name string) {
	if name == "" {
		return
	}
	t.mu.Lock()
	if a, ok := t.actuators[name]; ok {
		a.timeout.current.timeout = a.timeout.defaultC.timeout
	}
	t.mu.Unlock()
}

func (t *table) disableTimeout(name string) {
	t.setTimeout(name, NilValue)
}
