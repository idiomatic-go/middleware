package route

import (
	"golang.org/x/time/rate"
	"net/http"
	"time"
)

type MatchFn func(req *http.Request) (name string)

type Route struct {
	Name           string
	Timeout        int // milliseconds
	RateLimiter    *rate.Limiter
	WriteAccessLog bool
	Ping           bool
}

func (r *Route) IsTimeout() bool {
	return r != nil && r.Timeout != 0
}

func (r *Route) Duration() time.Duration {
	if r == nil {
		return 0
	}
	return time.Duration(r.Timeout) * time.Millisecond
}

func (r *Route) IsLogging() bool {
	return r != nil && r.WriteAccessLog
}

func (r *Route) Allow() bool {
	if r == nil || r.RateLimiter == nil {
		return true
	}
	return r.RateLimiter.Allow()
}

func (r *Route) IsRateLimiter() bool {
	return r != nil && r.RateLimiter != nil
}

/*
func (r *Route) Limiter() *rate.Limiter {
	if r == nil {
		return nil
	}
	return r.RateLimiter
}


*/
func (r *Route) NewLimiter(max rate.Limit, b int) {
	if r == nil {
		return
	}
	r.RateLimiter = rate.NewLimiter(max, b)
}

func (r *Route) RemoveLimiter() {
	if r != nil {
		r.RateLimiter = nil
	}
}
