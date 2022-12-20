package automation

func (t *table) enableFailover(name string, enabled bool) {
	if name == "" {
		return
	}
	t.mu.Lock()
	defer t.mu.Unlock()
	if a, ok := t.actuators[name]; ok {
		fc := cloneFailover(a.failover)
		fc.enabled = enabled
		clone := cloneActuator(a)
		clone.failover = fc
		t.update(name, clone)
	}
}

func (t *table) enableTimeout(name string, enabled bool) {
	if name == "" {
		return
	}
	t.mu.Lock()
	defer t.mu.Unlock()
	if a, ok := t.actuators[name]; ok {
		tc := cloneTimeout(a.timeout)
		tc.enabled = enabled
		clone := cloneActuator(a)
		clone.timeout = tc
		t.update(name, clone)
	}
}

func (t *table) setTimeout(name string, timeout int) {
	if name == "" {
		return
	}
	t.mu.Lock()
	defer t.mu.Unlock()
	if a, ok := t.actuators[name]; ok {
		if timeout <= 0 {
			timeout = NilValue
		}
		to := cloneTimeout(a.timeout)
		to.current.timeout = timeout
		clone := cloneActuator(a)
		clone.timeout = to
		t.update(name, clone)
	}
}

func (t *table) setRateLimiterCanary(name string, enable bool) {
	if name == "" {
		return
	}
	t.mu.Lock()
	defer t.mu.Unlock()
	if a, ok := t.actuators[name]; ok {
		clone := cloneActuator(a)
		clone.limiter = cloneRateLimiter(a.limiter)
		t.update(name, clone)
	}
}
