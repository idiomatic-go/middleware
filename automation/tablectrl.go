package automation

func (t *table) enableFailover(name string, enabled bool) {
	if name == "" {
		return
	}
	t.mu.Lock()
	defer t.mu.Unlock()
	if act, ok := t.actuators[name]; ok {
		fc := cloneFailover(act.failover)
		fc.enabled = enabled
		t.update(name, cloneActuator[*failover](act, fc))
	}
}

func (t *table) enableTimeout(name string, enabled bool) {
	if name == "" {
		return
	}
	t.mu.Lock()
	defer t.mu.Unlock()
	if act, ok := t.actuators[name]; ok {
		tc := cloneTimeout(act.timeout)
		tc.enabled = enabled
		t.update(name, cloneActuator[*timeout](act, tc))
	}
}

func (t *table) setTimeout(name string, timeout int) {
	if name == "" {
		return
	}
	t.mu.Lock()
	defer t.mu.Unlock()
	if act, ok := t.actuators[name]; ok {
		if timeout <= 0 {
			timeout = NilValue
		}
		tc := cloneTimeout(act.timeout)
		tc.current.timeout = timeout
		t.update(name, cloneActuator(act, tc))
	}
}

func (t *table) enableRateLimiter(name string, enabled bool) {
	if name == "" {
		return
	}
	t.mu.Lock()
	defer t.mu.Unlock()
	if act, ok := t.actuators[name]; ok {
		lc := cloneRateLimiter(act.rateLimiter)
		lc.enabled = enabled
		t.update(name, cloneActuator[*rateLimiter](act, lc))
	}
}

func (t *table) setRateLimiterCanary(name string, enable bool) {
	if name == "" {
		return
	}
	t.mu.Lock()
	defer t.mu.Unlock()
	if act, ok := t.actuators[name]; ok {
		//clone := cloneActuator(act)
		//clone.limiter = cloneRateLimiter(a.limiter)
		t.update(name, cloneActuator[*rateLimiter](act, nil))
	}
}
