package automation

import (
	"net/http"
)

type Matcher func(req *http.Request) (routeName string)

type Configuration interface {
	SetMatcher(fn Matcher)
	SetDefault(name string, tc *TimeoutConfig, rlc *RateLimiterConfig, fc *FailoverConfig)
	Add(name string, tc *TimeoutConfig, rlc *RateLimiterConfig, fc *FailoverConfig) bool
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
