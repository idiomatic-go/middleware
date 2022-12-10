package route

import (
	"errors"
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
	// Mangled for Idea
	default_       config
	writeAccessLog bool
	pingTraffic    bool
	rateLimiter    *rate.Limiter
}

func NewRoute(name string) Route {
	r, _ := NewRouteWithConfig(name, NilValue, NilValue, NilValue, false, false)
	return r
}

func NewRouteWithLogging(name string, accessLog bool) (Route, error) {
	return NewRouteWithConfig(name, NilValue, NilValue, NilValue, accessLog, false)
}

func NewRouteWithConfig(name string, timeout int, limit rate.Limit, burst int, accessLog, pingTraffic bool) (Route, error) {
	return newRouteWithConfig(name, timeout, limit, burst, accessLog, pingTraffic)
}

func newRouteWithConfig(name string, timeout int, limit rate.Limit, burst int, accessLog, pingTraffic bool) (*route, error) {
	if IsEmpty(name) {
		return nil, errors.New("invalid argument : route name is empty")
	}
	route := &route{name: name, default_: config{timeout: timeout, limit: limit, burst: burst}, writeAccessLog: accessLog, pingTraffic: pingTraffic}
	err := route.validate()
	if err != nil {
		return nil, err
	}
	route.current = route.default_
	if route.default_.limit != NilValue && route.default_.burst != NilValue {
		route.rateLimiter = rate.NewLimiter(route.default_.limit, route.default_.burst)
	}
	return route, nil
}

func (r *route) validate() error {
	if r.default_.timeout <= 0 {
		r.default_.timeout = NilValue
	}
	if r.default_.limit <= 0 {
		r.default_.limit = NilValue
	}
	if r.default_.burst <= 0 {
		r.default_.burst = NilValue
	}
	// Special handling for rate.Inf
	if r.default_.limit == rate.Inf {
		if r.default_.burst <= 0 {
			r.default_.burst = 1
		}
		return nil
	}
	if r.default_.limit == NilValue && r.default_.burst != NilValue {
		return errors.New("invalid argument : burst is configured but limit is not")
	}
	if r.default_.limit != NilValue && r.default_.burst == NilValue {
		return errors.New("invalid argument : limit is configured but burst is not")
	}
	return nil
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
	if !r.IsRateLimiter() {
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
