package automation

import (
	"net/http"
)

type Matcher func(req *http.Request) (routeName string)

type Configuration interface {
	SetDefault(name string, t *TimeoutConfig)
	SetMatcher(fn Matcher)
	Add(name string, t *TimeoutConfig) bool
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
