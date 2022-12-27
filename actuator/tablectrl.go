package actuator

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
		c := cloneFailover(act.failover)
		c.enabled = enabled
		t.update(name, cloneActuator[*failover](act, c))
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
/*
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


*/
func (t *table) setTimeout(name string, to time.Duration) {
	if name == "" {
		return
	}
	t.mu.Lock()
	defer t.mu.Unlock()
	if act, ok := t.actuators[name]; ok {
		c := cloneTimeout(act.timeout)
		c.config.timeout = to
		t.update(name, cloneActuator[*timeout](act, c))
	}
}

/*
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


*/
func (t *table) setRateLimit(name string, limit rate.Limit) {
	if name == "" {
		return
	}
	t.mu.Lock()
	defer t.mu.Unlock()
	if act, ok := t.actuators[name]; ok {
		c := cloneRateLimiter(act.rateLimiter)
		c.config.limit = limit
		// Not cloning the limiter as an old reference will not cause stale data when logging
		c.rateLimiter.SetLimit(limit)
		t.update(name, cloneActuator[*rateLimiter](act, c))
	}
}

func (t *table) setRateBurst(name string, burst int) {
	if name == "" {
		return
	}
	t.mu.Lock()
	defer t.mu.Unlock()
	if act, ok := t.actuators[name]; ok {
		c := cloneRateLimiter(act.rateLimiter)
		c.config.burst = burst
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
		c.config.limit = config.limit
		c.config.burst = config.burst
		c.rateLimiter = rate.NewLimiter(c.config.limit, c.config.burst)
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
	}
}

func (t *table) setRetryRateLimit(name string, limit rate.Limit, burst int) {
	if name == "" {
		return
	}
	t.mu.Lock()
	defer t.mu.Unlock()
	if act, ok := t.actuators[name]; ok {
		c := cloneRetry(act.retry)
		c.config.limit = limit
		// Not cloning the limiter as an old reference will not cause stale data when logging
		c.rateLimiter = rate.NewLimiter(limit, burst)
		t.update(name, cloneActuator[*retry](act, c))
	}
}
