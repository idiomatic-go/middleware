package automation

import (
	"golang.org/x/time/rate"
	"net/http"
)

type NewActions interface {
	NewTimeout(timeout int) TimeoutAction
	NewPing(enable bool) Action
	NewRateLimit(max rate.Limit, burst int) RateLimitAction
}

type Actuators interface {
	Lookup(req *http.Request) Actuator
	LookupByName(name string) Actuator

	Exists(name string) bool
	Add(name string, ping bool, t *TimeoutConfig, r *RateLimitConfig) bool
	Remove(name string)
}

type Automation interface {
	SetDefault(a Actuator)
	SetMatcher(fn Matcher)
	IsPing(name string) bool
	NewActions
	Actuators
}

var Ingress = NewTable()
var Egress = NewTable()
