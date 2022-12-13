package route

import "golang.org/x/time/rate"

func (t *table) SetTimeout(name string, timeout int) {
	if name == "" {
		return
	}
	t.mu.Lock()
	if r, ok := t.routes[name]; ok {
		if timeout <= 0 {
			timeout = NilValue
		}
		r.current.timeout = timeout
	}
	t.mu.Unlock()
}

func (t *table) ResetTimeout(name string) {
	if name == "" {
		return
	}
	t.mu.Lock()
	if r, ok := t.routes[name]; ok {
		r.current.timeout = r.default_.timeout
	}
	t.mu.Unlock()
}

func (t *table) DisableTimeout(name string) {
	t.SetTimeout(name, NilValue)
}

func (t *table) SetLimit(name string, max rate.Limit) {
	if name == "" {
		return
	}
	t.mu.Lock()
	defer t.mu.Unlock()
	if r, ok := t.routes[name]; ok {
		if r.IsRateLimiter() {
			r.validateLimiter(&max, nil)
			r.current.limit = max
			r.rateLimiter.SetLimit(max)
		}
	}
}

func (t *table) ResetLimit(name string) {
	if name == "" {
		return
	}
	t.mu.Lock()
	defer t.mu.Unlock()
	if r, ok := t.routes[name]; ok {
		if r.IsRateLimiter() {
			r.current.limit = r.default_.limit
			r.rateLimiter.SetLimit(r.current.limit)
		}
	}
}

func (t *table) DisableLimiter(name string) {
	t.SetLimit(name, rate.Inf)
}

func (t *table) SetBurst(name string, burst int) {
	if name == "" {
		return
	}
	t.mu.Lock()
	defer t.mu.Unlock()
	if r, ok := t.routes[name]; ok {
		if r.IsRateLimiter() {
			r.validateLimiter(nil, &burst)
			r.current.burst = burst
			r.rateLimiter.SetBurst(burst)
		}
	}
}

func (t *table) ResetBurst(name string) {
	if name == "" {
		return
	}
	t.mu.Lock()
	defer t.mu.Unlock()
	if r, ok := t.routes[name]; ok {
		if r.IsRateLimiter() {
			r.current.burst = r.default_.burst
			r.rateLimiter.SetBurst(r.current.burst)
		}
	}
}
