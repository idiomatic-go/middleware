package automation

import "golang.org/x/time/rate"

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

func (t *table) setLimit(name string, limit rate.Limit) {
	if name == "" {
		return
	}
	t.mu.Lock()
	defer t.mu.Unlock()
	if act, ok := t.actuators[name]; ok {
		lc := cloneRateLimiter(act.rateLimiter)
		lc.currentConfig.limit = limit
		t.update(name, cloneActuator[*rateLimiter](act, lc))
	}
}

func (t *table) setBurst(name string, burst int) {
	if name == "" {
		return
	}
	t.mu.Lock()
	defer t.mu.Unlock()
	if act, ok := t.actuators[name]; ok {
		lc := cloneRateLimiter(act.rateLimiter)
		lc.currentConfig.burst = burst
		t.update(name, cloneActuator[*rateLimiter](act, lc))
	}
}

func (t *table) setRateLimiter(name string, config RateLimiterConfig, canary bool) {
	if name == "" {
		return
	}
	t.mu.Lock()
	defer t.mu.Unlock()
	if act, ok := t.actuators[name]; ok {
		lc := cloneRateLimiter(act.rateLimiter)
		lc.currentConfig.limit = config.limit
		lc.currentConfig.burst = config.burst
		lc.canary = canary
		t.update(name, cloneActuator[*rateLimiter](act, lc))
	}
}
