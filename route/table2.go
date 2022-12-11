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
	if r, ok := t.routes[name]; ok {
		if r.IsRateLimiter() {
			r.validateLimiter(&max, nil)
			r.rateLimiter.SetLimit(max)
		}

	}
	t.mu.Unlock()
}

func (t *table) SetBurst(name string, burst int) {
	if name == "" {
		return
	}
	t.mu.Lock()
	if r, ok := t.routes[name]; ok {
		if r.IsRateLimiter() {
			r.validateLimiter(nil, &burst)
			r.rateLimiter.SetBurst(burst)
		}
	}
	t.mu.Unlock()
}

/*
func (t *table) SetLimiter(name string, max rate.Limit, burst int) error {
	if name == "" {
		return nil
	}
	t.mu.Lock()
	if r, ok := t.routes[name]; ok {
		err := r.validateLimiter(&max, &burst)
		if err != nil {
			return nil
		}
		if r.rateLimiter == nil {
			r.rateLimiter = rate.NewLimiter(max, burst)
			return nil
		}
		if r.rateLimiter != nil {
			if max >= 0 {
				r.current.limit = max
				r.rateLimiter.SetLimit(max)
			}
			if burst > 0 {
				r.current.burst = burst
				r.rateLimiter.SetBurst(burst)
			}
		}
	}
	t.mu.Unlock()
	return nil
}


*/
func (t *table) ResetLimiter(name string) {
	if name == "" {
		return
	}
	t.mu.Lock()
	if r, ok := t.routes[name]; ok {
		if r.rateLimiter != nil {
			r.current.limit = r.default_.limit
			r.current.burst = r.default_.burst
			r.rateLimiter.SetLimit(r.current.limit)
			r.rateLimiter.SetBurst(r.current.burst)
		}
	}
	t.mu.Unlock()
}

func (t *table) DisableLimiter(name string) {
	if name == "" {
		return
	}
	t.mu.Lock()
	if r, ok := t.routes[name]; ok {
		if r.rateLimiter != nil {
			r.current.limit = rate.Inf
			r.rateLimiter.SetLimit(r.current.limit)
		}
	}
	t.mu.Unlock()
}

/*
func (t *table) RemoveLimiter(name string) bool {
	if t == nil || name == "" {
		return false
	}
	t.mu.Lock()
	if r, ok := t.routes[name]; ok {
		r.rateLimiter = nil
		return true
	}
	t.mu.Unlock()
	return false
}

*/
