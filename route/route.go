package route

import (
	"errors"
	"golang.org/x/time/rate"
	"time"
)

const (
	DefaultName  = "*"
	NilValue     = -1
	DefaultBurst = 1
)

type Route interface {
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
}

type config struct {
	timeout int        // milliseconds
	limit   rate.Limit // A value of "0" will disallow all and a value of rate.Inf will allow all
	burst   int        // must not be "0", which will disallow all
}

type route struct {
	name    string
	current config
	// Mangled for Intellij
	default_       config
	writeAccessLog bool
	pingTraffic    bool
	rateLimiter    *rate.Limiter
}

func NewRoute(name string) (Route, error) {
	return NewRouteWithConfig(name, NilValue, NilValue, NilValue, false, false)
}

func newRoute(name string) Route {
	r, _ := NewRoute(name)
	return r
}

func NewRouteWithLogging(name string, accessLog bool) (Route, error) {
	return NewRouteWithConfig(name, NilValue, NilValue, NilValue, accessLog, false)
}

func NewRouteWithTimeout(name string, timeout int, accessLog, pingTraffic bool) (Route, error) {
	return newRouteWithConfig(name, timeout, rate.Inf, DefaultBurst, accessLog, pingTraffic)
}

func NewRouteWithConfig(name string, timeout int, limit rate.Limit, burst int, accessLog, pingTraffic bool) (Route, error) {
	return newRouteWithConfig(name, timeout, limit, burst, accessLog, pingTraffic)
}

func newRouteWithConfig(name string, timeout int, limit rate.Limit, burst int, accessLog, pingTraffic bool) (*route, error) {
	if IsEmpty(name) {
		return nil, errors.New("invalid argument : route name is empty")
	}
	route := &route{name: name, default_: config{timeout: timeout, limit: limit, burst: burst}, writeAccessLog: accessLog, pingTraffic: pingTraffic}
	route.validate()
	route.current = route.default_
	route.rateLimiter = rate.NewLimiter(route.default_.limit, route.default_.burst)
	return route, nil
}

func (r *route) validate() {
	if r.default_.timeout <= 0 {
		r.default_.timeout = NilValue
	}
	r.validateLimiter(&r.default_.limit, &r.default_.burst)
}

func (r *route) validateLimiter(max *rate.Limit, burst *int) {
	if max != nil && *max <= 0 {
		*max = rate.Inf
	}
	if burst != nil && *burst <= 0 {
		*burst = DefaultBurst
	}
	//if *max == NilValue && *burst != NilValue {
	//	return errors.New("invalid argument : burst is configured but limit is not")
	//}
	//if *max != NilValue && *burst == NilValue {
	//	return errors.New("invalid argument : limit is configured but burst is not")
	//}
	//return nil
}

func (r *route) IsDefault() bool {
	return r.name == DefaultName
}

func (r *route) IsTimeout() bool {
	return r.current.timeout != NilValue
}

func (r *route) Timeout() int {
	return r.current.timeout
}

func (r *route) IsLogging() bool {
	return r.writeAccessLog
}

func (r *route) IsRateLimiter() bool {
	return r.rateLimiter != nil
}

func (r *route) IsPingTraffic() bool {
	return r.pingTraffic
}

func (r *route) Duration() time.Duration {
	if r.current.timeout == NilValue {
		return 0
	}
	return time.Duration(r.current.timeout) * time.Millisecond
}

func (r *route) Allow() bool {
	if r.current.limit == rate.Inf {
		return true
	}
	return r.rateLimiter.Allow()
}

func (r *route) Limit() rate.Limit {
	return r.current.limit
}

func (r *route) Burst() int {
	return r.current.burst
}

func (r *route) Name() string {
	return r.name
}

func (r *route) newRateLimiter() {
	r.rateLimiter = rate.NewLimiter(r.current.limit, r.current.burst)
}

func (r *route) getRateLimiter() *rate.Limiter {
	return r.rateLimiter
}
