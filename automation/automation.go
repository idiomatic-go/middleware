package automation

import (
	"net/http"
)

/*
type NewActions interface {
	NewTimeout(timeout int) TimeoutAction
	NewPing(enable bool) Action
	NewRateLimit(max rate.Limit, burst int) RateLimitAction
}
*/

type Configuration interface {
	SetDefault(name string, t *TimeoutConfig, r *RateLimitConfig)
	SetMatcher(fn Matcher)
	IsPingEnabled(name string) bool
	Add(name string, p *PingConfig, t *TimeoutConfig, r *RateLimitConfig) bool
	//Exists(name string) bool
	//Remove(name string)
}

type Actuators interface {
	Lookup(req *http.Request) Actuator
	LookupByName(name string) Actuator
}

type Automation interface {
	Configuration
	Actuators
}

var Ingress = NewTable()
var Egress = NewTable()
