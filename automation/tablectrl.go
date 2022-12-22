package automation

import (
	"golang.org/x/time/rate"
)

func (t *table) enableFailover(name string, enabled bool) {
	if name == "" {
		return
	}
	t.mu.Lock()
	defer t.mu.Unlock()
	if act, ok := t.actuators[name]; ok {
		//fc := cloneFailover(act.failover)
		//fc.enabled = enabled
		//t.update(name, cloneActuator[*failover](act, fc))
		act.failover.enabled = enabled
		//t.actuators[name] = act
	}
}

func (t *table) setFailoverInvoke(name string, fn FailoverInvoke, enable bool) {
	if name == "" {
		return
	}
	t.mu.Lock()
	defer t.mu.Unlock()
	if act, ok := t.actuators[name]; ok {
		fc := cloneFailover(act.failover)
		fc.enabled = true
		fc.invoke = fn
		fc.enabled = enable
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
		//tc := cloneTimeout(act.timeout)
		//tc.enabled = enabled
		//t.update(name, cloneActuator[*timeout](act, tc))
		act.timeout.enabled = enabled
		//t.actuators[name] = act
	}
}

/*
func (t *table) setTimeout(name string, to time.Duration, enable bool) {
	if name == "" {
		return
	}
	t.mu.Lock()
	defer t.mu.Unlock()
	if act, ok := t.actuators[name]; ok {
		tc := cloneTimeout(act.timeout)
		tc.current.timeout = to
		tc.enabled = enable
		t.update(name, cloneActuator[*timeout](act, tc))
	}
}


*/
func (t *table) enableRateLimiter(name string, enabled bool) {
	if name == "" {
		return
	}
	t.mu.Lock()
	defer t.mu.Unlock()
	if act, ok := t.actuators[name]; ok {
		//lc := cloneRateLimiter(act.rateLimiter)
		//lc.enabled = enabled
		//t.update(name, cloneActuator[*rateLimiter](act, lc))
		act.rateLimiter.enabled = enabled
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

func (t *table) setRateLimiter(name string, config RateLimiterConfig) {
	if name == "" {
		return
	}
	t.mu.Lock()
	defer t.mu.Unlock()
	if act, ok := t.actuators[name]; ok {
		lc := cloneRateLimiter(act.rateLimiter)
		lc.currentConfig.limit = config.limit
		lc.currentConfig.burst = config.burst
		t.update(name, cloneActuator[*rateLimiter](act, lc))
	}
}

func (t *table) enableCircuitBreaker(name string, enabled bool) {
	if name == "" {
		return
	}
	t.mu.Lock()
	defer t.mu.Unlock()
	if act, ok := t.actuators[name]; ok {
		//c := cloneCircuitBreaker(act.circuitBreaker)
		//c.enabled = enabled
		//t.update(name, cloneActuator[*circuitBreaker](act, c))
		act.circuitBreaker.enabled = enabled
	}
}

func (t *table) enableRetry(name string, enabled bool) {
	if name == "" {
		return
	}
	t.mu.Lock()
	defer t.mu.Unlock()
	if act, ok := t.actuators[name]; ok {
		//c := cloneRetry(act.retry)
		//c.enabled = enabled
		//t.update(name, cloneActuator[*retry](act, c))
		act.retry.enabled = enabled
	}
}
