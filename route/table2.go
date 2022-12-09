package route

import "golang.org/x/time/rate"

func (t *table) SetTimeout(name string, timeout int) bool {
	if t == nil || IsEmpty(name) {
		return false
	}
	t.mu.Lock()
	if r, ok := t.routes[name]; ok {
		if timeout <= 0 {
			timeout = NilValue
		}
		r.current.timeout = timeout
	}
	t.mu.Unlock()
	return true
}

func (t *table) ResetTimeout(name string) bool {
	if t == nil || IsEmpty(name) {
		return false
	}
	t.mu.Lock()
	if r, ok := t.routes[name]; ok {
		r.current.timeout = r.original.timeout
	}
	t.mu.Unlock()
	return true
}

func (t *table) DisableTimeout(name string) bool {
	return t.SetTimeout(name, NilValue)
}

func (t *table) SetLimiter(name string, max rate.Limit, burst int) bool {
	if t == nil || IsEmpty(name) {
		return false
	}
	t.mu.Lock()
	if r, ok := t.routes[name]; ok {
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
	return true
}

func (t *table) ResetLimiter(name string) bool {
	if t == nil || IsEmpty(name) {
		return false
	}
	t.mu.Lock()
	if r, ok := t.routes[name]; ok {
		if r.rateLimiter != nil {
			r.current.limit = r.original.limit
			r.current.burst = r.original.burst
			r.rateLimiter.SetLimit(r.current.limit)
			r.rateLimiter.SetBurst(r.current.burst)
		}
	}
	t.mu.Unlock()
	return true
}

func (t *table) DisableLimiter(name string) bool {
	if t == nil || IsEmpty(name) {
		return false
	}
	t.mu.Lock()
	if r, ok := t.routes[name]; ok {
		if r.rateLimiter != nil {
			r.current.limit = rate.Inf
			r.rateLimiter.SetLimit(r.current.limit)
		}
	}
	t.mu.Unlock()
	return true
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
