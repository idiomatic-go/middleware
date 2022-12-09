package route

import (
	"golang.org/x/time/rate"
	"net/http"
	"time"
)

const (
	DefaultName = "*"
)

type MatchFn func(req *http.Request) (name string)

type Route interface {
	IsTimeout() bool
	IsLogging() bool
}

type config struct {
	timeout int // milliseconds
	limit   rate.Limit
	burst   int
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
	return &route{name: name, writeAccessLog: true}
}

func NewRouteWithConfig(name string, timeout int, limit rate.Limit, burst int) Route {
	if IsEmpty(name) {
		return nil
	}
	route := &route{name: name}
	route.original.timeout = timeout
	route.original.limit = limit
	route.original.burst = burst
	route.current = route.original
	return route
}

func (r *route) IsTimeout() bool {
	return r != nil && r.current.timeout != 0
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
	if r == nil {
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
