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
	WriteAccessLog bool
	PingTraffic    bool
	rateLimiter    *rate.Limiter
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
	if r == nil || r.rateLimiter == nil {
		return true
	}
	return r.rateLimiter.Allow()
}

func (r *Route) IsRateLimiter() bool {
	return r != nil && r.rateLimiter != nil
}
