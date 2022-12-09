package route

import (
	"golang.org/x/time/rate"
	"net/http"
	"time"
)

const (
	DefaultName = "*"
	NilValue    = -1
)

type Matcher func(req *http.Request) (name string)

type Route interface {
	IsDefault() bool
	IsTimeout() bool
	IsLogging() bool
	IsRateLimiter() bool
	IsPingTraffic() bool
	Timeout() int
	Duration() time.Duration
	Allow() bool
	Name() string
	Limit() rate.Limit
	Burst() int
	t() route
}

type config struct {
	timeout int // milliseconds
	limit   rate.Limit
	burst   int // must not be "0", which will disallow all
}

type route struct {
	name           string
	current        config
	original       config
	writeAccessLog bool
	pingTraffic    bool
	rateLimiter    *rate.Limiter
}

func NewRoute(name string) Route {
	if IsEmpty(name) {
		return nil
	}
	return NewRouteWithConfig(name, NilValue, NilValue, NilValue, false, false)
}

func NewRouteWithLogging(name string, accessLog bool) Route {
	if IsEmpty(name) {
		return nil
	}
	return NewRouteWithConfig(name, NilValue, NilValue, NilValue, accessLog, false)
}

func NewRouteWithConfig(name string, timeout int, limit rate.Limit, burst int, accessLog, pingTraffic bool) Route {
	if IsEmpty(name) {
		return nil
	}
	route := &route{name: name, writeAccessLog: accessLog, pingTraffic: pingTraffic}
	if timeout == 0 {
		timeout = NilValue
	}
	if limit == 0 {
		limit = NilValue
	}
	if burst == 0 {
		burst = NilValue
	}
	route.original.timeout = timeout
	route.original.limit = limit
	route.original.burst = burst
	route.current = route.original
	return route
}

func (r *route) IsDefault() bool {
	return r != nil && r.name == DefaultName
}

func (r *route) IsTimeout() bool {
	return r != nil && r.current.timeout != NilValue
}

func (r *route) Timeout() int {
	if r == nil {
		return NilValue
	}
	return r.current.timeout
}

func (r *route) IsLogging() bool {
	return r != nil && r.writeAccessLog
}

func (r *route) IsRateLimiter() bool {
	return r != nil && r.rateLimiter != nil
}

func (r *route) IsPingTraffic() bool {
	return r != nil && r.pingTraffic
}

func (r *route) Duration() time.Duration {
	if r == nil || r.current.timeout == NilValue {
		return 0
	}
	return time.Duration(r.current.timeout) * time.Millisecond
}

func (r *route) Allow() bool {
	if r == nil || r.rateLimiter == nil {
		return true
	}
	return r.rateLimiter.Allow()
}

func (r *route) Limit() rate.Limit {
	if r == nil || r.rateLimiter == nil {
		return 0
	}
	return r.current.limit
}

func (r *route) Burst() int {
	if r == nil || r.rateLimiter == nil {
		return 0
	}
	return r.current.burst
}

func (r *route) Name() string {
	return r.name
}

func (r *route) t() route {
	return *r
}

func (r *route) newRateLimiter() {
	r.rateLimiter = rate.NewLimiter(r.current.limit, r.current.burst)
}

func (r *route) getRateLimiter() *rate.Limiter {
	return r.rateLimiter
}
