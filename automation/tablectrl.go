package automation

import (
	"golang.org/x/time/rate"
	"time"
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

/*
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


*/
func (t *table) enableTimeout(name string, enabled bool) {
	if name == "" {
		return
	}
	t.mu.Lock()
	defer t.mu.Unlock()
	if act, ok := t.actuators[name]; ok {
		c := cloneTimeout(act.timeout)
		c.enabled = enabled
		t.update(name, cloneActuator[*timeout](act, c))
		//act.timeout.enabled = enabled
		//t.actuators[name] = act
	}
}

func (t *table) setTimeout(name string, to time.Duration) {
	if name == "" {
		return
	}
	t.mu.Lock()
	defer t.mu.Unlock()
	if act, ok := t.actuators[name]; ok {
		c := cloneTimeout(act.timeout)
		c.currentConfig.timeout = to
		t.update(name, cloneActuator[*timeout](act, c))
	}
}

func (t *table) enableRateLimiter(name string, enabled bool) {
	if name == "" {
		return
	}
	t.mu.Lock()
	defer t.mu.Unlock()
	if act, ok := t.actuators[name]; ok {
		c := cloneRateLimiter(act.rateLimiter)
		c.enabled = enabled
		t.update(name, cloneActuator[*rateLimiter](act, c))
	}
}

func (t *table) setLimit(name string, limit rate.Limit) {
	if name == "" {
		return
	}
	t.mu.Lock()
	defer t.mu.Unlock()
	if act, ok := t.actuators[name]; ok {
		c := cloneRateLimiter(act.rateLimiter)
		c.currentConfig.limit = limit
		// Not cloning the limiter as an old reference will not cause stale data when logging
		c.rateLimiter.SetLimit(limit)
		t.update(name, cloneActuator[*rateLimiter](act, c))
	}
}

func (t *table) setBurst(name string, burst int) {
	if name == "" {
		return
	}
	t.mu.Lock()
	defer t.mu.Unlock()
	if act, ok := t.actuators[name]; ok {
		c := cloneRateLimiter(act.rateLimiter)
		c.currentConfig.burst = burst
		// Not cloning the limiter as an old reference will not cause stale data when logging
		c.rateLimiter.SetBurst(burst)
		t.update(name, cloneActuator[*rateLimiter](act, c))
	}
}

func (t *table) setRateLimiter(name string, config RateLimiterConfig) {
	if name == "" {
		return
	}
	t.mu.Lock()
	defer t.mu.Unlock()
	if act, ok := t.actuators[name]; ok {
		c := cloneRateLimiter(act.rateLimiter)
		c.currentConfig.limit = config.limit
		c.currentConfig.burst = config.burst
		c.rateLimiter = rate.NewLimiter(c.currentConfig.limit, c.currentConfig.burst)
		t.update(name, cloneActuator[*rateLimiter](act, c))
	}
}

func (t *table) enableRetry(name string, enabled bool) {
	if name == "" {
		return
	}
	t.mu.Lock()
	defer t.mu.Unlock()
	if act, ok := t.actuators[name]; ok {
		c := cloneRetry(act.retry)
		c.enabled = enabled
		t.update(name, cloneActuator[*retry](act, c))
		//act.retry.enabled = enabled
	}
}
